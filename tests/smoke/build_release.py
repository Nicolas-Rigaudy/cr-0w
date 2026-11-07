"""Fake release (0.1.0, 0.2.0), API answers and a data repo for the container smoke test."""
import hashlib
import json
import subprocess
import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parent.parent))
from test_install_update import tarball

out = Path(sys.argv[1])
for v in ("0.1.0", "0.2.0"):
    d = out / "rel/download" / f"v{v}"; d.mkdir(parents=True, exist_ok=True)
    b = tarball(v, {f"changelogs/CHANGELOGS_{v.replace('.','_')}.md": f"# {v}\n- smoke notes {v}\n"})
    (d / f"crow-{v}.tar.gz").write_bytes(b)
    (d / "SHA256SUMS").write_text(f"{hashlib.sha256(b).hexdigest()}  crow-{v}.tar.gz\n")
(out / "api/releases/tags").mkdir(parents=True, exist_ok=True)
(out / "api/releases/latest").write_text(json.dumps({"tag_name": "v0.1.0", "body": "first"}))
(out / "api/releases/tags/v0.2.0").write_text(json.dumps({"tag_name": "v0.2.0", "body": "- smoke notes 0.2.0"}))
d = out / "data"; d.mkdir(exist_ok=True)
subprocess.run(["git", "init", "-q", "-b", "main", str(d)], check=True)
(d / "tokens.yaml").write_text("tokens:\n  - service: forge\n    user: a@example.com\n")
(d / "accounts.yaml").write_text("default: main\naccounts:\n  - {name: main, policy: use}\n")
subprocess.run(["git", "-C", str(d), "add", "-A"], check=True)
subprocess.run(["git", "-C", str(d), "-c", "user.name=t", "-c", "user.email=t@example.com", "commit", "-q", "-m", "d"], check=True)
