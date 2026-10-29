package main

// Pages other than Home are placeholders in the skeleton: the real boxes arrive with
// the later slices of req_007. They show the chrome, the planned boxes and the slice.

type pageInfo struct {
	slice string
	boxes []string
}

var pagePlans = map[int]pageInfo{
	1: {"item_017 (real data on Home and Agents)", []string{
		"Agents: needs you, running, background runs, finished today",
		"Live tail with the agent's question and the a x e c answer",
		"Who worked when: one line per agent, hour ruler, ▼ now",
	}},
	2: {"item_018 (MRs from the forges, read only)", []string{
		"Per project list of every open MR or PR assigned to you",
		"Diff, findings or threads, one card, a x e c actions",
	}},
	3: {"item_019 (Projects and Usage)", []string{
		"One card per project: NOW, NEXT, BLOCK, ROAD, MILE, OPEN",
		"Timeline with ▼ today and ◆ milestones",
	}},
	4: {"item_019 (Projects and Usage)", []string{
		"Accounts and policy: 5h and weekly gauges",
		"Forecast, daily spend, cost today",
	}},
}

func (lk *look) placeholder(page, w, h int, sync string) []string {
	plan := pagePlans[page]
	sp := boxSpec{title: pageNames[page], index: 1, focused: true, sel: -1}
	if sync == stLoading {
		sp.rows = []line{{{"syncing…", lk.sub}}}
	} else {
		sp.rows = append(sp.rows, line{{"Not built yet. Planned in " + plan.slice + ".", lk.sub}}, nil)
		for _, b := range plan.boxes {
			sp.rows = append(sp.rows, line{{"· ", lk.dim}, {b, lk.plain}})
		}
	}
	out := lk.box(sp, w, min(h, lk.boxH(len(sp.rows))))
	for len(out) < h {
		out = append(out, blank(w))
	}
	return out
}
