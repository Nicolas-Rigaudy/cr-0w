"""crow-focus and the crow_term backends: window matching, the dry run and the print fallback.
No terminal is ever opened: kitty tests fake setsid, the rest use the print backend."""
import os
import sys

import pytest

from conftest import BIN, run

sys.path.insert(0, str(BIN))
import crow_term as ct  # noqa: E402

LS = [{"tabs": [{"windows": [
    {"id": 1, "pid": 100, "foreground_processes": [{"pid": 100}]},
    {"id": 2, "pid": 200, "foreground_processes": [{"pid": 250}]},
]}]}]


def test_window_found_by_ancestor_or_foreground_process():
    assert ct.kitty_find_window(LS, 999, chain=[999, 300, 200]) == 2
    assert ct.kitty_find_window(LS, 250, chain=[250]) == 2
    assert ct.kitty_find_window(LS, 999, chain=[999, 1]) is None


def test_ancestors_stop_at_init(tmp_path):
    for pid, ppid in ((30, 20), (20, 10), (10, 1)):
        d = tmp_path / str(pid)
        d.mkdir()
        (d / "stat").write_text(f"{pid} (a b) S {ppid} 0")
    assert ct.ancestors(30, proc=str(tmp_path)) == [30, 20, 10]


def test_dry_run_needs_account_and_cwd(hub):
    ok = run(hub, "crow-focus", "--pid", "5", "--account", "main", "--cwd", "/tmp", "--dry-run")
    assert ok.returncode == 0 and "cdx main --dir /tmp --resume" in ok.stdout
    assert run(hub, "crow-focus", "--pid", "5", "--dry-run").returncode != 0


def test_send_dry_run_and_command(hub):
    ok = run(hub, "crow-focus", "--pid", "5", "--send", "yes", "--dry-run", env={"CROW_TERMINAL": "kitty"})
    assert ok.returncode == 0 and "'yes'" in ok.stdout and "pid 5" in ok.stdout and "kitty" in ok.stdout
    assert ct.kitty_send_args("unix:/s", 7, "no") == ["kitten", "@", "--to", "unix:/s", "send-text", "--match", "id:7", "--", "no\r"]
    assert run(hub, "crow-focus", "--send", "yes", "--dry-run").returncode != 0


# ---------- backend choice ----------

def test_choice_order_and_override(monkeypatch):
    have = set()
    monkeypatch.setattr(ct.shutil, "which", lambda name, **_: f"/bin/{name}" if name in have else None)
    assert ct.choose({}) == "print"
    have.add("tmux")
    assert ct.choose({}) == "print" and ct.choose({"TMUX": "/s,1,0"}) == "tmux"
    have.add("kitty")
    assert ct.choose({"TMUX": "/s,1,0"}) == "kitty"
    assert ct.choose({"CROW_TERMINAL": "Print"}) == "print"
    with pytest.raises(SystemExit, match="not a backend"):
        ct.choose({"CROW_TERMINAL": "xterm"})


def test_clean_env_drops_session_markers(monkeypatch):
    monkeypatch.setenv("CLAUDECODE", "1")
    monkeypatch.setenv("CDX_SESSION_NAME", "main")
    monkeypatch.setenv("ANTHROPIC_MODEL", "x")
    env = ct.clean_env({"CROW_ACCOUNT": "spare"})
    assert not any(k.startswith(("CLAUDE", "CDX_", "ANTHROPIC_")) for k in env)
    assert env["CROW_ACCOUNT"] == "spare" and env["HOME"]


# ---------- tmux (commands only) ----------

def test_tmux_pane_match_and_commands(monkeypatch):
    panes = "%1 100\n%7 200\nbad line\n"
    assert ct.tmux_find_pane(panes, [999, 200]) == "%7"
    assert ct.tmux_find_pane(panes, [999]) is None
    assert ct.tmux_send_args("%7", "yes") == [["tmux", "send-keys", "-t", "%7", "-l", "yes"],
                                             ["tmux", "send-keys", "-t", "%7", "Enter"]]
    monkeypatch.setenv("CDX_SESSION_NAME", "main")
    env = ct.clean_env({"CROW_ACCOUNT": "spare"})
    args = ct.tmux_open_args(["crow", "open", "demo"], "crow · demo", "/w", env)
    assert args[:6] == ["tmux", "new-window", "-n", "crow · demo", "-c", "/w"]
    wrap = args[6:]
    assert wrap[0] == "env" and wrap[-3:] == ["crow", "open", "demo"]
    assert "-u" in wrap and "CDX_SESSION_NAME" in wrap and "CROW_ACCOUNT=spare" in wrap


# ---------- print fallback ----------

def test_print_backend_prints_the_line():
    line = ct.open_window(["cdx", "main", "--dir", "/a b", "--resume"], "t", "/a b",
                          {"CROW_ACCOUNT": "spare"}, ["CROW_ACCOUNT"], backend="print")
    assert line == "cd '/a b' && CROW_ACCOUNT=spare cdx main --dir '/a b' --resume"
    assert ct.focus(5, backend="print") is None and ct.send(5, "yes", backend="print") is False


def test_focus_without_kitty_prints_the_resume_command(hub, tmp_path):
    env = {"CROW_TERMINAL": "print"}
    r = run(hub, "crow-focus", "--pid", "5", "--account", "main", "--cwd", str(tmp_path), env=env)
    assert r.returncode == 0, r.stderr
    assert f"cdx main --dir {tmp_path} --resume" in r.stdout and "new terminal" in r.stdout
    s = run(hub, "crow-focus", "--pid", "5", "--send", "yes", env=env)
    assert s.returncode == 1 and "no print window" in s.stderr


def test_script_open_prints_on_the_print_backend(tmp_path):
    import subprocess
    e = {**os.environ, "CROW_TERMINAL": "print", "CDX_SESSION_NAME": "x"}
    r = subprocess.run([sys.executable, str(BIN / "crow_term.py"), "open", "--title", "t", "--cwd", str(tmp_path),
                        "--set", "CROW_ACCOUNT=spare", "--", "crow", "open", "demo"], capture_output=True, text=True, env=e)
    assert r.returncode == 0, r.stderr
    assert r.stdout.strip() == f"cd {tmp_path} && CROW_ACCOUNT=spare crow open demo"
    b = subprocess.run([sys.executable, str(BIN / "crow_term.py"), "backend"], capture_output=True, text=True, env=e)
    assert b.stdout.strip() == "print"


def test_ancestors_use_ps_without_proc(tmp_path, monkeypatch):
    parents = {30: "20", 20: "1"}

    class R:
        def __init__(self, out):
            self.stdout = out
    monkeypatch.setattr(ct.subprocess, "run", lambda cmd, **_: R(parents.get(int(cmd[-1]), "")))
    assert ct.ancestors(30, proc=str(tmp_path / "no-proc")) == [30, 20]
