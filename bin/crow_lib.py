"""crow_lib: helpers shared by the crow scripts."""
import json
import os
import pwd
import re
import subprocess
import sys


def real_home():
    return pwd.getpwuid(os.getuid()).pw_dir


def profiles_dir(env=None):
    """cdx profiles dir: CDX_PROFILES, else CDX_HOME/profiles, else ~/.cdx/profiles."""
    env = os.environ if env is None else env
    if env.get("CDX_PROFILES"):
        return os.path.expanduser(env["CDX_PROFILES"])
    root = env.get("CDX_HOME") or os.path.join(real_home(), ".cdx")
    return os.path.join(os.path.expanduser(root), "profiles")


def _json_objects(text):
    """Objects from a JSON document, a JSON array or JSON lines; other lines are skipped."""
    try:
        doc = json.loads(text)
        docs = doc if isinstance(doc, list) else [doc]
    except ValueError:
        docs = []
        for line in text.splitlines():
            try:
                docs.append(json.loads(line))
            except ValueError:
                continue
    return [d for d in docs if isinstance(d, dict)]


def final_text(text, provider="claude"):
    """Final message of a cdx run's stdout. Claude: the `result` envelope. Codex: the last
    agent_message event (`item.completed` or the older `msg` shape). Other output: its tail."""
    objs = _json_objects(text)
    found = None
    for o in objs:
        if provider == "claude" and o.get("type") == "result":
            found = o.get("result")
        elif provider == "codex":
            item = o.get("item") if o.get("type") == "item.completed" else o.get("msg")
            if isinstance(item, dict) and item.get("type") == "agent_message":
                found = item.get("text") or item.get("message")
    return found if found else text[-4000:]


# ---------- providers ----------

PROVIDERS = ("claude", "codex")
CLAUDE_MODEL = re.compile(r"^(claude|opus|sonnet|haiku|fable)", re.I)
CODEX_MODEL = re.compile(r"^(gpt|o\d|codex)", re.I)
# Read first by a provider that does not load CLAUDE.md or @imports itself.
READ_FIRST = ("CLAUDE.md", "AGENTS.md", "LOGICS.md", "local/owner-prefs.md", "local/context.md",
              "local/memory/MEMORY.md")


def provider_of(account, status=None):
    """Provider of a cdx account (claude when cdx does not say)."""
    if status is None:
        try:
            r = subprocess.run(["cdx", "status", "--cached", "--json"], capture_output=True, text=True)
            status = json.loads(r.stdout).get("rows", [])
        except (OSError, ValueError):
            status = []
    row = next((x for x in status if x.get("session_name") == account), {})
    return row.get("provider") if row.get("provider") in PROVIDERS else "claude"


def check_model(provider, model):
    """Refuse a model name that belongs to the other provider; return None when fine."""
    if model and provider == "codex" and CLAUDE_MODEL.match(model):
        return f"model '{model}' is a Claude model; the account runs Codex"
    if model and provider == "claude" and CODEX_MODEL.match(model):
        return f"model '{model}' is a Codex model; the account runs Claude"
    return None


def read_first(hub, cwd=None):
    """Instruction files Codex must read first (it does not follow CLAUDE.md or @imports)."""
    refs = [os.path.join(hub, f) for f in READ_FIRST if os.path.exists(os.path.join(hub, f))]
    if cwd and os.path.realpath(cwd) != os.path.realpath(hub):
        refs += [os.path.join(cwd, f) for f in ("AGENTS.md", "CLAUDE.md") if os.path.exists(os.path.join(cwd, f))]
    return refs


def skill_file(hub, name):
    for d in (".claude/skills", "skills", "local/skills"):
        p = os.path.join(hub, d, name, "SKILL.md")
        if os.path.exists(p):
            return p
    return None


def oneshot_prompt(provider, hub, skill, args=""):
    """Prompt that runs a hub skill. Claude takes the slash command; Codex gets the skill file
    and the read-first list. Exits when Codex has no skill file to follow."""
    if provider == "claude":
        return f"/{skill} {args}".strip()
    path = skill_file(hub, skill)
    if not path:
        sys.exit(f"crow: no skill file for '{skill}' under {hub}; Codex cannot run it")
    refs = ", ".join(read_first(hub))
    return (f"Read these files first: {refs}. Then read {path} and follow it as the `{skill}` skill. "
            f"{args}").strip()


def permission(provider):
    """cdx permission for hub one-shots: Codex keeps writes in the workspace, no network."""
    return "auto" if provider == "claude" else "default"


def hub_session_refusal(provider):
    """Why an interactive hub session cannot run on this provider (None when it can)."""
    if provider == "codex":
        return ("the hub guard (crow-hub-guard) is a Claude hook; an interactive Codex hub session would run "
                "unguarded. Use a Claude account, or `crow open <project>` for repo work.")
    return None


# ---------- keyring (secret-tool, attributes service + user) ----------

def secret_lookup(service, user=None):
    """Value stored for service (and user), "" when missing or the keyring is unreachable."""
    args = ["secret-tool", "lookup", "service", service] + (["user", user] if user else [])
    try:
        r = subprocess.run(args, capture_output=True, text=True)
    except OSError:
        return ""
    return r.stdout.rstrip("\n") if r.returncode == 0 else ""


def secret_store(service, user, value, label):
    """Store value on stdin (never argv); True on success."""
    args = ["secret-tool", "store", "--label", label, "service", service, "user", user]
    try:
        return subprocess.run(args, input=value, capture_output=True, text=True).returncode == 0
    except OSError:
        return False


def secret_clear(service, user):
    args = ["secret-tool", "clear", "service", service, "user", user]
    try:
        return subprocess.run(args, capture_output=True, text=True).returncode == 0
    except OSError:
        return False


def secret_items(service=None):
    """Keyring items as dicts (label, created, modified, service, user), values dropped.
    Without service: every generic item that carries a service attribute."""
    query = ["service", service] if service else ["xdg:schema", "org.freedesktop.Secret.Generic"]
    try:
        # secret-tool prints attributes on stderr: one pipe keeps each item's lines together.
        r = subprocess.run(["secret-tool", "search", "--all", *query], stdout=subprocess.PIPE,
                           stderr=subprocess.STDOUT, text=True)
    except OSError:
        return []
    items, cur = [], None
    for line in r.stdout.splitlines():
        if line.startswith("["):
            cur = {}
            items.append(cur)
            continue
        key, sep, val = line.partition(" = ")
        if cur is None or not sep or key == "secret":
            continue
        if key in ("label", "created", "modified"):
            cur[key] = val
        elif key.startswith("attribute."):
            cur[key[len("attribute."):]] = val
    return [i for i in items if i.get("service")]


# ---------- account picking (limits, policies, shared-account caps) ----------

NEARLY_OUT_PCT = 95   # 5h use at which crow stops starting new runs on an account
AUTO_FALLBACK_PCT = 90  # every `use` account past this: a `fallback: auto` account needs no OK


def used_pct(row, key):
    """Points used of a cdx status window (remaining_5h_pct / remaining_week_pct), None when unknown."""
    left = (row or {}).get(key)
    return None if left is None else 100 - left


def minutes_to(iso, now):
    if not iso:
        return None
    from datetime import datetime
    try:
        return (datetime.fromisoformat(iso) - now).total_seconds() / 60
    except ValueError:
        return None


def hard_block(row):
    """Why cdx cannot serve work on the account at all, else None."""
    if not row:
        return "unknown to cdx"
    if row.get("enabled") is False:
        return "disabled"
    if row.get("auth_status") not in (None, "authenticated"):
        return f"auth {row['auth_status']}"
    if row.get("remaining_5h_pct") == 0:
        return f"5h limit (resets {row.get('reset_5h_at')})"
    if row.get("remaining_week_pct") == 0:
        return f"weekly limit (resets {row.get('reset_week_at')})"
    return None


def cap_block(acct, row, now):
    """Why crow leaves a shared account alone (max_5h_pct / max_week_pct in accounts.yaml), else None.
    Past the 5h cap the account stays usable when its 5h window resets within free_5h_min minutes
    and the week is under its cap: nobody else is likely to need the rest of that window."""
    u5, uw = used_pct(row, "remaining_5h_pct"), used_pct(row, "remaining_week_pct")
    mx5, mxw = acct.get("max_5h_pct"), acct.get("max_week_pct")
    if mxw is not None and uw is not None and uw >= mxw:
        return f"week {uw}% used, crow's cap is {mxw}%"
    if mx5 is not None and u5 is not None and u5 >= mx5:
        left = minutes_to((row or {}).get("reset_5h_at_iso"), now)
        free = acct.get("free_5h_min")
        if free and left is not None and 0 <= left <= free:
            return None
        return f"5h {u5}% used, crow's cap is {mx5}% (resets {(row or {}).get('reset_5h_at')})"
    return None


def pick_accounts(accts, rows, current, now):
    """Every allowed account with its state and the targets crow may pick.

    accts: crow-accounts json entries; rows: cdx status rows by name. Returns
    {current, current_blocked, accounts[], auto_target, recommended}: auto_target needs no OK
    (policy use, or `fallback: auto` once every use account is blocked or past AUTO_FALLBACK_PCT);
    recommended is the best unblocked account, ask-first included."""
    out = []
    for a in accts:
        if a["policy"] == "never":
            continue
        row = rows.get(a["name"])
        out.append({"name": a["name"], "policy": a["policy"], "fallback": a.get("fallback"),
                    "current": a["name"] == current,
                    "blocked": hard_block(row) or cap_block(a, row, now),
                    "used_5h_pct": used_pct(row, "remaining_5h_pct"), "used_week_pct": used_pct(row, "remaining_week_pct"),
                    "remaining_5h_pct": (row or {}).get("remaining_5h_pct"),
                    "remaining_week_pct": (row or {}).get("remaining_week_pct"),
                    "reset_5h_at": (row or {}).get("reset_5h_at"), "reset_5h_at_iso": (row or {}).get("reset_5h_at_iso"),
                    "reset_week_at": (row or {}).get("reset_week_at")})
    uses = [a for a in out if a["policy"] == "use"]
    uses_spent = all(a["blocked"] or (a["used_5h_pct"] or 0) >= AUTO_FALLBACK_PCT for a in uses)
    for a in out:
        a["auto_ok"] = a["policy"] == "use" or (a["fallback"] == "auto" and a["policy"] == "ask-first" and uses_spent)
    usable = [a for a in out if not a["blocked"] and not a["current"]]
    usable.sort(key=lambda a: (a["policy"] != "use", not a["auto_ok"],
                               -min(a["remaining_5h_pct"] or 0, a["remaining_week_pct"] or 0)))
    cur = next((a for a in out if a["current"]), None)
    nearly_out = bool(cur and (cur["used_5h_pct"] or 0) >= NEARLY_OUT_PCT)
    return {"current": current, "current_blocked": (cur["blocked"] if cur else hard_block(rows.get(current)))
            or (f"5h {cur['used_5h_pct']}% used" if nearly_out else None),
            "accounts": out, "auto_target": next((a["name"] for a in usable if a["auto_ok"]), None),
            "recommended": usable[0]["name"] if usable else None}


def account_state(current=None, cached=False):
    """pick_accounts on live data: crow-accounts json and `cdx status` (--cached for hooks)."""
    from datetime import datetime
    bin_dir = os.path.dirname(os.path.realpath(__file__))
    accts = json.loads(subprocess.run([os.path.join(bin_dir, "crow-accounts"), "json"], capture_output=True,
                                      text=True, timeout=30).stdout)
    current = current or os.environ.get("CDX_SESSION_NAME") or os.environ.get("CROW_ACCOUNT") or accts.get("default")
    try:
        out = subprocess.run(["cdx", "status", *(["--cached"] if cached else []), "--json"],
                             capture_output=True, text=True, timeout=30).stdout
        rows = {r["session_name"]: r for r in json.loads(out).get("rows") or [] if r.get("session_name")}
    except (OSError, ValueError, subprocess.TimeoutExpired):
        rows = {}
    return pick_accounts(accts["accounts"], rows, current, datetime.now().astimezone())


def room(a):
    """Quota left on an account row of pick_accounts: the smaller of its 5h and weekly windows."""
    return min(a["remaining_5h_pct"] or 0, a["remaining_week_pct"] or 0)


def launch_pick(state, busy, default):
    """Account for an interactive launch (crow open, term, handoff).

    state: pick_accounts output with no current account; busy: names with a live cdx session.
    Returns {account, ask, skipped}: account is None when none is free and has room; ask is True
    when it is an ask-first account the owner must OK. The default stays the pick while it is free,
    policy use and under NEARLY_OUT_PCT; otherwise the free account with the most room, use first."""
    skipped, free = [], []
    for a in state["accounts"]:
        why = ("live session" if a["name"] in busy else None) or a["blocked"]
        if not why and a["policy"] == "use" and (a["used_5h_pct"] or 0) >= NEARLY_OUT_PCT:
            why = f"5h {a['used_5h_pct']}% used"
        (skipped.append({"name": a["name"], "why": why}) if why else free.append(a))
    free.sort(key=lambda a: (a["policy"] != "use", not a["auto_ok"], -room(a), a["name"]))
    dflt = next((a for a in free if a["name"] == default and a["policy"] == "use"), None)
    top = dflt or (free[0] if free else None)
    return {"account": top and top["name"], "ask": bool(top and not top["auto_ok"]),
            "room": top and room(top), "skipped": skipped}


def live_sessions():
    """Names of cdx accounts with a live interactive session (`cdx sessions --json`)."""
    try:
        out = subprocess.run(["cdx", "sessions", "--json"], capture_output=True, text=True, timeout=30).stdout
        return {s["session"] for s in json.loads(out).get("sessions") or [] if s.get("active") and s.get("session")}
    except (OSError, ValueError, subprocess.TimeoutExpired):
        return set()


# ---------- headless run flags (req_017 levers 3, 4, 7, 16) ----------

CLAUDE_MODEL_IDS = {"opus": "claude-opus-5-5", "sonnet": "claude-sonnet-5-5", "haiku": "claude-haiku-5-5",
                    "fable": "claude-fable-5-1"}
# kind: default model, effort, budget in USD, wall-clock cap in minutes. Env CROW_<KIND>_BUDGET / _TIMEOUT_MIN override.
RUN_KINDS = {
    "task": {"model": "opus", "power": "medium", "budget": 25.0, "timeout_min": 120},
    "review": {"model": "opus", "power": "high", "budget": 10.0, "timeout_min": 60},
    "retro": {"model": "sonnet", "power": "medium", "budget": 2.0, "timeout_min": 20},
    "oneshot": {"model": "opus", "power": "medium", "budget": 3.0, "timeout_min": 20},
}
# crow kinds that map to a cdx --kind (cdx knows assistant|code-review|retro; the rest keep its default).
CDX_KINDS = {"retro": "retro"}
# Headless traps: no background tasks (phantom results), no claude.ai connectors (~40k tokens of tool defs per call).
HEADLESS_ENV = {"CLAUDE_CODE_DISABLE_BACKGROUND_TASKS": "1", "ENABLE_CLAUDEAI_MCP_SERVERS": "false"}


def model_id(model):
    """Exact Claude model id for an alias (opus -> claude-opus-5-5); other values pass through."""
    return CLAUDE_MODEL_IDS.get((model or "").strip().lower(), model)


def run_flags(kind, provider="claude", model=None, power=None, continue_run=None, env=None, budget=None):
    """cdx run arguments and environment every headless launch carries.
    With continue_run, cdx pins the parent's model, effort and account, so none of them nor --failover is passed.
    --failover moves a run stopped by a usage limit to another allowed account. cdx run has no --fallback-model
    flag: the fallback model is a per-account setting (`cdx set <account> --fallback-model M`).
    budget (USD) overrides the kind's budget for this run; Claude only."""
    env = os.environ if env is None else env
    k = RUN_KINDS[kind]
    budget = float(budget if budget is not None else env.get(f"CROW_{kind.upper()}_BUDGET") or k["budget"])
    timeout = int(env.get(f"CROW_{kind.upper()}_TIMEOUT_MIN") or k["timeout_min"])
    args = []
    if not continue_run:
        m = model or (k["model"] if provider == "claude" else None)
        args += (["--model", model_id(m) if provider == "claude" else m] if m else []) + ["--power", power or k["power"]]
    args += ["--timeout-seconds", str(timeout * 60)] + ([] if continue_run else ["--failover"])
    if kind in CDX_KINDS:
        args += ["--kind", CDX_KINDS[kind]]
    if provider == "claude":  # cdx measures spend from Claude usage; Codex budgets fail as unverifiable without model evidence
        args += ["--budget", f"{budget:g}"]
    return args, (dict(HEADLESS_ENV) if provider == "claude" else {})


def run_env(extra):
    """os.environ plus the headless overrides, for subprocess calls."""
    return {**os.environ, **extra}
