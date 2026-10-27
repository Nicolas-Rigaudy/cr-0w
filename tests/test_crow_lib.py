import sys

sys.path.insert(0, str(__import__("pathlib").Path(__file__).resolve().parent.parent / "bin"))
import crow_lib  # noqa: E402


def test_profiles_dir_prefers_cdx_profiles():
    assert crow_lib.profiles_dir({"CDX_PROFILES": "/p", "CDX_HOME": "/h"}) == "/p"


def test_profiles_dir_uses_cdx_home():
    assert crow_lib.profiles_dir({"CDX_HOME": "/h"}) == "/h/profiles"


def test_profiles_dir_defaults_to_dot_cdx():
    assert crow_lib.profiles_dir({}) == crow_lib.real_home() + "/.cdx/profiles"
