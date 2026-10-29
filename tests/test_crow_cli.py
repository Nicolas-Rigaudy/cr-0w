"""crow: the help shows five groups and old command names keep working with a hint."""
from conftest import run


def test_help_shows_the_five_groups(hub):
    out = run(hub, "crow", "--help")
    assert out.returncode == 0
    lines = out.stdout.splitlines()
    groups = [l for l in lines if l in ("Start", "See", "Delegate", "Rituals", "Data")]
    assert groups == ["Start", "See", "Delegate", "Rituals", "Data"]
    for name in ("crow projects", "crow doctor", "crow runs sessions", "crow runs relaunch"):
        assert name in out.stdout
    assert "set -euo" not in out.stdout


def test_old_names_print_a_hint_on_stderr(hub):
    for old, new in (("sessions", "crow runs sessions"), ("relaunch", "crow runs relaunch"),
                     ("releases", "crow projects")):
        out = run(hub, "crow", old, "--help")
        assert f"'crow {old}' is now '{new}'" in out.stderr
        assert "is now" not in out.stdout


def test_new_names_print_no_hint(hub):
    for args in (("runs", "sessions", "--help"), ("projects",)):
        out = run(hub, "crow", *args)
        assert "is now" not in out.stderr
