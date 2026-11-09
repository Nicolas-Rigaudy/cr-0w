"""Shared helpers for the vault capture scripts (crow-note, crow-notes)."""
import datetime as dt
import os
import re
from pathlib import Path

import crow_lib
import yaml

ROOT = Path(crow_lib.data_home())
REGISTRY = ROOT / "local" / "projects.yaml"

CAPTURE = "📥 Capture"
MEETINGS = "📆 Meetings & notes"
EOD = "🌙 End of day"
CROW = "CR-0W"  # crow's activity section in the daily note
SECTIONS = {"capture": CAPTURE, "meetings": MEETINGS, "eod": EOD}

DAILY_DIR = "40 Journal/Daily"
TEMPLATE = "90 Templates/Daily Note.md"
INBOX = "00 Inbox"
WORK = "10 Work"


def vault_path():
    """Vault root: CROW_VAULT overrides the registry's vault.path."""
    env = os.environ.get("CROW_VAULT")
    if env:
        return Path(env)
    try:
        data = yaml.safe_load(REGISTRY.read_text()) or {}
        path = (data.get("vault") or {}).get("path")
    except OSError:
        path = None
    return Path(path) if path else None


def daily_path(vault, day):
    return vault / DAILY_DIR / f"{day.isoformat()}.md"


def _tp_date(match):
    """Resolve a Templater tp.date.now(fmt, offset, ref, ref_fmt) against the file date."""
    args = [a.strip().strip('"') for a in match.group(1).split(",")]
    fmt, offset, ref = args[0], int(args[1]) if len(args) > 1 else 0, args[2]
    day = dt.date.fromisoformat(ref) + dt.timedelta(days=offset)
    if fmt == "YYYY-MM-DD":
        return day.isoformat()
    names = ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"]
    months = ["January", "February", "March", "April", "May", "June", "July",
              "August", "September", "October", "November", "December"]
    return f"{names[day.weekday()]} {day.day} {months[day.month - 1]} {day.year}"


def render_template(text, day):
    """Fill the Templater fields the daily template uses (dates, title, cursor)."""
    text = text.replace("tp.file.title", f'"{day.isoformat()}"')
    text = re.sub(r"<%\s*tp\.date\.now\((.*?)\)\s*%>", _tp_date, text)
    return re.sub(r"<%\s*tp\.file\.cursor\(\)\s*%>", "", text)


def section_bounds(lines, heading):
    """(start, end) line indexes of the body under '## <heading>', or None."""
    start = None
    for i, line in enumerate(lines):
        if line.startswith("## "):
            if start is not None:
                return start, i
            if line[3:].strip().startswith(heading):
                start = i + 1
    return (start, len(lines)) if start is not None else None


def known_template(text):
    """True when every Templater call is one render_template fills (dates, title, cursor)."""
    calls = re.findall(r"<%(.*?)%>", text, re.DOTALL)
    return all(re.fullmatch(r"\s*tp\.(date\.now\(.*\)|file\.cursor\(\))\s*", c, re.DOTALL) for c in calls)


def crow_lines(vault, day):
    """Activity lines under '## CR-0W' in that day's daily note ([] when none)."""
    try:
        lines = daily_path(vault, day).read_text().split("\n")
    except OSError:
        return []
    bounds = section_bounds(lines, CROW)
    if bounds is None:
        return []
    return [l for l in lines[bounds[0]:bounds[1]] if l.startswith("- ")]
