"""crow_lib: helpers shared by the crow scripts."""
import os
import pwd


def real_home():
    return pwd.getpwuid(os.getuid()).pw_dir


def profiles_dir(env=None):
    """cdx profiles dir: CDX_PROFILES, else CDX_HOME/profiles, else ~/.cdx/profiles."""
    env = os.environ if env is None else env
    if env.get("CDX_PROFILES"):
        return os.path.expanduser(env["CDX_PROFILES"])
    root = env.get("CDX_HOME") or os.path.join(real_home(), ".cdx")
    return os.path.join(os.path.expanduser(root), "profiles")


def _json_objects(text):
    """Objects from a JSON document, a JSON array or JSON lines; other lines are skipped."""
    import json
    try:
        doc = json.loads(text)
        docs = doc if isinstance(doc, list) else [doc]
    except ValueError:
        docs = []
        for line in text.splitlines():
            try:
                docs.append(json.loads(line))
            except ValueError:
                continue
    return [d for d in docs if isinstance(d, dict)]


def final_text(text, provider="claude"):
    """Final message of a cdx run's stdout. Claude: the `result` envelope. Codex: the last
    agent_message event (`item.completed` or the older `msg` shape). Other output: its tail."""
    objs = _json_objects(text)
    found = None
    for o in objs:
        if provider == "claude" and o.get("type") == "result":
            found = o.get("result")
        elif provider == "codex":
            item = o.get("item") if o.get("type") == "item.completed" else o.get("msg")
            if isinstance(item, dict) and item.get("type") == "agent_message":
                found = item.get("text") or item.get("message")
    return found if found else text[-4000:]
