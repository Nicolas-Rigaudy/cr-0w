"""crow-mrs: both forges mapped to rows, threads and diffs, from invented API payloads."""
import importlib.machinery
import importlib.util

from conftest import BIN, PYTHON, REPO  # noqa: F401


def load():
    loader = importlib.machinery.SourceFileLoader("crow_mrs", str(BIN / "crow-mrs"))
    spec = importlib.util.spec_from_loader("crow_mrs", loader)
    mod = importlib.util.module_from_spec(spec)
    loader.exec_module(mod)
    return mod


PROJS = [
    {"name": "demo-service", "aliases": ["demo"], "forge": "gitlab", "remote": "grp/demo-service"},
    {"name": "demo-tool", "forge": "github", "remote": "me/demo-tool"},
]


def gitlab_call(cmd):
    path = cmd[2]
    if path == "user":
        return {"username": "me"}
    if path.startswith("merge_requests?"):
        mine = {"project_id": 7, "iid": 5, "title": "typed ids", "author": {"username": "me"}, "reviewers": [{"username": "sam"}],
                "assignees": [], "draft": False, "updated_at": "2026-10-29T10:00:00Z", "created_at": "2026-10-28T10:00:00Z",
                "source_branch": "feat/ids", "target_branch": "main", "web_url": "https://x/5",
                "references": {"full": "grp/demo-service!5"}}
        theirs = {**mine, "iid": 6, "title": "retry policy", "author": {"username": "sam"}, "reviewers": [{"username": "me"}],
                  "references": {"full": "grp/demo-service!6"}}
        other = {**mine, "iid": 9, "references": {"full": "elsewhere/unrelated!9"}}
        if "created_by_me" in path:
            return [mine, other]
        return [theirs] if "reviewer_username" in path else []
    if path.endswith("/discussions?per_page=100"):
        note = lambda resolved: [{"resolvable": True, "resolved": resolved, "author": {"username": "sam"}, "body": "guard it",
                                  "created_at": "2026-10-29T09:00:00Z", "position": {"new_path": "a.py", "new_line": 3}}]
        return [{"notes": note(False)}, {"notes": note(True)}, {"notes": [{"resolvable": False}]}]
    return {"head_pipeline": {"status": "failed"}}


def test_gitlab_rows_roles_threads_and_unregistered_repos():
    rows = load().gitlab_rows(PROJS, call=gitlab_call)
    by_ref = {r["ref"]: r for r in rows}
    assert sorted(by_ref) == ["!5", "!6"]
    assert by_ref["!5"]["role"] == "AUTHOR" and by_ref["!5"]["project"] == "demo"
    assert (by_ref["!5"]["threads_open"], by_ref["!5"]["threads_resolved"]) == (1, 1)
    assert by_ref["!5"]["pipeline"] == "failed" and by_ref["!5"]["branch"] == "feat/ids"
    assert by_ref["!6"]["role"] == "REVIEWER" and by_ref["!6"]["threads_open"] == 0


def github_call(cmd):
    if cmd[:3] == ["gh", "api", "user"]:
        return {"login": "me"}
    if cmd[:3] == ["gh", "search", "prs"]:
        pr = {"repository": {"nameWithOwner": "me/demo-tool"}, "number": 12, "title": "pill", "author": {"login": "me"},
              "isDraft": True, "updatedAt": "2026-10-29T10:00:00Z", "createdAt": "2026-10-28T10:00:00Z", "url": "https://x/12"}
        return [pr] if "--author=@me" in cmd else []
    if cmd[:3] == ["gh", "pr", "view"]:
        return {"headRefName": "feat/pill", "baseRefName": "main", "statusCheckRollup": [{"conclusion": "SUCCESS"}, {"status": "IN_PROGRESS"}]}
    return {"data": {"repository": {"pullRequest": {"reviewThreads": {"nodes": [
        {"isResolved": False, "path": "b.go", "line": 4, "comments": {"nodes": [{"author": {"login": "kim"}, "body": "rename", "createdAt": "2026-10-29T08:00:00Z"}]}}]}}}}}


def test_github_rows_draft_pipeline_and_threads():
    (r,) = load().github_rows(PROJS, call=github_call)
    assert (r["ref"], r["role"], r["draft"], r["pipeline"]) == ("#12", "AUTHOR", True, "running")
    assert (r["threads_open"], r["branch"], r["project"]) == (1, "feat/pill", "demo-tool")


def test_one_forge_failing_keeps_the_other_and_names_the_error():
    mod = load()

    def broken(_projs):
        raise RuntimeError("401 unauthorized")
    out = mod.list_all(PROJS, gl=broken, gh=lambda p: mod.github_rows(p, call=github_call))
    assert out["forges"]["gitlab"] == "gitlab sync failed: 401 unauthorized" and out["forges"]["github"] == "ok"
    assert [r["ref"] for r in out["mrs"]] == ["#12"]


def test_show_returns_files_with_counts_and_threads():
    mod = load()
    diff = "@@ -1,2 +1,3 @@\n a\n-b\n+c\n+d\n"

    def call(cmd):
        if cmd[2].endswith("/diffs?per_page=100"):
            return [{"new_path": "a.py", "diff": diff}]
        return [{"notes": [{"resolvable": True, "resolved": False, "author": {"username": "sam"}, "body": "x",
                            "created_at": "2026-10-29T09:00:00Z", "position": {"new_path": "a.py", "new_line": 2}}]}]
    out = mod.gitlab_show("grp/demo-service", "5", call=call)
    assert out["files"] == [{"path": "a.py", "add": 2, "del": 1, "diff": diff}]
    assert out["threads"][0]["line"] == 2 and out["threads"][0]["resolved"] is False
