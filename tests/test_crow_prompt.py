"""crow-prompt: the PermissionRequest hook record, the pending check and the answer refusals.
No window is touched: answers run dry or hit the refusal paths before any key press."""
import importlib.machinery
import importlib.util
import json

from conftest import BIN

loader = importlib.machinery.SourceFileLoader("crow_prompt", str(BIN / "crow-prompt"))
spec = importlib.util.spec_from_loader("crow_prompt", loader)
cp = importlib.util.module_from_spec(spec)
loader.exec_module(cp)

PROJECTS = [
    {"name": "lambda", "scope": "work", "path": "/w/lambda", "remote": "git@gitlab.example:team/lambda.git"},
    {"name": "toy", "scope": "personal", "path": "/p/toy"},
]
WAITING = {"status": "waiting", "waitingFor": "permission prompt", "sessionId": "s1"}


def event(cmd, cwd="/p/toy", tool="Bash"):
    return {"hook_event_name": "PermissionRequest", "session_id": "s1", "cwd": cwd, "tool_name": tool,
            "tool_input": {"command": cmd, "description": "do it"} if tool == "Bash" else cmd}


def test_hook_records_the_exact_command_and_prints_nothing(tmp_path, capsys):
    assert cp.hook(json.dumps(event("date +%s > stamp.txt")), home=tmp_path, pid=42, projects=PROJECTS, now=1.0) == 0
    assert capsys.readouterr().out == ""
    p = json.loads((tmp_path / "local/state/prompts/42.json").read_text())
    assert (p["tool"], p["detail"], p["description"], p["work_approval"]) == ("Bash", "date +%s > stamp.txt", "do it", None)
    assert [o["key"] for o in p["options"]] == ["allow", "deny"]


def test_each_prompt_gets_a_new_id(tmp_path):
    cp.hook(json.dumps(event("ls")), home=tmp_path, pid=42, projects=PROJECTS, now=1.0)
    first = cp.pending(42, home=tmp_path, meta=WAITING)["id"]
    cp.hook(json.dumps(event("ls")), home=tmp_path, pid=42, projects=PROJECTS, now=2.0)
    assert cp.pending(42, home=tmp_path, meta=WAITING)["id"] != first


def test_hook_survives_garbage(tmp_path, capsys):
    assert cp.hook("not json", home=tmp_path, pid=42, projects=PROJECTS) == 0
    assert capsys.readouterr().out == ""


def test_detail_falls_back_to_file_then_raw_input():
    assert cp.detail_of("Edit", {"file_path": "/a/b.py", "old_string": "x"}) == "/a/b.py"
    assert cp.detail_of("mcp__x__y", {"n": 1}) == '{"n": 1}'


def test_work_push_or_mr_is_an_approval():
    assert cp.work_approval("Bash", {"command": "git push origin feat"}, "/w/lambda.wt/feat", PROJECTS) == "lambda"
    assert cp.work_approval("Bash", {"command": "glab mr create --fill"}, "/w/lambda", PROJECTS) == "lambda"
    assert cp.work_approval("Bash", {"command": "git -C /w/lambda push"}, "/p/toy", PROJECTS) == "lambda"
    assert cp.work_approval("Bash", {"command": "git push origin main"}, "/p/toy", PROJECTS) is None
    assert cp.work_approval("Bash", {"command": "git status"}, "/w/lambda", PROJECTS) is None
    rec = cp.record(event("git push", cwd="/w/lambda"), 42, PROJECTS, now=1.0)
    assert [o["key"] for o in rec["options"]] == ["deny"]


def test_pending_needs_the_session_waiting_on_a_permission(tmp_path):
    cp.hook(json.dumps(event("ls")), home=tmp_path, pid=42, projects=PROJECTS, now=1.0)
    assert cp.pending(42, home=tmp_path, meta=WAITING)
    assert cp.pending(42, home=tmp_path, meta={"status": "busy", "sessionId": "s1"}) is None
    assert cp.pending(42, home=tmp_path, meta={**WAITING, "sessionId": "other"}) is None
    assert cp.pending(43, home=tmp_path, meta=WAITING) is None


def test_answer_refuses_stale_or_work_allow_and_presses_nothing(tmp_path):
    cp.hook(json.dumps(event("ls")), home=tmp_path, pid=42, projects=PROJECTS, now=1.0)
    pid_ = cp.pending(42, home=tmp_path, meta=WAITING)["id"]
    assert cp.answer(42, "old", "allow", home=tmp_path, meta=WAITING)[0] == 1
    assert cp.answer(42, pid_, "allow", home=tmp_path, meta={"status": "idle"})[0] == 1
    assert cp.answer(42, pid_, "maybe", home=tmp_path, meta=WAITING)[0] == 2
    code, msg = cp.answer(42, pid_, "allow", dry=True, backend="tmux", home=tmp_path, meta=WAITING)
    assert code == 0 and "press 1" in msg
    assert cp.answer(42, pid_, "allow", backend="print", home=tmp_path, meta=WAITING)[0] == 1

    cp.hook(json.dumps(event("git push", cwd="/w/lambda")), home=tmp_path, pid=42, projects=PROJECTS, now=3.0)
    wid = cp.pending(42, home=tmp_path, meta=WAITING)["id"]
    code, msg = cp.answer(42, wid, "allow", home=tmp_path, meta=WAITING)
    assert code == 1 and "work approval on lambda" in msg
    assert cp.answer(42, wid, "deny", dry=True, backend="tmux", home=tmp_path, meta=WAITING)[0] == 0


def test_claude_pid_walks_up_to_the_claude_process(tmp_path):
    for pid, ppid, argv in ((30, 20, b"bash\0-c\0hook"), (20, 10, b"node\0/usr/bin/claude\0"), (10, 1, b"fish")):
        d = tmp_path / str(pid)
        d.mkdir()
        (d / "stat").write_text(f"{pid} (x) S {ppid} 0")
        (d / "cmdline").write_bytes(argv)
    assert cp.claude_pid(30, proc=str(tmp_path)) == 20
