package main

import (
	"strings"
	"time"
)

// hiddenPages holds the pages of bricks switched off in crow.yaml (adr_002), as declared by
// each manifest's deck.pages. Tabs keep their numbers; a hidden page shows no tab and its
// keys do nothing.
var hiddenPages = map[int]bool{}

func pageHidden(i int) bool { return hiddenPages[i] }

type brickPages struct {
	Name  string   `json:"name"`
	On    bool     `json:"on"`
	Pages []string `json:"deck_pages"`
}

// loadHiddenPages reads the brick switches once at start. A failure keeps every page.
func loadHiddenPages() {
	var r struct {
		Modules []brickPages `json:"modules"`
	}
	if runJSON(5*time.Second, nil, &r, "crow-modules", "json") != nil {
		return
	}
	setHiddenPages(r.Modules)
}

// setHiddenPages hides the pages of every brick that is off.
func setHiddenPages(mods []brickPages) {
	hiddenPages = map[int]bool{}
	for _, b := range mods {
		if b.On {
			continue
		}
		for _, name := range b.Pages {
			for i, n := range pageNames {
				if strings.EqualFold(n, name) {
					hiddenPages[i] = true
				}
			}
		}
	}
}

// goPage switches to page i unless its brick is off.
func (m *model) goPage(i int) {
	if i >= 0 && i < len(pageNames) && !pageHidden(i) {
		m.page = i
	}
}
