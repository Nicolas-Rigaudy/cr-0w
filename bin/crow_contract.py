"""Parse and validate the run contract files (spec: logics/specs/req_037_run_contract).

parse_report_file(path)   -> ContractResult for a plan.vN.md, report.vN.md or mockup.vN.md
parse_questions_file(path) -> ContractResult for outbox/questions.json
Every problem is a plain message in result.errors; nothing raises on bad input.
"""
import json
import re
import sys
from dataclasses import dataclass, field
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parent))
import yaml
from crow_schema_lite import Validator

SCHEMAS = Path(__file__).resolve().parent.parent / "logics" / "specs" / "req_037_run_contract" / "schemas"
FENCE_OPEN = re.compile(r"^```crow[ \t]*$")
FENCE_CLOSE = re.compile(r"^```[ \t]*$")
NAME = re.compile(r"^(plan|report|mockup)\.v([0-9]+)\.md$")


@dataclass
class ContractResult:
    path: str
    data: dict | None = None
    errors: list = field(default_factory=list)

    @property
    def ok(self):
        return not self.errors


def _validator(name):
    return Validator(json.loads((SCHEMAS / name).read_text()), SCHEMAS)


def extract_blocks(text):
    """Bodies of every fenced ```crow block, in order. An unclosed fence is reported as None."""
    blocks, cur = [], None
    for line in text.splitlines():
        if cur is None:
            if FENCE_OPEN.match(line):
                cur = []
        elif FENCE_CLOSE.match(line):
            blocks.append("\n".join(cur))
            cur = None
        else:
            cur.append(line)
    if cur is not None:
        blocks.append(None)
    return blocks


def _load_body(body):
    """YAML is a superset of JSON, so one loader reads both."""
    return yaml.safe_load(body)


def file_name_parts(path):
    """(phase, version) from the file name. Examples use <label>.<phase>.vN.md too."""
    name = Path(path).name
    m = NAME.match(name)
    if not m:
        m = re.search(r"(?:^|\.)(plan|report|mockup)\.v([0-9]+)\.md$", name)
    return (m.group(1), int(m.group(2))) if m else None


def parse_report_text(text, path="<text>"):
    res = ContractResult(str(path))
    blocks = extract_blocks(text)
    if not blocks:
        res.errors.append("no ```crow block found (exactly one is required)")
        return res
    if len(blocks) > 1:
        res.errors.append(f"found {len(blocks)} ```crow blocks, exactly one is required")
        return res
    if blocks[0] is None:
        res.errors.append("the ```crow block is not closed")
        return res
    try:
        data = _load_body(blocks[0])
    except yaml.YAMLError as e:
        res.errors.append(f"crow block is not valid YAML or JSON: {e}".replace("\n", " "))
        return res
    if not isinstance(data, dict):
        res.errors.append("crow block must be a mapping")
        return res
    res.data = data
    res.errors.extend(_validator("crow-block.schema.json").errors(data))
    return res


def parse_report_file(path):
    path = Path(path)
    try:
        text = path.read_text()
    except OSError as e:
        return ContractResult(str(path), errors=[f"cannot read file: {e}"])
    res = parse_report_text(text, path)
    parts = file_name_parts(path)
    if parts is None:
        res.errors.append(f"file name {path.name!r} is not <phase>.v<N>.md (phase plan, report or mockup)")
    elif isinstance(res.data, dict):
        phase, version = parts
        if res.data.get("phase") != phase:
            res.errors.append(f"file name says phase {phase!r}, block says {res.data.get('phase')!r}")
        if res.data.get("version") != version:
            res.errors.append(f"file name says version {version}, block says {res.data.get('version')!r}")
    return res


def parse_questions_file(path):
    path = Path(path)
    res = ContractResult(str(path))
    try:
        data = json.loads(path.read_text())
    except OSError as e:
        res.errors.append(f"cannot read file: {e}")
        return res
    except ValueError as e:
        res.errors.append(f"not valid JSON: {e}")
        return res
    res.data = data
    res.errors.extend(_validator("questions.schema.json").errors(data))
    return res


def parse_file(path):
    """Pick the parser from the file name: questions.json or a phase file."""
    return parse_questions_file(path) if Path(path).suffix == ".json" else parse_report_file(path)
