"""Parser for the run contract: spec examples pass, broken copies fail with clear messages."""
import json
import re
import subprocess
import sys

import pytest
import yaml
from conftest import BIN, REPO

sys.path.insert(0, str(BIN))
import crow_contract as cc

SPEC = REPO / "logics" / "specs" / "req_037_run_contract"
EX = SPEC / "examples"
PHASE_FILES = sorted(p for p in EX.rglob("*.md") if cc.file_name_parts(p))
FENCE = re.compile(r"^```crow\n(.*?)^```", re.DOTALL | re.MULTILINE)


def block(path):
    return yaml.safe_load(FENCE.search(path.read_text()).group(1))


def write(tmp_path, name, data, extra=""):
    p = tmp_path / name
    p.write_text("# t\n\n```crow\n" + yaml.safe_dump(data, sort_keys=False) + "```\n" + extra)
    return p


def test_examples_found():
    assert len(PHASE_FILES) >= 7


@pytest.mark.parametrize("path", PHASE_FILES, ids=lambda p: p.name)
def test_spec_examples_parse(path):
    res = cc.parse_report_file(path)
    assert res.ok, res.errors
    assert res.data["contract"] == 1


def test_spec_questions_validate():
    res = cc.parse_questions_file(EX / "gates" / "questions.json")
    assert res.ok, res.errors


def test_ready_report_without_proof_fails(tmp_path):
    d = block(EX / "gates" / "report-gate.report.v1.md")
    del d["proof"]
    res = cc.parse_report_file(write(tmp_path, "report.v1.md", d))
    assert any("proof" in e for e in res.errors)


def test_v2_without_responses_fails(tmp_path):
    d = block(EX / "20261109-build-item-095" / "plan.v2.md")
    del d["responses"]
    res = cc.parse_report_file(write(tmp_path, "plan.v2.md", d))
    assert any("responses" in e for e in res.errors)


def test_two_blocks_fail(tmp_path):
    d = block(EX / "gates" / "report-gate.report.v1.md")
    p = write(tmp_path, "report.v1.md", d, "\n```crow\ncontract: 1\n```\n")
    res = cc.parse_report_file(p)
    assert res.errors == ["found 2 ```crow blocks, exactly one is required"]


def test_zero_blocks_fail(tmp_path):
    p = tmp_path / "report.v1.md"
    p.write_text("# no block\n")
    assert "no ```crow block" in cc.parse_report_file(p).errors[0]


def test_unclosed_block_fails(tmp_path):
    p = tmp_path / "report.v1.md"
    p.write_text("```crow\ncontract: 1\n")
    assert "not closed" in cc.parse_report_file(p).errors[0]


def test_bad_yaml_fails(tmp_path):
    p = tmp_path / "report.v1.md"
    p.write_text("```crow\na: [1, 2\n```\n")
    assert "not valid YAML or JSON" in cc.parse_report_file(p).errors[0]


def test_json_block_accepted(tmp_path):
    d = block(EX / "gates" / "report-gate.report.v1.md")
    p = tmp_path / "report.v1.md"
    p.write_text("```crow\n" + json.dumps(d) + "\n```\n")
    assert cc.parse_report_file(p).ok


def test_name_must_match_version(tmp_path):
    d = block(EX / "gates" / "report-gate.report.v1.md")
    res = cc.parse_report_file(write(tmp_path, "report.v2.md", d))
    assert any("file name says version 2, block says 1" in e for e in res.errors)


def test_name_must_match_phase(tmp_path):
    d = block(EX / "gates" / "report-gate.report.v1.md")
    res = cc.parse_report_file(write(tmp_path, "plan.v1.md", d))
    assert any("file name says phase 'plan'" in e for e in res.errors)


def test_bad_file_name(tmp_path):
    d = block(EX / "gates" / "report-gate.report.v1.md")
    res = cc.parse_report_file(write(tmp_path, "notes.md", d))
    assert any("not <phase>.v<N>.md" in e for e in res.errors)


def test_question_without_blocking_fails(tmp_path):
    d = json.loads((EX / "gates" / "questions.json").read_text())
    del d["questions"][0]["blocking"]
    p = tmp_path / "questions.json"
    p.write_text(json.dumps(d))
    res = cc.parse_questions_file(p)
    assert any("blocking" in e for e in res.errors)


def test_questions_bad_json(tmp_path):
    p = tmp_path / "questions.json"
    p.write_text("{")
    assert "not valid JSON" in cc.parse_questions_file(p).errors[0]


def test_missing_file(tmp_path):
    assert "cannot read" in cc.parse_report_file(tmp_path / "report.v1.md").errors[0]


def test_cli_exit_codes(tmp_path):
    good = EX / "gates" / "report-gate.report.v1.md"
    bad = tmp_path / "report.v1.md"
    bad.write_text("nothing\n")
    ok = subprocess.run([sys.executable, str(BIN / "crow-contract"), str(good)], capture_output=True, text=True, check=False)
    assert ok.returncode == 0 and ok.stdout.startswith("ok ")
    ko = subprocess.run([sys.executable, str(BIN / "crow-contract"), str(bad)], capture_output=True, text=True, check=False)
    assert ko.returncode == 1 and "no ```crow block" in ko.stdout
