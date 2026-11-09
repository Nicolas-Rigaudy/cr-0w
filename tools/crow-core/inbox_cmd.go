package main

import (
	"flag"
	"fmt"
	"io"
	"slices"
	"strings"
	"text/tabwriter"

	"crow-core/internal/mission"
)

// inboxCmd runs `inbox list|show|answer`.
func inboxCmd(env mission.Env, args []string, stdout io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("inbox: want list, show or answer")
	}
	sub, args := args[0], args[1:]
	fs := flag.NewFlagSet("inbox "+sub, flag.ContinueOnError)
	asJSON := fs.Bool("json", false, "print JSON")
	all := fs.Bool("all", false, "list answered and dropped items too")
	switch sub {
	case "list", "show", "answer":
	default:
		return fmt.Errorf("inbox: unknown subcommand %q", sub)
	}
	// Flags may follow the positionals. Words after -- are never flags.
	var tail []string
	if i := slices.Index(args, "--"); i >= 0 {
		args, tail = args[:i], args[i+1:]
	}
	var pos []string
	for {
		if err := fs.Parse(args); err != nil {
			return err
		}
		if fs.NArg() == 0 {
			break
		}
		pos = append(pos, fs.Arg(0))
		args = fs.Args()[1:]
	}
	pos = append(pos, tail...)

	switch sub {
	case "list":
		if len(pos) != 0 {
			return fmt.Errorf("inbox list: unexpected argument %q", pos[0])
		}
		items, err := env.Inbox(*all)
		if err != nil {
			return err
		}
		if *asJSON {
			return printJSON(stdout, items)
		}
		w := tabwriter.NewWriter(stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ITEM\tKIND\tTYPE\tVERBS\tASK")
		for _, it := range items {
			ask := it.Ask
			if it.Blocking {
				ask = "[blocking] " + ask
			}
			if it.State != "open" {
				ask = "(" + it.State + ") " + ask
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", it.ID, it.Kind, it.Type, strings.Join(it.Verbs, " "), ask)
		}
		return w.Flush()
	case "show":
		if len(pos) != 1 {
			return fmt.Errorf("inbox show: want one item id")
		}
		d, err := env.InboxShow(pos[0])
		if err != nil {
			return err
		}
		if *asJSON {
			return printJSON(stdout, d)
		}
		printItem(stdout, d)
		return nil
	default: // answer
		if len(pos) < 2 {
			return fmt.Errorf("inbox answer: want <item> <verb> [words], or <item> <words> for a question")
		}
		d, err := env.InboxShow(pos[0])
		if err != nil {
			return err
		}
		verb, words := pos[1], pos[2:]
		if d.Item.Source == "questions" && verb != "answer" {
			verb, words = "answer", pos[1:]
		}
		r := mission.Reply{Verb: verb, Surface: "cli"}
		if w := strings.Join(words, " "); w != "" {
			r.Words = []string{w}
		}
		res, err := env.InboxReply(d.Item.ID, r)
		if err != nil {
			return err
		}
		if *asJSON {
			return printJSON(stdout, res)
		}
		fmt.Fprintf(stdout, "%s  %s -> %s, mission %s\n", res.Item, res.Verb, res.File, res.State)
		if len(res.Validations) > 0 {
			fmt.Fprintf(stdout, "  validated %s\n", strings.Join(res.Validations, ", "))
		}
		if len(res.Dropped) > 0 {
			fmt.Fprintf(stdout, "  dropped   %s\n", strings.Join(res.Dropped, ", "))
		}
		return nil
	}
}

func printItem(w io.Writer, d mission.ItemDetail) {
	it := d.Item
	fmt.Fprintf(w, "%s  %s %s, %s\n", it.ID, it.Type, it.Source, it.State)
	fmt.Fprintf(w, "  mission   %s (%s, %s)\n", it.Mission, it.Kind, it.Project)
	fmt.Fprintf(w, "  file      %s\n", d.Path)
	fmt.Fprintf(w, "  verbs     %s\n", strings.Join(it.Verbs, " "))
	if d.Summary != "" {
		fmt.Fprintf(w, "  summary   %s\n", strings.Join(strings.Fields(d.Summary), " "))
	} else {
		fmt.Fprintf(w, "  ask       %s\n", it.Ask)
	}
	for _, q := range d.Questions {
		flag := ""
		if q.Blocking {
			flag = " [blocking]"
		}
		fmt.Fprintf(w, "  %s%s  %s\n", q.ID, flag, q.Text)
		if q.Choices == nil {
			continue
		}
		for _, c := range *q.Choices {
			mark := " "
			if q.Recommended != nil && *q.Recommended == c.ID {
				mark = "*"
			}
			line := fmt.Sprintf("    %s %s  %s", mark, c.ID, c.Label)
			if c.Note != nil {
				line += ": " + *c.Note
			}
			if c.Preview != nil {
				line += " (" + *c.Preview + ")"
			}
			fmt.Fprintln(w, line)
		}
	}
	for _, p := range d.Publish {
		fmt.Fprintf(w, "  publish   %s %s\n", p.Action, p.Target)
	}
	for _, a := range it.Artifacts {
		fmt.Fprintf(w, "  artifact  %s\n", a)
	}
	for i, p := range d.Answers {
		fmt.Fprintf(w, "  a%d        %s\n", i+1, p)
	}
}
