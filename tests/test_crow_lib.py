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
