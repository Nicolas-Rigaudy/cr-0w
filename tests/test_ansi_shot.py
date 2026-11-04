import importlib.machinery
import importlib.util
import os
import sys

import pytest

from conftest import BIN

pytest.importorskip("PIL")
pytest.importorskip("fontTools")
needs_font = pytest.mark.skipif(not os.path.exists("/usr/share/fonts/TTF/FiraCodeNerdFontMono-Regular.ttf"),
                                reason="default font missing")


def _load():
    loader = importlib.machinery.SourceFileLoader("ansi_shot", str(BIN / "ansi-shot"))
    spec = importlib.util.spec_from_loader("ansi_shot", loader)
    mod = importlib.util.module_from_spec(spec)
    loader.exec_module(mod)
    return mod


def _render(mod, *argv):
    old = sys.argv
    sys.argv = ["ansi-shot", *argv]
    try:
        mod.main()
    finally:
        sys.argv = old


@needs_font
def test_reverse_video_on_default_background_renders(tmp_path):
    mod = _load()
    src = tmp_path / "in.ansi"
    src.write_text("\x1b[7mrev\x1b[0m plain\n")
    out = tmp_path / "out.png"
    _render(mod, str(src), str(out))
    assert out.exists()


def test_bundled_theme_when_no_kitty_theme(tmp_path, monkeypatch):
    mod = _load()
    monkeypatch.setattr(mod, "KITTY_THEME", str(tmp_path / "missing.conf"))
    assert mod.theme_path(None) == mod.DEFAULT_THEME and mod.theme_path("default") == mod.DEFAULT_THEME
    th = mod.kitty_theme(mod.DEFAULT_THEME)
    assert all(f"color{i}" in th for i in range(16)) and "foreground" in th and "background" in th
    kitty = tmp_path / "kitty.conf"
    kitty.write_text("color1 #ff0000\n")
    monkeypatch.setattr(mod, "KITTY_THEME", str(kitty))
    assert mod.theme_path(None) == str(kitty)
    with pytest.raises(SystemExit, match="not found"):
        mod.theme_path(str(tmp_path / "nope.conf"))


def test_renders_with_the_bundled_theme_and_any_installed_font(tmp_path, monkeypatch):
    mod = _load()
    if not any(os.path.exists(f) for f in mod.FALLBACKS):
        pytest.skip("no fallback font installed")
    monkeypatch.setattr(mod, "KITTY_THEME", str(tmp_path / "missing.conf"))
    src = tmp_path / "in.ansi"
    src.write_text("\x1b[31mred\x1b[0m plain\n")
    out = tmp_path / "out.png"
    _render(mod, str(src), str(out), "--font", str(tmp_path / "no-such-font.ttf"))
    from PIL import Image
    img = Image.open(out).convert("RGB")
    assert img.getpixel((1, 1)) == (0, 0, 0)  # bundled background #000000
