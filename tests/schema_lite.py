"""Test-side alias: the validator lives in bin/crow_schema_lite.py."""
import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parent.parent / "bin"))
from crow_schema_lite import *
from crow_schema_lite import ANNOTATIONS, APPLIED, UnsupportedKeyword, Validator  # noqa: F401
