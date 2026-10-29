"""crow-tidy: a freshly cut branch is not treated as merged."""
import importlib.machinery
import importlib.util
import subprocess

from conftest import BIN

loader = importlib.machinery.SourceFileLoader("crow_tidy", str(BIN / "crow-tidy"))
spec = importlib.util.spec_from_loader("crow_tidy", loader)
ct = importlib.util.module_from_spec(spec)
loader.exec_module(ct)


def sh(repo, *args):
    subprocess.run(["git", "-C", str(repo), "-c", "user.name=t", "-c", "user.email=t@t", *args],
                   check=True, capture_output=True)


def test_merged_needs_commits_of_its_own(tmp_path):
    sh(tmp_path, "init", "-b", "main")
    sh(tmp_path, "commit", "--allow-empty", "-m", "base")
    sh(tmp_path, "checkout", "-b", "done")
    sh(tmp_path, "commit", "--allow-empty", "-m", "work")
    sh(tmp_path, "checkout", "main")
    sh(tmp_path, "merge", "--no-ff", "-m", "merge", "done")
    sh(tmp_path, "branch", "fresh")
    assert not ct.merged(tmp_path, "fresh", "main")
    assert ct.merged(tmp_path, "done", "main")
