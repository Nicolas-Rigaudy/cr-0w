#!/usr/bin/env python3
"""Base for the Python sidecars of crow-core (spec: logics/specs/core_contract).

A sidecar reads one JSON request per line on stdin and writes one JSON reply per line on
stdout. Logs go to stderr. Every sidecar answers hello, health and shutdown; a real sidecar
adds its own messages with register().

    crow_sidecar.py [--name base] [--version 0.0.0]
"""
import argparse
import json
import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parent))
from crow_schema_lite import Validator

CONTRACT = 1
SCHEMAS = Path(__file__).resolve().parent.parent / "logics" / "specs" / "core_contract" / "schemas"
ENVELOPE = ("v", "id", "type", "data")


class Sidecar:
    def __init__(self, name="base", version="0.0.0"):
        self.name = name
        self.version = version
        self.handlers = {}
        self.stopping = False
        self.register("hello", self._hello)
        self.register("health", lambda data: {"status": "ok"})
        self.register("shutdown", self._shutdown)

    def register(self, type_, fn):
        """fn(data) returns the reply data dict, or raises SidecarError."""
        self.handlers[type_] = fn

    @property
    def messages(self):
        return list(self.handlers)

    def _hello(self, data):
        if data["contract"] != CONTRACT:
            raise SidecarError("unsupported_contract",
                               f"core speaks contract {data['contract']}, this sidecar speaks {CONTRACT}")
        return {"name": self.name, "version": self.version, "contract": CONTRACT, "messages": self.messages}

    def _shutdown(self, data):
        self.stopping = True
        return {}

    def handle(self, msg):
        """Reply dict for one decoded request."""
        if not isinstance(msg, dict):
            return _error("?", "?", "bad_request", "request is not a JSON object")
        rid = msg.get("id") if isinstance(msg.get("id"), str) and msg.get("id") else "?"
        rtype = msg.get("type") if isinstance(msg.get("type"), str) and msg.get("type") else "?"
        missing = [k for k in ENVELOPE if k not in msg]
        extra = sorted(set(msg) - set(ENVELOPE))
        if missing or extra or not isinstance(msg["data"], dict):
            return _error(rid, rtype, "bad_request",
                          f"envelope: missing {missing}, unknown {extra}, data must be an object")
        if msg["v"] != CONTRACT:
            return _error(rid, rtype, "unsupported_contract", f"v {msg['v']} is not {CONTRACT}")
        fn = self.handlers.get(rtype)
        if fn is None:
            return _error(rid, rtype, "unknown_type", f"unknown message type: {rtype}")
        problems = _data_errors(rtype, msg["data"])
        if problems:
            return _error(rid, rtype, "bad_request", "; ".join(problems))
        try:
            return {"v": CONTRACT, "id": rid, "type": rtype, "ok": True, "data": fn(msg["data"])}
        except SidecarError as e:
            return _error(rid, rtype, e.code, e.message)
        except Exception as e:  # noqa: BLE001 - one bad request must not kill the sidecar
            print(f"{self.name}: {rtype} failed: {e!r}", file=sys.stderr)
            return _error(rid, rtype, "internal", f"{type(e).__name__}: {e}")

    def serve(self, stdin=sys.stdin, stdout=sys.stdout):
        """Answer requests until shutdown or end of input."""
        for line in stdin:
            if not line.strip():
                continue
            try:
                msg = json.loads(line)
            except json.JSONDecodeError as e:
                reply = _error("?", "?", "bad_request", f"not JSON: {e.msg}")
            else:
                reply = self.handle(msg)
            stdout.write(json.dumps(reply, separators=(",", ":")) + "\n")
            stdout.flush()
            if self.stopping:
                return


class SidecarError(Exception):
    def __init__(self, code, message):
        super().__init__(message)
        self.code = code
        self.message = message


def _error(rid, rtype, code, message):
    return {"v": CONTRACT, "id": rid, "type": rtype, "ok": False, "error": {"code": code, "message": message}}


_SCHEMA = None


def _data_errors(rtype, data):
    """Schema errors for the data of a known message; messages without a schema entry pass."""
    global _SCHEMA
    if _SCHEMA is None:
        _SCHEMA = json.loads((SCHEMAS / "request.schema.json").read_text())
    if rtype not in _SCHEMA["$defs"]:
        return []
    schema = {"$ref": f"request.schema.json#/$defs/{rtype}"}
    return [e.replace("$", "data", 1) for e in Validator(schema, SCHEMAS).errors(data)]


def main(argv=None):
    ap = argparse.ArgumentParser(description=__doc__.splitlines()[0])
    ap.add_argument("--name", default="base")
    ap.add_argument("--version", default="0.0.0")
    a = ap.parse_args(argv)
    Sidecar(a.name, a.version).serve()
    return 0


if __name__ == "__main__":
    sys.exit(main())
