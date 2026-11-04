"""crow-usage by_model: each model gets its own turns, matching a transcript recount (item_029 AC2)."""
from test_crow_usage_projects import setup, turn, usage


def recount(events):
    seen, tot = set(), {}
    for e in events:
        m = e["message"]
        if m["id"] in seen:
            continue
        seen.add(m["id"])
        name = m["model"].replace("claude-", "")
        tot[name] = tot.get(name, 0) + m["usage"]["input_tokens"] + m["usage"]["output_tokens"]
    return tot


def test_multi_model_session_is_split_by_turn(hub):
    events = [turn("a", 0, out=9000, model="claude-opus-5-5"),
              turn("b", 0, out=500, model="claude-haiku-5-5"),
              turn("c", 0, out=3000, model="claude-opus-5-5"),
              turn("c", 0, out=3000, model="claude-opus-5-5")]  # streamed duplicate
    setup(hub, "m1", events)
    models = {e["name"]: e for e in usage(hub)["today"]["by_model"]}
    expected = recount(events)
    assert set(models) == set(expected)
    for name, tokens in expected.items():
        assert abs(models[name]["tokens"] - tokens) <= 0.02 * tokens
    assert models["haiku-5-5"]["tokens"] == 500  # not the whole session
