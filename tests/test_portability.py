"""Scripts must not assume a user name or a hub-private file."""
import re

from conftest import BIN

HARD_HOME = re.compile(r"""["'/]home/(?!me\b|\$|<|\*|user\b)[a-z]\w*""")


def test_no_hard_coded_home_dir():
    bad = [f"{p.name}: {m.group(0)}" for p in BIN.iterdir() if p.is_file() and p.suffix != ".pyc"
           for m in HARD_HOME.finditer(p.read_text(errors="replace"))]
    assert bad == []
