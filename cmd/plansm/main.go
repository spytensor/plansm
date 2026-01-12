package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "os"
    "path/filepath"
    "strings"

    "github.com/spytensor/plansm/internal/plan"
    "github.com/spytensor/plansm/internal/verify"
)

const defaultPlanPath = "plan.json"

func main() {
    if len(os.Args) < 2 {
        usage()
        os.Exit(2)
    }
    switch os.Args[1] {
    case "init":
        cmdInit(os.Args[2:])
    case "status":
        cmdStatus(os.Args[2:])
    case "current":
        cmdCurrent(os.Args[2:])
    case "verify":
        cmdVerify(os.Args[2:])
    case "advance":
        cmdAdvance(os.Args[2:])
    case "doctor":
        cmdDoctor(os.Args[2:])
    default:
        fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
        usage()
        os.Exit(2)
    }
}

func usage() {
    fmt.Print(`plansm — Plan as a Verifiable State Machine

Usage:
  plansm init [--plan plan.json] [--claude]
  plansm status [--plan plan.json]
  plansm current [--plan plan.json] [--json]
  plansm verify [--plan plan.json] (--current|--all) [--json]
  plansm advance [--plan plan.json]
  plansm doctor [--plan plan.json]
`)
}

func cmdInit(args []string) {
    fs := flag.NewFlagSet("init", flag.ExitOnError)
    planPath := fs.String("plan", defaultPlanPath, "path to plan.json")
    withClaude := fs.Bool("claude", false, "create Claude Code integration (.claude/ + .claude-plugin/)")
    _ = fs.Parse(args)

    if _, err := os.Stat(*planPath); err == nil {
        fmt.Fprintf(os.Stderr, "refusing to overwrite existing %s\n", *planPath)
        os.Exit(1)
    }

    p := plan.Plan{
        Version:     1,
        CurrentStep: "STEP_001",
        Invariants: []string{
            "Do not mark VERIFIED without running plansm verify.",
            "Only work on current_step unless explicitly unlocked.",
        },
        Steps: []plan.Step{
            {
                ID:        "STEP_001",
                Objective: "Initialize project skeleton",
                Status:    plan.StatusPending,
                Verify:    []plan.VerifyRule{{Type: "command", Cmd: "test -d ."}},
            },
            {
                ID:        "STEP_002",
                Objective: "Add first real proof rules (edit this step in plan.json)",
                Status:    plan.StatusLocked,
                DependsOn: []string{"STEP_001"},
                Verify:    []plan.VerifyRule{{Type: "command", Cmd: "echo \"edit plan.json proofs\""}},
            },
        },
    }

    if err := plan.Save(*planPath, &p); err != nil {
        fatal(err)
    }
    fmt.Printf("created %s\n", *planPath)

    if *withClaude {
        if err := installClaudeProjectFiles(filepath.Dir(*planPath)); err != nil {
            fatal(err)
        }
        fmt.Println("created .claude/commands and .claude-plugin/ (project-local)")
        fmt.Println("Tip: restart Claude Code session for hooks to take effect.")
    }
}

func cmdStatus(args []string) {
    fs := flag.NewFlagSet("status", flag.ExitOnError)
    planPath := fs.String("plan", defaultPlanPath, "path to plan.json")
    _ = fs.Parse(args)

    p := mustLoad(*planPath)
    p.UnlockReady()

    fmt.Printf("current_step: %s\n\n", p.CurrentStep)
    fmt.Printf("%-10s  %-9s  %s\n", "STEP", "STATUS", "OBJECTIVE")
    fmt.Printf("%s\n", strings.Repeat("-", 80))
    for _, s := range p.Steps {
        obj := s.Objective
        if len(obj) > 55 {
            obj = obj[:55] + "…"
        }
        fmt.Printf("%-10s  %-9s  %s\n", s.ID, s.Status, obj)
    }
}

func cmdCurrent(args []string) {
    fs := flag.NewFlagSet("current", flag.ExitOnError)
    planPath := fs.String("plan", defaultPlanPath, "path to plan.json")
    asJSON := fs.Bool("json", false, "json output")
    _ = fs.Parse(args)

    p := mustLoad(*planPath)
    p.UnlockReady()
    st, _, err := p.Current()
    if err != nil {
        fatal(err)
    }

    out := map[string]any{
        "current_step": st.ID,
        "objective":    st.Objective,
        "status":       st.Status,
        "allow_paths":  st.AllowPaths,
        "verify":       st.Verify,
    }

    if *asJSON {
        enc := json.NewEncoder(os.Stdout)
        enc.SetIndent("", "  ")
        _ = enc.Encode(out)
        return
    }

    fmt.Printf("CURRENT_STEP: %s\n", st.ID)
    fmt.Printf("STATUS: %s\n", st.Status)
    fmt.Printf("OBJECTIVE: %s\n", st.Objective)
    if len(st.AllowPaths) > 0 {
        fmt.Println("ALLOW_PATHS:")
        for _, p := range st.AllowPaths {
            fmt.Printf("  - %s\n", p)
        }
    }
    fmt.Println("VERIFY:")
    for _, v := range st.Verify {
        switch v.Type {
        case "command":
            fmt.Printf("  - command: %s\n", v.Cmd)
        case "file_exists":
            fmt.Printf("  - file_exists: %s\n", v.File)
        case "file_contains":
            fmt.Printf("  - file_contains: %s pattern=%s\n", v.File, v.Pattern)
        case "http":
            fmt.Printf("  - http: %s\n", v.URL)
        default:
            fmt.Printf("  - %s\n", v.Type)
        }
    }
}

func cmdVerify(args []string) {
    fs := flag.NewFlagSet("verify", flag.ExitOnError)
    planPath := fs.String("plan", defaultPlanPath, "path to plan.json")
    current := fs.Bool("current", false, "verify current step")
    all := fs.Bool("all", false, "verify all PENDING/FAILED steps (in order)")
    asJSON := fs.Bool("json", false, "json output")
    _ = fs.Parse(args)

    if (*current && *all) || (!*current && !*all) {
        fmt.Fprintln(os.Stderr, "must pass exactly one of --current or --all")
        os.Exit(2)
    }

    p := mustLoad(*planPath)
    p.UnlockReady()
    wd := filepath.Dir(*planPath)

    var stepResults []verify.StepResult
    if *current {
        st, _, err := p.Current()
        if err != nil {
            fatal(err)
        }
        if st.Status == plan.StatusLocked {
            fatal(fmt.Errorf("current step is LOCKED; dependencies not met"))
        }
        r := verify.RunStep(wd, *st)
        stepResults = append(stepResults, r)
        applyVerifyResult(&p, st.ID, r.Ok)
    } else {
        for i := range p.Steps {
            if p.Steps[i].Status == plan.StatusPending || p.Steps[i].Status == plan.StatusFailed {
                r := verify.RunStep(wd, p.Steps[i])
                stepResults = append(stepResults, r)
                applyVerifyResult(&p, p.Steps[i].ID, r.Ok)
                if !r.Ok {
                    break
                }
            }
        }
    }

    if err := plan.Save(*planPath, &p); err != nil {
        fatal(err)
    }

    overall := true
    for _, sr := range stepResults {
        if !sr.Ok {
            overall = false
        }
    }

    if *asJSON {
        enc := json.NewEncoder(os.Stdout)
        enc.SetIndent("", "  ")
        _ = enc.Encode(map[string]any{"ok": overall, "results": stepResults})
    } else {
        for _, sr := range stepResults {
            fmt.Printf("STEP %s: %s\n", sr.StepID, ternary(sr.Ok, "OK", "FAILED"))
            for _, rr := range sr.Results {
                if rr.Ok {
                    fmt.Printf("  ✓ %s\n", ruleShort(rr.Rule))
                } else {
                    fmt.Printf("  ✗ %s — %s\n", ruleShort(rr.Rule), rr.Detail)
                }
            }
        }
        fmt.Printf("\nOVERALL: %s\n", ternary(overall, "OK", "FAILED"))
    }

    if !overall {
        os.Exit(1)
    }
}

func cmdAdvance(args []string) {
    fs := flag.NewFlagSet("advance", flag.ExitOnError)
    planPath := fs.String("plan", defaultPlanPath, "path to plan.json")
    _ = fs.Parse(args)

    p := mustLoad(*planPath)
    p.UnlockReady()
    st, idx, err := p.Current()
    if err != nil {
        fatal(err)
    }
    if st.Status != plan.StatusVerified {
        fatal(fmt.Errorf("cannot advance: current step %s status=%s (must be VERIFIED). Run `plansm verify --current` first", st.ID, st.Status))
    }

    p.UnlockReady()
    next, ok := p.NextPendingAfter(idx)
    if !ok {
        fmt.Println("No more pending steps. Plan complete.")
        return
    }
    p.CurrentStep = next
    if err := plan.Save(*planPath, &p); err != nil {
        fatal(err)
    }
    fmt.Printf("advanced to %s\n", next)
}

func cmdDoctor(args []string) {
    fs := flag.NewFlagSet("doctor", flag.ExitOnError)
    planPath := fs.String("plan", defaultPlanPath, "path to plan.json")
    _ = fs.Parse(args)

    fmt.Println("plansm doctor")
    fmt.Println("------------")
    p, err := plan.Load(*planPath)
    if err != nil {
        fmt.Printf("❌ plan load/validate: %v\n", err)
        os.Exit(1)
    }
    fmt.Printf("✅ plan load/validate ok (version=%d)\n", p.Version)

    st, _, err := p.Current()
    if err != nil {
        fmt.Printf("❌ current step: %v\n", err)
        os.Exit(1)
    }
    fmt.Printf("✅ current step ok: %s (%s)\n", st.ID, st.Status)

    wd := filepath.Dir(*planPath)
    if _, err := os.Stat(filepath.Join(wd, ".claude", "commands")); err == nil {
        fmt.Println("✅ .claude/commands present")
    } else {
        fmt.Println("ℹ️  .claude/commands not found (run `plansm init --claude` if desired)")
    }
    if _, err := os.Stat(filepath.Join(wd, ".claude-plugin", "plugin.json")); err == nil {
        fmt.Println("✅ .claude-plugin present")
    } else {
        fmt.Println("ℹ️  .claude-plugin not found (run `plansm init --claude` if desired)")
    }
}

func mustLoad(path string) plan.Plan {
    p, err := plan.Load(path)
    if err != nil {
        fatal(err)
    }
    return *p
}

func fatal(err error) {
    fmt.Fprintln(os.Stderr, "error:", err)
    os.Exit(1)
}

func applyVerifyResult(p *plan.Plan, stepID string, ok bool) {
    st, _ := p.StepByID(stepID)
    if st == nil {
        return
    }
    if ok {
        st.Status = plan.StatusVerified
        p.UnlockReady()
    } else {
        st.Status = plan.StatusFailed
    }
}

func ruleShort(r plan.VerifyRule) string {
    switch r.Type {
    case "command":
        return "command: " + r.Cmd
    case "file_exists":
        return "file_exists: " + r.File
    case "file_contains":
        return "file_contains: " + r.File
    case "http":
        return "http: " + r.URL
    default:
        return r.Type
    }
}

func ternary[T any](cond bool, a, b T) T {
    if cond {
        return a
    }
    return b
}

func installClaudeProjectFiles(projectDir string) error {
    cmdsDir := filepath.Join(projectDir, ".claude", "commands")
    if err := os.MkdirAll(cmdsDir, 0755); err != nil {
        return err
    }
    write := func(name, body string) error {
        return os.WriteFile(filepath.Join(cmdsDir, name), []byte(body), 0644)
    }

    if err := write("pwork.md", "---\ndescription: Show the current step (low token).\n---\nRun: `plansm current`\n\nRules:\n- Do NOT edit plan.json status fields manually.\n- Only work on CURRENT_STEP.\n"); err != nil { return err }

    if err := write("pverify.md", "---\ndescription: Verify current step proofs (machine gate).\n---\nRun: `plansm verify --current`\n\nIf this fails, do NOT claim completion. Fix the root cause and run again.\n"); err != nil { return err }

    if err := write("pstatus.md", "---\ndescription: Show plan status table.\n---\nRun: `plansm status`\n"); err != nil { return err }

    if err := write("pnext.md", "---\ndescription: Advance to next step (only if current step is VERIFIED).\n---\nRun: `plansm advance`\n"); err != nil { return err }

    // Plugin skeleton for open-source distribution
    pluginDir := filepath.Join(projectDir, ".claude-plugin")
    if err := os.MkdirAll(filepath.Join(pluginDir, "commands"), 0755); err != nil { return err }
    if err := os.MkdirAll(filepath.Join(pluginDir, "hooks"), 0755); err != nil { return err }

    pluginJSON := `{
  "name": "plansm",
  "version": "0.1.0",
  "description": "Plan-as-state-machine with proof-based verification (no markdown).",
  "commands": {
    "pwork": "commands/pwork.md",
    "pverify": "commands/pverify.md",
    "pstatus": "commands/pstatus.md",
    "pnext": "commands/pnext.md"
  },
  "hooks": "hooks/hooks.json"
}
`
    if err := os.WriteFile(filepath.Join(pluginDir, "plugin.json"), []byte(pluginJSON), 0644); err != nil { return err }

    for _, n := range []string{"pwork.md", "pverify.md", "pstatus.md", "pnext.md"} {
        b, err := os.ReadFile(filepath.Join(cmdsDir, n))
        if err != nil { return err }
        if err := os.WriteFile(filepath.Join(pluginDir, "commands", n), b, 0644); err != nil { return err }
    }

    // NOTE: Hook schema can vary by Claude Code version.
    // Users can configure via `/hooks` UI and set Stop hook to: `plansm verify --current`.
    hooksJSON := `{
  "hooks": {
    "Stop": [
      {
        "matcher": "",
        "command": "plansm verify --current",
        "blocking": true
      }
    ]
  }
}
`
    _ = os.WriteFile(filepath.Join(pluginDir, "hooks", "hooks.json"), []byte(hooksJSON), 0644)

    // add ignore line for optional local settings
    gi := filepath.Join(projectDir, ".gitignore")
    addGitignoreLine(gi, ".claude/settings.local.json")

    return nil
}

func addGitignoreLine(path, line string) {
    b, _ := os.ReadFile(path)
    s := string(b)
    if strings.Contains(s, line) {
        return
    }
    if len(s) > 0 && !strings.HasSuffix(s, "\n") {
        s += "\n"
    }
    s += line + "\n"
    _ = os.WriteFile(path, []byte(s), 0644)
}
