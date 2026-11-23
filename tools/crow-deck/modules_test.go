package main

import "testing"

func TestForgeOffHidesMRsPage(t *testing.T) {
	setHiddenPages([]brickPages{{Name: "forge", Pages: []string{"MRs"}}})
	t.Cleanup(func() { setHiddenPages(nil) })
	m := newModel(defaultConfig())
	m.w, m.h = 120, 40
	for _, k := range []string{"3", "\"", ":tab MRs"} {
		n := m
		if k[0] == ':' {
			n.exec(k[1:])
		} else {
			n = send(m, keyMsg(k))
		}
		if n.page == 2 {
			t.Errorf("%q opened the hidden MRs page", k)
		}
	}
	if n := send(m, keyMsg("4")); n.page != 3 {
		t.Errorf("4 went to page %d, want 3 (numbers stay)", n.page)
	}
	if hasText(m.View(), "MRs") {
		t.Error("tab bar still shows MRs")
	}
}

func TestForgeOnKeepsMRsPage(t *testing.T) {
	setHiddenPages([]brickPages{{Name: "forge", On: true, Pages: []string{"MRs"}}})
	t.Cleanup(func() { setHiddenPages(nil) })
	if n := send(newModel(defaultConfig()), keyMsg("3")); n.page != 2 {
		t.Errorf("3 went to page %d, want 2", n.page)
	}
}

func TestMailOffSkipsTodaySource(t *testing.T) {
	setHiddenPages([]brickPages{{Name: "mail", Boxes: []string{"today"}}, {Name: "team", On: true, Boxes: []string{"checks"}}})
	t.Cleanup(func() { setHiddenPages(nil) })
	if fetchCmd(srcToday) != nil {
		t.Error("mail off: today is still fetched")
	}
	if fetchCmd(srcChecks) == nil {
		t.Error("team on: checks not fetched")
	}
	m := newModel(defaultConfig())
	m.w, m.h = 120, 40
	if pageHidden(0) {
		t.Error("a box switch hid a page")
	}
}

func TestOffBricksNamesOnlyOffBricks(t *testing.T) {
	setHiddenPages([]brickPages{{Name: "retro"}, {Name: "afk", On: true}})
	t.Cleanup(func() { setHiddenPages(nil) })
	if !brickOff("retro") || brickOff("afk") || brickOff("unknown") {
		t.Fatalf("offBricks = %v", offBricks)
	}
}
