from conftest import run


def test_migrate_moves_notes_and_sets_crow_dir(hub, tmp_path):
    home, _ = hub
    vault = tmp_path / "vault"
    (vault / "40 Journal" / "Daily").mkdir(parents=True)
    (vault / "CR-0W" / "Rollups").mkdir(parents=True)
    (vault / "CR-0W" / "Rollups" / "2026-W44.md").write_text("w44\n")
    daily = vault / "40 Journal" / "Daily" / "2026-11-05.md"
    daily.write_text("# Day\n\n## 📥 Capture\n- a\n\n## 🌙 End of day\n- x\n")
    (home / "local" / "projects.yaml").write_text(f"vault:\n  path: {vault}\n  crow_dir: null\nprojects: []\n")
    st = home / "local" / "state"
    (st / "open-threads.md").write_text("# Open threads\n\n## Me\n- [ ] 2026-11-01 [hub] x\n")
    (st / "status.md").write_text("# Status\n")
    (st / "decisions.md").write_text("# Decisions\n\n## 2026-10-15 [hub] A/B split\nwhy a\n\n## 2026-10-16 Other\nwhy b\n")
    (st / "retros").mkdir()
    (st / "retros" / "2026-11-05-batch.md").write_text("r\n")
    (st / "retros" / "proposals.md").write_text("- [ ] p\n")
    (st / "activity").mkdir()
    (st / "activity" / "2026-11-05.md").write_text("- 09:00 [hub] brief\n")
    (st / "activity" / "2026-10-01.md").write_text("- 09:00 [hub] old\n")
    (st / "runs.jsonl").write_text("{}\n")

    dry = run(hub, "crow-vault-migrate")
    assert dry.returncode == 0, dry.stderr and (st / "open-threads.md").exists()
    r = run(hub, "crow-vault-migrate", "--apply")
    assert r.returncode == 0, r.stderr
    notes = vault / "CR-0W"
    assert (notes / "Open Threads.md").read_text().startswith("---\ntype: crow-threads\n---\n\n# Open threads")
    assert (notes / "Status.md").exists() and not (st / "status.md").exists()
    dec = sorted(p.name for p in (notes / "Decisions").iterdir())
    assert dec == ["2026-10-15 A-B split.md", "2026-10-16 Other.md"]
    assert "project: hub" in (notes / "Decisions" / "2026-10-15 A-B split.md").read_text()
    assert (notes / "Retros" / "2026-11-05-batch.md").exists() and (st / "retros" / "proposals.md").exists()
    assert (notes / "Rollups" / "2026-W44 rollup.md").exists()
    text = daily.read_text()
    assert "## CR-0W\n- 09:00 [hub] brief\n\n## 🌙 End of day" in text
    assert (st / "activity" / "2026-10-01.md").exists() and not (st / "activity" / "2026-11-05.md").exists()
    assert (st / "runs.jsonl").exists()
    assert "crow_dir: CR-0W" in (home / "local" / "projects.yaml").read_text()
    s = run(hub, "crow-state", "path", "threads")
    assert s.stdout.strip() == str(notes / "Open Threads.md")
