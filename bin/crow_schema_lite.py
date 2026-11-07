"""Stdlib validator for the JSON Schema subset the run contract spec uses.

Supported: type, enum, const, required, properties, additionalProperties, pattern, minimum,
maximum, minLength, maxLength, items, minItems, uniqueItems, $ref (local pointer or sibling
file) with $defs, allOf, if/then, propertyNames. Annotations are ignored. Any other keyword
raises UnsupportedKeyword, so a schema edit cannot pass silently.
"""
import json
import re
from pathlib import Path

ANNOTATIONS = {"$schema", "$id", "title", "description", "$comment", "default", "examples", "format", "$defs"}
APPLIED = {"type", "enum", "const", "required", "properties", "additionalProperties", "pattern",
           "minimum", "maximum", "minLength", "maxLength", "items", "minItems", "uniqueItems",
           "$ref", "allOf", "if", "then", "propertyNames"}

TYPES = {
    "object": lambda v: isinstance(v, dict),
    "array": lambda v: isinstance(v, list),
    "string": lambda v: isinstance(v, str),
    "boolean": lambda v: isinstance(v, bool),
    "null": lambda v: v is None,
    "integer": lambda v: isinstance(v, int) and not isinstance(v, bool),
    "number": lambda v: isinstance(v, (int, float)) and not isinstance(v, bool),
}


class UnsupportedKeyword(Exception):
    pass


def _canon(v):
    return json.dumps(v, sort_keys=True)


class Validator:
    def __init__(self, schema, base_dir=None):
        self.root = schema
        self.base_dir = Path(base_dir) if base_dir else None
        self._files = {}

    def errors(self, instance):
        """List of 'json.path: message' strings; empty when valid."""
        out = []
        self._check(self.root, self.root, instance, "$", out)
        return out

    def _resolve(self, ref, root):
        file_part, _, pointer = ref.partition("#")
        if file_part:
            if self.base_dir is None:
                raise UnsupportedKeyword(f"$ref to a file needs base_dir: {ref}")
            path = self.base_dir / file_part
            if path not in self._files:
                self._files[path] = json.loads(path.read_text())
            root = self._files[path]
        node = root
        for part in [p for p in pointer.split("/") if p]:
            node = node[part.replace("~1", "/").replace("~0", "~")]
        return node, root

    def _check(self, schema, root, inst, path, out):
        if schema is True:
            return
        if not isinstance(schema, dict):
            raise UnsupportedKeyword(f"{path}: schema is not an object")
        for key in schema:
            if key not in APPLIED and key not in ANNOTATIONS:
                raise UnsupportedKeyword(f"unsupported keyword {key!r}")
        if "$ref" in schema:
            target, troot = self._resolve(schema["$ref"], root)
            self._check(target, troot, inst, path, out)
        if "type" in schema:
            want = schema["type"] if isinstance(schema["type"], list) else [schema["type"]]
            for t in want:
                if t not in TYPES:
                    raise UnsupportedKeyword(f"unsupported type {t!r}")
            if not any(TYPES[t](inst) for t in want):
                out.append(f"{path}: expected {'|'.join(want)}, got {type(inst).__name__}")
                return
        if "enum" in schema and not any(_canon(inst) == _canon(e) for e in schema["enum"]):
            out.append(f"{path}: {inst!r} not in {schema['enum']}")
        if "const" in schema and _canon(inst) != _canon(schema["const"]):
            out.append(f"{path}: expected {schema['const']!r}, got {inst!r}")
        if TYPES["string"](inst):
            if "pattern" in schema and not re.search(schema["pattern"], inst):
                out.append(f"{path}: {inst!r} does not match {schema['pattern']}")
            if "minLength" in schema and len(inst) < schema["minLength"]:
                out.append(f"{path}: shorter than {schema['minLength']}")
            if "maxLength" in schema and len(inst) > schema["maxLength"]:
                out.append(f"{path}: longer than {schema['maxLength']}")
        if TYPES["number"](inst):
            if "minimum" in schema and inst < schema["minimum"]:
                out.append(f"{path}: {inst} below minimum {schema['minimum']}")
            if "maximum" in schema and inst > schema["maximum"]:
                out.append(f"{path}: {inst} above maximum {schema['maximum']}")
        if isinstance(inst, dict):
            self._object(schema, root, inst, path, out)
        if isinstance(inst, list):
            self._array(schema, root, inst, path, out)
        for sub in schema.get("allOf", []):
            self._check(sub, root, inst, path, out)
        if "if" in schema and "then" in schema:
            probe = []
            self._check(schema["if"], root, inst, path, probe)
            if not probe:
                self._check(schema["then"], root, inst, path, out)
        elif "if" in schema or "then" in schema:
            raise UnsupportedKeyword("if and then must come together")

    def _object(self, schema, root, inst, path, out):
        for name in schema.get("required", []):
            if name not in inst:
                out.append(f"{path}: missing required field {name!r}")
        props = schema.get("properties", {})
        extra = schema.get("additionalProperties", True)
        for name, value in inst.items():
            here = f"{path}.{name}"
            if name in props:
                self._check(props[name], root, value, here, out)
            elif extra is False:
                out.append(f"{here}: unknown field")
            elif isinstance(extra, dict):
                self._check(extra, root, value, here, out)
            if "propertyNames" in schema:
                self._check(schema["propertyNames"], root, name, f"{path}[{name!r}]", out)

    def _array(self, schema, root, inst, path, out):
        if "minItems" in schema and len(inst) < schema["minItems"]:
            out.append(f"{path}: fewer than {schema['minItems']} items")
        if schema.get("uniqueItems") and len({_canon(v) for v in inst}) != len(inst):
            out.append(f"{path}: duplicate items")
        if "items" in schema:
            for i, v in enumerate(inst):
                self._check(schema["items"], root, v, f"{path}[{i}]", out)
