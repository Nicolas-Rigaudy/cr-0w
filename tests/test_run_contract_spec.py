"""The run contract spec (req_037) stays valid against its own schemas."""
import json
import re
from pathlib import Path

import pytest
import yaml

from conftest import REPO
from schema_lite import UnsupportedKeyword, Validator

SPEC = REPO / "logics" / "specs" / "req_037_run_contract"
SCHEMAS = SPEC / "schemas"
FENCE = re.compile(r"^```crow\n(.*?)^```", re.S | re.M)


def validator(name):
    path = SCHEMAS / name
    return Validator(json.loads(path.read_text()), SCHEMAS)


def rel(p):
    return str(p.relative_to(REPO))


def check(name, data, label):
    errs = validator(name).errors(data)
    assert not errs, f"{label}: " + "; ".join(errs)


KINDS = sorted((SPEC / "kinds").glob("*.yaml"))
BLOCK_FILES = sorted(SPEC.glob("examples/**/*.v[0-9]*.md"))
BLOCK_FILES = [p for p in BLOCK_FILES if p.name.startswith(("plan", "report", "mockup")) or ".plan." in p.name
               or ".report." in p.name or ".mockup." in p.name]


def test_eight_kind_files():
    assert len(KINDS) == 8


@pytest.mark.parametrize("path", KINDS, ids=lambda p: p.name)
def test_kind_file_valid(path):
    data = yaml.safe_load(path.read_text())
    assert data["name"] == path.stem
    check("kind.schema.json", data, rel(path))


@pytest.mark.parametrize("path", KINDS, ids=lambda p: p.name)
def test_ask_actions_are_published_actions(path):
    data = yaml.safe_load(path.read_text())
    publishes = [data["publish"]] + [o["publish"] for o in data.get("by_scope", {}).values() if "publish" in o]
    for pub in publishes:
        assert set(pub.get("ask_actions", [])) <= set(pub["actions"]), rel(path)


@pytest.mark.parametrize("path", BLOCK_FILES, ids=lambda p: p.name)
def test_crow_block_valid(path):
    m = FENCE.search(path.read_text())
    assert m, f"{rel(path)}: no crow block"
    check("crow-block.schema.json", yaml.safe_load(m.group(1)), rel(path))


def test_example_questions_valid():
    p = SPEC / "examples" / "gates" / "questions.json"
    check("questions.schema.json", json.loads(p.read_text()), rel(p))


def test_example_events_validation_lines_valid():
    p = SPEC / "examples" / "20261109-build-item-095" / "events.jsonl"
    lines = [json.loads(x) for x in p.read_text().splitlines() if x.strip()]
    vals = [x for x in lines if x.get("type") == "validation"]
    for v in vals:
        check("validation-event.schema.json", v, rel(p))


def test_every_schema_is_covered():
    assert {p.name for p in SCHEMAS.glob("*.json")} == {
        "kind.schema.json", "crow-block.schema.json", "questions.schema.json", "validation-event.schema.json"}


# negative cases

def build():
    return yaml.safe_load((SPEC / "kinds" / "build.yaml").read_text())


def test_rejects_bad_phase_with_path():
    d = build()
    d["phase"] = 9
    errs = validator("kind.schema.json").errors(d)
    assert any(e.startswith("$.phase:") for e in errs)


def test_rejects_duplicate_gates():
    d = build()
    d["gates"] = ["plan", "plan"]
    assert any(e.startswith("$.gates:") and "duplicate" in e for e in validator("kind.schema.json").errors(d))


def test_rejects_unknown_field():
    d = build()
    d["surprise"] = 1
    assert "$.surprise: unknown field" in validator("kind.schema.json").errors(d)


def test_rejects_unknown_nested_field_and_missing_required():
    d = build()
    del d["caps"]["budget_usd"]
    d["account"]["name"] = "x"
    errs = validator("kind.schema.json").errors(d)
    assert "$.account.name: unknown field" in errs
    assert any(e.startswith("$.caps:") and "budget_usd" in e for e in errs)


def test_conditional_requires_proof_on_ready_report():
    m = FENCE.search((SPEC / "examples" / "gates" / "report-gate.report.v1.md").read_text())
    d = yaml.safe_load(m.group(1))
    del d["proof"]
    assert any("proof" in e for e in validator("crow-block.schema.json").errors(d))


def test_unknown_keyword_fails():
    with pytest.raises(UnsupportedKeyword):
        Validator({"type": "object", "patternProperties": {}}).errors({})
    with pytest.raises(UnsupportedKeyword):
        Validator({"properties": {"a": {"oneOf": []}}}).errors({"a": 1})


def test_every_spec_schema_uses_only_the_subset():
    # a keyword outside the subset would raise while validating an instance that reaches it
    for p in SCHEMAS.glob("*.json"):
        walk(json.loads(p.read_text()))


def walk(node, in_props=False):
    from schema_lite import ANNOTATIONS, APPLIED
    if isinstance(node, dict):
        if not in_props:
            for k in node:
                assert k in APPLIED or k in ANNOTATIONS, f"unsupported keyword {k!r}"
        for k, v in node.items():
            if k in ("properties", "$defs"):
                for sub in v.values():
                    walk(sub)
            elif k in ("items", "additionalProperties", "propertyNames", "if", "then"):
                walk(v)
            elif k == "allOf":
                for sub in v:
                    walk(sub)
    return None
