package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"crow-core/contract"
	"crow-core/internal/mission"
)

// stopRun stops a live run through `crow-runs stop` of the engine checkout.
func stopRun(engine string) func(string) error {
	if engine == "" {
		return nil
	}
	return func(run string) error {
		out, err := exec.Command(filepath.Join(engine, "bin", "crow-runs"), "stop", run).CombinedOutput()
		if err != nil {
			return fmt.Errorf("crow-runs stop: %v: %s", err, strings.TrimSpace(string(out)))
		}
		return nil
	}
}

// missionCmd runs `mission new|list|show|cancel`.
func missionCmd(env mission.Env, args []string, stdout io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("mission: want new, list, show or cancel")
	}
	sub, args := args[0], args[1:]
	fs := flag.NewFlagSet("mission "+sub, flag.ContinueOnError)
	asJSON := fs.Bool("json", false, "print JSON")
	var spec mission.Spec
	var gates, reason string
	switch sub {
	case "new":
		fs.StringVar(&spec.Goal, "goal", "", "what the mission is for")
		fs.StringVar(&spec.Kind, "kind", "", "kind file name")
		fs.StringVar(&spec.Project, "project", "", "registry name or alias")
		fs.StringVar(&spec.EndState, "end-state", "", "what done looks like")
		fs.StringVar(&spec.Key, "key", "", "short part of the id (default: from the goal)")
		fs.StringVar(&spec.Scope, "scope", "", "work, personal or crow (default: the project's)")
		fs.StringVar(&spec.DataClass, "data-class", "", "work, personal or private")
		fs.StringVar(&gates, "gates", "\x00", "comma list of plan, report, mockup (empty for none)")
		fs.StringVar(&spec.Autonomy, "autonomy", "", "ask, notify or auto")
	case "cancel":
		fs.StringVar(&reason, "reason", "", "why the mission stops")
	case "list", "show":
	default:
		return fmt.Errorf("mission: unknown subcommand %q", sub)
	}
	// Flags may follow the positional id.
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
	var ref string
	switch sub {
	case "show", "cancel":
		if len(pos) != 1 {
			return fmt.Errorf("mission %s: want one mission id", sub)
		}
		ref = pos[0]
	default:
		if len(pos) != 0 {
			return fmt.Errorf("mission %s: unexpected argument %q", sub, pos[0])
		}
	}

	switch sub {
	case "new":
		if gates != "\x00" {
			g := []string{}
			for _, x := range strings.Split(gates, ",") {
				if x = strings.TrimSpace(x); x != "" {
					g = append(g, x)
				}
			}
			spec.Gates = &g
		}
		m, err := env.New(spec)
		if err != nil {
			return err
		}
		return printMission(stdout, m, *asJSON)
	case "list":
		ms, err := env.List()
		if err != nil {
			return err
		}
		if *asJSON {
			if ms == nil {
				ms = []contract.Mission{}
			}
			return printJSON(stdout, ms)
		}
		w := tabwriter.NewWriter(stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tSTATE\tPROJECT\tGOAL")
		for _, m := range ms {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", m.ID, m.State, m.Project, m.Goal)
		}
		return w.Flush()
	case "show":
		m, err := env.Show(ref)
		if err != nil {
			return err
		}
		return printMission(stdout, m, *asJSON)
	default: // cancel
		m, err := env.Cancel(ref, reason, func(s string) { fmt.Fprintln(os.Stderr, "crow-core: warning:", s) })
		if err != nil {
			return err
		}
		return printMission(stdout, m, *asJSON)
	}
}

func printMission(w io.Writer, m contract.Mission, asJSON bool) error {
	if asJSON {
		return printJSON(w, m)
	}
	fmt.Fprintf(w, "%s  %s\n", m.ID, m.State)
	fmt.Fprintf(w, "  project   %s (%s, data %s)\n", m.Project, m.Scope, m.DataClass)
	fmt.Fprintf(w, "  goal      %s\n", m.Goal)
	fmt.Fprintf(w, "  end state %s\n", m.EndState)
	fmt.Fprintf(w, "  autonomy  %s, gates [%s]\n", m.Autonomy, strings.Join(m.Gates, ", "))
	fmt.Fprintf(w, "  budget    %d min, $%.2f, %d resumes\n", m.Budget.TimeMin, m.Budget.BudgetUSD, m.Budget.Resumes)
	fmt.Fprintf(w, "  runs      %d, version %d\n", len(m.Runs), m.Version)
	return nil
}
