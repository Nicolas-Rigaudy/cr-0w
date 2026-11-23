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
NOTES_DAILY_DIR = "daily"  # dailies under CROW_HOME/notes while the vault brick is off
NOTES_DIR = "notes"
# Daily note of CROW_HOME/notes (vault off): the sections crow-note and crow-log write to.
NOTES_DAILY = "# {day}\n\n## 📥 Capture\n- \n\n## 📆 Meetings & notes\n- \n\n## 🌙 End of day\n- \n"
TEMPLATE = "90 Templates/Daily Note.md"
INBOX = "00 Inbox"
WORK = "10 Work"


def notes_mode():
    """True while the vault brick is off: notes live as plain Markdown under CROW_HOME/notes."""
    return not crow_lib.brick_on("vault")


def vault_path():
    """Vault root: CROW_VAULT overrides the registry's vault.path. Vault brick off: CROW_HOME/notes."""
    if notes_mode():
        return Path(crow_lib.data_home()) / NOTES_DIR
    env = os.environ.get("CROW_VAULT")
    if env:
        return Path(env)
    try:
        data = yaml.safe_load(REGISTRY.read_text()) or {}
        path = (data.get("vault") or {}).get("path")
    except OSError:
        path = None
    return Path(path) if path else None


def daily_dir(vault):
    """Folder of the daily notes: the vault journal, or CROW_HOME/notes/daily while the vault brick is off."""
    return vault / (NOTES_DAILY_DIR if notes_mode() else DAILY_DIR)


def daily_path(vault, day):
    return daily_dir(vault) / f"{day.isoformat()}.md"


def new_daily(vault, day, strict=False):
    """Text of a missing daily: the vault template rendered for day, else the plain built-in daily while the
    vault brick is off, else None. strict refuses a template with Templater calls crow cannot fill."""
    tpl = vault / TEMPLATE
    if tpl.exists():
        raw = tpl.read_text()
        return None if strict and not known_template(raw) else render_template(raw, day)
    return NOTES_DAILY.format(day=day.isoformat()) if notes_mode() else None


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
