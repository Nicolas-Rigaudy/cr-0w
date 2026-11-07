"""Background run sandbox: per-kind allowlists, credential-free env, denials (item_077)."""
import shlex
import sys

from conftest import BIN

sys.path.insert(0, str(BIN))
import crow_lib  # noqa: E402


def test_every_dispatch_kind_has_an_allowlist():
    for kind in ("review", "task", "retro", "oneshot"):
        assert kind in crow_lib.RUN_SANDBOX
        assert crow_lib.RUN_SANDBOX[kind]["mcp"] == []


def test_bash_rules_get_an_rtk_twin():
    allow = crow_lib.sandbox_allow("review")
    assert "Bash(git status:*)" in allow and "Bash(rtk git status:*)" in allow
    assert "Read" in allow and "Write" not in allow
    assert not any(a.startswith("mcp__") for a in allow)


def test_task_may_commit_but_never_push():
    allow = crow_lib.sandbox_allow("task")
    assert "Bash(git commit:*)" in allow and "Edit" in allow
    assert not any("push" in a for a in allow)
    assert "Bash(git push:*)" in crow_lib.RUN_DENY and "Bash(rtk git push:*)" in crow_lib.RUN_DENY


def test_sandbox_args_strict_mcp_and_lists():
    args = crow_lib.sandbox_args("retro")
    assert args[0] == "--strict-mcp-config"
    allow = args[args.index("--allowedTools") + 1].split(",")
    assert allow == crow_lib.sandbox_allow("retro")
    assert "WebFetch" in args[args.index("--disallowedTools") + 1].split(",")
    assert "--mcp-config" not in args


def test_listed_mcp_server_is_passed(monkeypatch):
    monkeypatch.setitem(crow_lib.RUN_SANDBOX, "retro", {**crow_lib.RUN_SANDBOX["retro"], "mcp": ["echo"]})
    monkeypatch.setattr(crow_lib, "MCP_SERVERS", {"echo": {"command": "echo"}})
    args = crow_lib.sandbox_args("retro")
    assert "mcp__echo" in args[args.index("--allowedTools") + 1]
    assert '"echo"' in args[args.index("--mcp-config") + 1]


def test_sandbox_env_drops_credentials_keeps_config_dirs():
    env = {"GH_TOKEN": "x", "GITLAB_TOKEN": "x", "AWS_SECRET_ACCESS_KEY": "x", "AWS_PROFILE": "p",
           "FAKE_PROBE_TOKEN": "x", "OPENAI_API_KEY": "x", "SSH_AUTH_SOCK": "/s",
           "GH_CONFIG_DIR": "/c", "GLAB_CONFIG_DIR": "/c", "PATH": "/bin", "HOME": "/h",
           "CLAUDE_CODE_DISABLE_BACKGROUND_TASKS": "1"}
    clean, dropped = crow_lib.sandbox_env(env, "task")
    assert dropped == sorted(["GH_TOKEN", "GITLAB_TOKEN", "AWS_SECRET_ACCESS_KEY", "AWS_PROFILE",
                              "FAKE_PROBE_TOKEN", "OPENAI_API_KEY", "SSH_AUTH_SOCK"])
    assert set(clean) == {"GH_CONFIG_DIR", "GLAB_CONFIG_DIR", "PATH", "HOME", "CLAUDE_CODE_DISABLE_BACKGROUND_TASKS"}


def test_keep_list_lets_an_approved_token_through(monkeypatch):
    monkeypatch.setattr(crow_lib, "RUN_KEEP_ENV", {"review": ["GH_TOKEN"]})
    assert "GH_TOKEN" in crow_lib.sandbox_env({"GH_TOKEN": "x"}, "review")[0]
    assert "GH_TOKEN" not in crow_lib.sandbox_env({"GH_TOKEN": "x"}, "task")[0]


def test_run_env_strips_only_with_a_kind(monkeypatch):
    monkeypatch.setenv("GH_TOKEN", "secret")
    assert crow_lib.run_env({"A": "1"})["GH_TOKEN"] == "secret"
    env = crow_lib.run_env({"A": "1"}, "review")
    assert "GH_TOKEN" not in env and env["A"] == "1"


def test_profile_args_roundtrip_through_shlex():
    assert crow_lib.run_profile("acct") == "acct-runs"
    parts = shlex.split(crow_lib.run_profile_args())
    assert parts[0] == "--strict-mcp-config" and len(parts) == 5
    allow = parts[2].split(",")
    for kind in crow_lib.RUN_SANDBOX:
        assert set(crow_lib.sandbox_allow(kind)) <= set(allow)
    assert len(allow) == len(set(allow))
    assert parts[4].split(",") == crow_lib.RUN_DENY


def test_run_denials_from_cdx_warnings():
    run = {"final_payload": {"warnings": [
        {"code": "permission_denials", "permission_denials": [
            {"tool_name": "Bash", "tool_input": {"command": "touch x"}},
            {"tool_name": "Write", "tool_input": {"file_path": "/p/f"}},
            {"tool_name": "mcp__echo__ping", "tool_input": {}}]},
        {"code": "other"}]}}
    assert crow_lib.run_denials(run) == ["Bash: touch x", "Write: /p/f", "mcp__echo__ping"]
    assert crow_lib.run_denials({}) == [] and crow_lib.run_denials(None) == []


def test_dispatch_hands_cdx_a_credential_free_env(tmp_path, monkeypatch):
    """crow-runs dispatch starts cdx without the hub's tokens (the run inherits cdx's env)."""
    import importlib.machinery
    import importlib.util
    (tmp_path / "local").mkdir()
    (tmp_path / "local" / "projects.yaml").write_text(f"projects:\n  - {{name: demo, path: {tmp_path}}}\n")
    (tmp_path / "bin").mkdir()
    monkeypatch.setenv("CROW_HOME", str(tmp_path))
    monkeypatch.setenv("CROW_ACCOUNT", "main")
    monkeypatch.setenv("GH_TOKEN", "secret")
    loader = importlib.machinery.SourceFileLoader("crow_runs_sb", str(BIN / "crow-runs"))
    spec = importlib.util.spec_from_loader("crow_runs_sb", loader)
    cr = importlib.util.module_from_spec(spec)
    loader.exec_module(cr)
    seen = []

    class R:
        returncode = 0
        stdout = '{"ok": true, "run_id": "sb000001-full"}'

    monkeypatch.setattr(cr.subprocess, "run", lambda argv, **k: seen.append((argv, k.get("env"))) or R())
    monkeypatch.setattr(cr.subprocess, "Popen", lambda *a, **k: None)
    monkeypatch.setattr(cr, "running_runs", lambda: [])
    monkeypatch.setattr(cr.crow_lib, "provider_of", lambda acct: "claude")
    monkeypatch.setattr(cr, "sh", lambda *a, **k: "")
    monkeypatch.setattr(sys, "argv", ["crow-runs", "dispatch", "demo", "t-1", "--prompt", "p"])
    cr.main()
    argv, env = next(x for x in seen if x[0][:2] == ["cdx", "run"])
    assert "GH_TOKEN" not in env and env["CLAUDE_CODE_DISABLE_BACKGROUND_TASKS"] == "1"


def test_crow_runs_lists_denials_under_a_finished_run():
    import importlib.machinery
    import importlib.util
    import os
    os.environ.setdefault("CROW_ACCOUNT", "test")
    loader = importlib.machinery.SourceFileLoader("crow_runs_dn", str(BIN / "crow-runs"))
    spec = importlib.util.spec_from_loader("crow_runs_dn", loader)
    cr = importlib.util.module_from_spec(spec)
    loader.exec_module(cr)
    row = {"run_id": "aaaaaaaa", "project": "p", "label": "l", "status": "done", "duration_s": 60, "ended_at": None,
           "commits": [], "cost_usd": None, "budget_usd": None, "denials": ["Bash: touch x"]}
    out = "\n".join(cr.render({"agents": [], "slots": {"busy": 0, "max": 2}, "finished": [row]}, color=False, width=200))
    assert "denied Bash: touch x" in out
    row["denials"] = []
    assert "denied" not in "\n".join(cr.render({"agents": [], "slots": {"busy": 0, "max": 2}, "finished": [row]},
                                                color=False, width=200))
