package main

import (
	"strings"
	"time"
)

// hiddenPages holds the pages of bricks switched off in crow.yaml (adr_002), as declared by
// each manifest's deck.pages. Tabs keep their numbers; a hidden page shows no tab and its
// keys do nothing.
var hiddenPages = map[int]bool{}

// hiddenBoxes holds the deck sources of off bricks (deck.boxes, item_237). They are never
// fetched, so their boxes stay empty and hide as they do with no data.
var hiddenBoxes = map[string]bool{}

func pageHidden(i int) bool { return hiddenPages[i] }

func boxHidden(src string) bool { return hiddenBoxes[src] }

type brickPages struct {
	Name  string   `json:"name"`
	On    bool     `json:"on"`
	Pages []string `json:"deck_pages"`
	Boxes []string `json:"deck_boxes"`
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

// setHiddenPages hides the pages and boxes of every brick that is off.
func setHiddenPages(mods []brickPages) {
	hiddenPages, hiddenBoxes = map[int]bool{}, map[string]bool{}
	for _, b := range mods {
		if b.On {
			continue
		}
		for _, src := range b.Boxes {
			hiddenBoxes[src] = true
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
