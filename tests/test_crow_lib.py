import sys

sys.path.insert(0, str(__import__("pathlib").Path(__file__).resolve().parent.parent / "bin"))
import crow_lib  # noqa: E402


def test_profiles_dir_prefers_cdx_profiles():
    assert crow_lib.profiles_dir({"CDX_PROFILES": "/p", "CDX_HOME": "/h"}) == "/p"


def test_profiles_dir_uses_cdx_home():
    assert crow_lib.profiles_dir({"CDX_HOME": "/h"}) == "/h/profiles"


def test_profiles_dir_defaults_to_dot_cdx():
    assert crow_lib.profiles_dir({}) == crow_lib.real_home() + "/.cdx/profiles"


FIX = __import__("pathlib").Path(__file__).parent / "fixtures"


def _read(name):
    return (FIX / name).read_text()


def test_final_text_claude_result_envelope():
    assert crow_lib.final_text(_read("run-claude.jsonl"), "claude") == "Shipped: claude report"


def test_final_text_claude_single_json_document():
    assert crow_lib.final_text('{"type":"result","result":"one doc"}', "claude") == "one doc"


def test_final_text_codex_last_agent_message():
    assert crow_lib.final_text(_read("run-codex.jsonl"), "codex") == "Shipped: codex report"


def test_final_text_codex_legacy_shape():
    assert crow_lib.final_text(_read("run-codex-legacy.jsonl"), "codex") == "old shape report"


def test_final_text_unknown_output_falls_back_to_tail():
    assert crow_lib.final_text("plain output", "codex") == "plain output"


ROWS = [{"session_name": "main", "provider": "claude"}, {"session_name": "cx", "provider": "codex"}]


def test_provider_of_reads_cdx_rows_and_defaults_to_claude():
    assert crow_lib.provider_of("cx", ROWS) == "codex"
    assert crow_lib.provider_of("main", ROWS) == "claude"
    assert crow_lib.provider_of("unknown", ROWS) == "claude"


def test_check_model_refuses_the_other_providers_names():
    assert crow_lib.check_model("codex", "sonnet")
    assert crow_lib.check_model("codex", "claude-opus-5-5")
    assert crow_lib.check_model("claude", "gpt-5.1-codex")
    assert crow_lib.check_model("codex", "gpt-5.1-codex") is None
    assert crow_lib.check_model("claude", "haiku") is None
    assert crow_lib.check_model("codex", None) is None


def _hub(tmp_path):
    for f in ("CLAUDE.md", "AGENTS.md", "LOGICS.md", "local/owner-prefs.md"):
        p = tmp_path / f
        p.parent.mkdir(exist_ok=True)
        p.write_text("x")
    s = tmp_path / ".claude/skills/brief"
    s.mkdir(parents=True)
    (s / "SKILL.md").write_text("skill")
    return str(tmp_path)


def test_oneshot_prompt_claude_is_the_slash_command(tmp_path):
    assert crow_lib.oneshot_prompt("claude", _hub(tmp_path), "daily", "Draft only") == "/daily Draft only"


def test_oneshot_prompt_codex_reads_instructions_then_the_skill(tmp_path):
    hub = _hub(tmp_path)
    p = crow_lib.oneshot_prompt("codex", hub, "brief")
    assert p.startswith("Read these files first: ")
    for f in ("CLAUDE.md", "AGENTS.md", "LOGICS.md", "local/owner-prefs.md"):
        assert f"{hub}/{f}" in p
    assert f"{hub}/.claude/skills/brief/SKILL.md" in p and "`brief` skill" in p
    assert "local/context.md" not in p  # only files that exist


def test_oneshot_prompt_codex_without_skill_file_exits(tmp_path):
    import pytest
    with pytest.raises(SystemExit):
        crow_lib.oneshot_prompt("codex", _hub(tmp_path), "weekly")


def test_permission_and_hub_session_guard_per_provider():
    assert crow_lib.permission("claude") == "auto"
    assert crow_lib.permission("codex") == "default"
    assert crow_lib.hub_session_refusal("claude") is None
    assert "guard" in crow_lib.hub_session_refusal("codex")
