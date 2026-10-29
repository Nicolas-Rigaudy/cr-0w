"""crow-runs inbox: finished runs are listed once, then marked seen."""
import importlib.machinery
import importlib.util
import os

from conftest import BIN

os.environ.setdefault("CROW_ACCOUNT", "test")
loader = importlib.machinery.SourceFileLoader("crow_runs", str(BIN / "crow-runs"))
spec = importlib.util.spec_from_loader("crow_runs", loader)
cr = importlib.util.module_from_spec(spec)
loader.exec_module(cr)


def setup(monkeypatch, tmp_path):
    monkeypatch.setattr(cr, "SEEN", tmp_path / "seen.json")
    monkeypatch.setattr(cr, "ledger", lambda: [
        {"run_id": "aaaa1111", "project": "p", "label": "one"},
        {"run_id": "bbbb2222", "project": "p", "label": "two"},
    ])
    monkeypatch.setattr(cr, "finished", lambda since, ttl=8: {
        "aaaa1111": {"run_id": "aaaa1111", "exit_code": 0, "ended_at": "t1"},
    })


def test_inbox_lists_finished_run_once(monkeypatch, tmp_path, capsys):
    setup(monkeypatch, tmp_path)
    cr.inbox(hook=True)
    out = capsys.readouterr().out
    assert "UNREPORTED" in out and "aaaa1111" in out and "bbbb2222" not in out
    cr.inbox(hook=True)
    assert capsys.readouterr().out == ""


def test_peek_does_not_mark_seen(monkeypatch, tmp_path, capsys):
    setup(monkeypatch, tmp_path)
    cr.inbox(peek=True)
    cr.inbox(peek=True)
    assert capsys.readouterr().out.count("aaaa1111") == 2
