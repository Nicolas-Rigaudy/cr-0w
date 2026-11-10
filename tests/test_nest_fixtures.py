"""Nest extra-state fixtures (item_137) match the core contract schemas."""
import json

import pytest
from conftest import REPO
from schema_lite import Validator

SCHEMAS = REPO / "logics" / "specs" / "core_contract" / "schemas"
STATES = REPO / "tools" / "nest" / "fixtures" / "states"


def load(p):
    return json.loads(p.read_text())


def errors(schema, instance):
    return Validator(load(SCHEMAS / schema), SCHEMAS).errors(instance)


@pytest.mark.parametrize("name,schema", [("missions.json", "mission.schema.json"), ("inbox.json", "inbox.schema.json")])
def test_states_valid(name, schema):
    for record in load(STATES / name):
        assert not errors(schema, record), record["id"]


def test_states_cover_every_ui_state():
    missions = load(STATES / "missions.json")
    states = {m["state"] for m in missions}
    assert {"failed", "queued", "waiting", "review"} <= states
    assert any(m["state"] == "queued" and m["runs"] and m["runs"][-1]["cause"] == "usage_limit" for m in missions)


def test_inbox_rows_point_at_known_missions():
    ids = {m["id"] for m in load(STATES / "missions.json")}
    assert all(row["mission"] in ids for row in load(STATES / "inbox.json"))
