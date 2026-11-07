"""The wrap and retro skills must tell the agent to keep private specifics out of tracked runbooks."""
from pathlib import Path

ROOT = Path(__file__).resolve().parent.parent


def test_wrap_and_retro_route_specifics_to_local_runbooks():
    for rel in ("skills/wrap/SKILL.md", ".claude/skills/retro/SKILL.md"):
        text = (ROOT / rel).read_text()
        assert "local/runbooks/" in text, rel
        assert "runbooks/" in text and "tracked" in text, rel
