package main

import "testing"

func TestAzertyTabKeys(t *testing.T) {
	keys := [][2]string{{"&", "1"}, {"é", "2"}, {"\"", "3"}, {"'", "4"}, {"(", "5"}}
	for i, k := range keys {
		for _, key := range k {
			m := send(newModel(defaultConfig()), resize(120, 40), keyMsg(key))
			if m.page != i {
				t.Errorf("key %q: page %d, want %d", key, m.page, i)
			}
		}
	}
}

func TestHourRangeKeysWithoutAltGr(t *testing.T) {
	m := send(newModel(defaultConfig()), resize(120, 40), keyMsg("2"))
	if m.now.Hour() < 5 {
		t.Skip("needs a clock past 05:00 to move back")
	}
	m = send(m, keyMsg("<"))
	if m.agentShift != 1 {
		t.Fatalf("< agentShift %d, want 1", m.agentShift)
	}
	m = send(m, keyMsg(">"))
	if m.agentShift != 0 {
		t.Fatalf("> agentShift %d, want 0", m.agentShift)
	}
}
