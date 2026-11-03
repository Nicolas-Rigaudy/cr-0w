"""crow-secret: keyring tokens via secret-tool; values never in argv or list output."""
import json
import os
import pty
import select
import subprocess
import sys
import textwrap

import pytest

from conftest import BIN, PYTHON, run

FAKE = textwrap.dedent('''\
    #!/usr/bin/env python3
    """Fake secret-tool: items in $FAKE_KEYRING (JSON), every argv appended to $FAKE_KEYRING.argv."""
    import json, os, sys
    path = os.environ["FAKE_KEYRING"]
    with open(path + ".argv", "a") as f:
        f.write(json.dumps(sys.argv[1:]) + "\\n")
    items = json.load(open(path)) if os.path.exists(path) else []
    args = sys.argv[1:]
    cmd = args.pop(0)
    label = None
    if args and args[0] == "--label":
        label = args[1]; args = args[2:]
    if args and args[0] == "--all":
        args = args[1:]
    attrs = dict(zip(args[::2], args[1::2]))
    match = [i for i in items if all(i["attrs"].get(k) == v for k, v in attrs.items() if k != "xdg:schema")]
    if cmd == "store":
        items = [i for i in items if i["attrs"] != attrs]
        items.append({"label": label, "attrs": attrs, "secret": sys.stdin.read()})
        json.dump(items, open(path, "w"))
    elif cmd == "lookup":
        if not match:
            sys.exit(1)
        sys.stdout.write(match[0]["secret"])
    elif cmd == "clear":
        json.dump([i for i in items if i not in match], open(path, "w"))
    elif cmd == "search":
        for n, i in enumerate(match):
            print(f"[/{n}]\\nlabel = {i['label']}\\nsecret = {i['secret']}\\ncreated = 2026-11-03 10:00:00\\n"
                  f"modified = 2026-11-03 11:00:00\\nschema = org.freedesktop.Secret.Generic", flush=True)
            for k, v in i["attrs"].items():
                print(f"attribute.{k} = {v}", file=sys.stderr, flush=True)
''')

REGISTRY = """\
tokens:
  - service: sonar
    user: me@example.com
    expiry: 2000-01-01
  - service: jira
    user: me@example.com
"""


@pytest.fixture
def keyring(hub, tmp_path):
    fakebin = tmp_path / "fakebin"
    fakebin.mkdir()
    tool = fakebin / "secret-tool"
    tool.write_text(FAKE.replace("/usr/bin/env python3", PYTHON or sys.executable))
    tool.chmod(0o755)
    store = tmp_path / "keyring.json"
    env = {"PATH": f"{fakebin}:{os.environ['PATH']}", "FAKE_KEYRING": str(store)}
    (hub[0] / "local" / "tokens.yaml").write_text(REGISTRY)
    return store, env


def seed(store, service, user, secret, label="crow token"):
    items = json.loads(store.read_text()) if store.exists() else []
    items.append({"label": label, "attrs": {"service": service, "user": user}, "secret": secret})
    store.write_text(json.dumps(items))


def secret(hub, env, *args, **kw):
    return run(hub, "crow-secret", *args, env=env, **kw)


def set_on_tty(hub, env, *args, typed="tok-123\n"):
    """Run `crow-secret set` with a pty as controlling terminal and type the value."""
    home, profiles = hub
    e = {k: v for k, v in os.environ.items() if k != "CROW_ACCOUNT"}
    e.update(CROW_HOME=str(home), CDX_PROFILES=str(profiles), **env)
    pid, fd = pty.fork()
    if pid == 0:
        os.execve(PYTHON, [PYTHON, str(BIN / "crow-secret"), "set", *args], e)
    out = b""
    while b"hidden" not in out:
        out += os.read(fd, 1024)
    os.write(fd, typed.encode())
    while True:
        r, _, _ = select.select([fd], [], [], 5)
        if not r:
            break
        try:
            chunk = os.read(fd, 1024)
        except OSError:
            break
        if not chunk:
            break
        out += chunk
    _, status = os.waitpid(pid, 0)
    return os.waitstatus_to_exitcode(status), out.decode()


def test_set_reads_the_tty_and_stores_with_registry_user(hub, keyring):
    store, env = keyring
    code, out = set_on_tty(hub, env, "sonar")
    assert code == 0, out
    assert "tok-123" not in out
    items = json.loads(store.read_text())
    assert items == [{"label": "crow sonar token", "attrs": {"service": "sonar", "user": "me@example.com"},
                      "secret": "tok-123"}]
    assert all("tok-123" not in line for line in (store.parent / "keyring.json.argv").read_text().splitlines())


def test_set_hints_the_registry_for_an_unknown_service(hub, keyring):
    _, env = keyring
    code, out = set_on_tty(hub, env, "other", "--user", "x@example.com")
    assert code == 0
    assert "Add 'other' to local/tokens.yaml" in out


def test_set_refuses_an_empty_value(hub, keyring):
    store, env = keyring
    code, out = set_on_tty(hub, env, "sonar", typed="\n")
    assert code == 1
    assert "empty value" in out
    assert not store.exists()


def test_set_refuses_a_value_on_argv_without_echoing_it(hub, keyring):
    store, env = keyring
    out = secret(hub, env, "set", "sonar", "tok-on-argv")
    assert out.returncode == 1
    assert "tok-on-argv" not in out.stdout + out.stderr
    assert not store.exists()


def test_set_refuses_without_a_terminal(hub, keyring):
    store, env = keyring
    home, profiles = hub
    e = {**os.environ, "CROW_HOME": str(home), **env}
    out = subprocess.run([PYTHON, str(BIN / "crow-secret"), "set", "sonar"], input="tok\n",
                         capture_output=True, text=True, env=e, start_new_session=True)
    assert out.returncode == 1
    assert "real terminal" in out.stderr
    assert not store.exists()


def test_get_prints_the_value_only(hub, keyring):
    store, env = keyring
    seed(store, "sonar", "me@example.com", "s3cret")
    out = secret(hub, env, "get", "sonar")
    assert (out.returncode, out.stdout) == (0, "s3cret")


def test_get_uses_the_only_stored_user_when_not_in_registry(hub, keyring):
    store, env = keyring
    seed(store, "adhoc", "a@example.com", "v1")
    assert secret(hub, env, "get", "adhoc").stdout == "v1"
    seed(store, "adhoc", "b@example.com", "v2")
    out = secret(hub, env, "get", "adhoc")
    assert out.returncode == 1 and "pass --user" in out.stderr
    assert secret(hub, env, "get", "adhoc", "--user", "b@example.com").stdout == "v2"


def test_get_exits_2_when_missing(hub, keyring):
    _, env = keyring
    out = secret(hub, env, "get", "sonar")
    assert out.returncode == 2
    assert "crow secret set sonar" in out.stderr


def test_list_shows_labels_dates_expiry_never_values(hub, keyring):
    store, env = keyring
    seed(store, "sonar", "me@example.com", "s3cret", label="crow sonar token")
    out = secret(hub, env, "list")
    assert out.returncode == 0
    assert "s3cret" not in out.stdout + out.stderr
    rows = {l.split()[0]: l for l in out.stdout.splitlines()[1:]}
    assert "crow sonar token" in rows["sonar"] and "2026-11-03" in rows["sonar"]
    assert "EXPIRED" in rows["sonar"]
    assert "NOT STORED" in rows["jira"]


def test_rm_clears_only_the_resolved_user(hub, keyring):
    store, env = keyring
    seed(store, "sonar", "me@example.com", "a")
    seed(store, "sonar", "other@example.com", "b")
    out = secret(hub, env, "rm", "sonar")
    assert out.returncode == 0
    assert [i["attrs"]["user"] for i in json.loads(store.read_text())] == ["other@example.com"]
    assert secret(hub, env, "rm", "sonar").returncode == 2


def test_crow_dispatches_secret(hub, keyring):
    _, env = keyring
    out = run(hub, "crow", "secret", "list", env=env)
    assert out.returncode == 0
    assert out.stdout.startswith("SERVICE")
