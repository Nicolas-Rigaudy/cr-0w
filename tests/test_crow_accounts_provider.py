"""crow-accounts provider commands against a fake cdx on PATH."""
import json
import stat

from conftest import run

STATUS = {"rows": [{"session_name": "main", "provider": "claude"}, {"session_name": "spare", "provider": "codex"}]}


def fake_cdx(tmp_path):
    d = tmp_path / "fakebin"
    d.mkdir()
    f = d / "cdx"
    f.write_text(f"#!/bin/sh\ncat <<'EOF'\n{json.dumps(STATUS)}\nEOF\n")
    f.chmod(f.stat().st_mode | stat.S_IEXEC)
    return {"PATH": f"{d}:/usr/bin:/bin"}


def test_provider_and_permission(hub, tmp_path):
    env = fake_cdx(tmp_path)
    assert run(hub, "crow-accounts", "provider", "spare", env=env).stdout.strip() == "codex"
    assert run(hub, "crow-accounts", "permission", "spare", env=env).stdout.strip() == "default"
    assert run(hub, "crow-accounts", "permission", "main", env=env).stdout.strip() == "auto"


def test_hub_session_refused_on_codex_only(hub, tmp_path):
    env = fake_cdx(tmp_path)
    assert run(hub, "crow-accounts", "hub-session", "main", env=env).returncode == 0
    r = run(hub, "crow-accounts", "hub-session", "spare", env=env)
    assert r.returncode == 1 and "hub guard" in r.stderr


def test_check_model_mismatch_exits(hub, tmp_path):
    env = fake_cdx(tmp_path)
    assert run(hub, "crow-accounts", "check-model", "spare", "sonnet", env=env).returncode == 1
    assert run(hub, "crow-accounts", "check-model", "spare", "gpt-5.1", env=env).returncode == 0
