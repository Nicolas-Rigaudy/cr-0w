import json

from conftest import run, write_jsonl
from test_crow_usage_projects import day, usage


def ctx_turn(mid, cache_read, new=0, cache_write=0):
    return {"type": "assistant", "cwd": "/w/hub", "timestamp": f"{day(0)}T12:00:00.000Z",
            "message": {"id": mid, "role": "assistant", "model": "m", "content": [{"type": "text", "text": "ok"}],
                        "usage": {"input_tokens": new, "cache_read_input_tokens": cache_read,
                                  "cache_creation_input_tokens": cache_write, "output_tokens": 10}}}


def session(hub, sid, *turns):
    home, profiles = hub
    write_jsonl(profiles / "main" / "claude-home" / ".claude" / "projects" / "-w-hub" / f"{sid}.jsonl",
                [{"type": "user", "cwd": "/w/hub", "timestamp": f"{day(0)}T08:00:00.000Z",
                  "message": {"role": "user", "content": "go"}}, *turns])


def test_usage_reports_last_turn_context_and_level(hub):
    session(hub, "short", ctx_turn("a", 30_000, 10, 500))
    session(hub, "long", ctx_turn("b", 50_000), ctx_turn("c", 119_000, 900, 100))  # 120_000 exactly
    session(hub, "over", ctx_turn("d", 160_000, 1))
    by_id = {s["id"]: s for s in usage(hub)["sessions"]}
    assert by_id["short"]["context_tokens"] == 30_510 and by_id["short"]["context_level"] == "ok"
    assert by_id["long"]["context_tokens"] == 120_000 and by_id["long"]["context_level"] == "warn"
    assert by_id["over"]["context_level"] == "over"


def test_usage_context_ignores_subagent_turns(hub):
    home, profiles = hub
    session(hub, "parent", ctx_turn("a", 40_000))
    sub = profiles / "main" / "claude-home" / ".claude" / "projects" / "-w-hub" / "parent" / "subagents" / "x.jsonl"
    write_jsonl(sub, [ctx_turn("z", 200_000)])
    s = {s["id"]: s for s in usage(hub)["sessions"]}["parent"]
    assert s["context_tokens"] == 40_000


def fake_usage(tmp_path, sessions):
    f = tmp_path / "usage.json"
    f.write_text(json.dumps({"sessions": sessions}))
    return {"CROW_USAGE_JSON": str(f)}


def sess(sid, tokens, active=True):
    return {"id": sid, "account": "main", "project": "hub", "title": f"t-{sid}", "active": active, "context_tokens": tokens}


def test_ctx_lists_only_live_sessions_past_the_warning(hub, tmp_path):
    env = fake_usage(tmp_path, [sess("a1", 119_999), sess("b2", 130_000), sess("c3", 160_000), sess("d4", 300_000, active=False)])
    out = run(hub, "crow-ctx", "--json", env=env)
    rows = json.loads(out.stdout)
    assert [r["id"] for r in rows] == ["c3", "b2"]
    assert [r["level"] for r in rows] == ["over", "warn"]
    assert "/wrap" in rows[0]["advice"] and "handoff" in rows[0]["advice"]


def test_ctx_brief_and_empty(hub, tmp_path):
    r = run(hub, "crow-ctx", "--brief", env=fake_usage(tmp_path, [sess("b2", 130_000)]))
    assert r.stdout.strip() == "- hub: t-b2 is at 130k of context; wrap and split"
    assert run(hub, "crow-ctx", "--brief", env=fake_usage(tmp_path, [sess("a1", 10_000)])).stdout == ""
    assert "no live session" in run(hub, "crow-ctx", env=fake_usage(tmp_path, [])).stdout


def settings(profiles, name):
    return profiles / name / "claude-home" / ".claude" / "settings.json"


def test_cap_apply_sets_allowed_profiles_only_and_keeps_other_settings(hub):
    _, profiles = hub
    for name in ("main", "shared", "mate"):
        settings(profiles, name).parent.mkdir(parents=True)
    settings(profiles, "main").write_text(json.dumps({"env": {"KEEP": "1"}, "hooks": {}}))
    assert run(hub, "crow-ctx", "cap", "check").returncode == 1
    r = run(hub, "crow-ctx", "cap", "apply")
    assert r.returncode == 0 and "main: set 150000" in r.stdout
    data = json.loads(settings(profiles, "main").read_text())
    assert data["env"] == {"KEEP": "1", "CLAUDE_CODE_AUTO_COMPACT_WINDOW": "150000"} and data["hooks"] == {}
    assert json.loads(settings(profiles, "shared").read_text())["env"]["CLAUDE_CODE_AUTO_COMPACT_WINDOW"] == "150000"
    assert not settings(profiles, "mate").exists()  # policy never
    assert run(hub, "crow-ctx", "cap", "check").returncode == 0
    again = run(hub, "crow-ctx", "cap", "apply")
    assert "main: already 150000" in again.stdout
