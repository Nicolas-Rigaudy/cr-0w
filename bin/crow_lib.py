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
