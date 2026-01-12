package verify

import (
    "bytes"
    "fmt"
    "io"
    "net/http"
    "os"
    "os/exec"
    "regexp"
    "strings"
    "time"

    "github.com/spytensor/plansm/internal/plan"
)

type RuleResult struct {
    Rule   plan.VerifyRule `json:"rule"`
    Ok     bool            `json:"ok"`
    Detail string          `json:"detail,omitempty"`
}

type StepResult struct {
    StepID  string       `json:"step_id"`
    Ok      bool         `json:"ok"`
    Results []RuleResult `json:"results"`
}

func RunStep(workdir string, step plan.Step) StepResult {
    res := StepResult{StepID: step.ID, Ok: true}
    for _, r := range step.Verify {
        rr := runRule(workdir, r)
        res.Results = append(res.Results, rr)
        if !rr.Ok {
            res.Ok = false
        }
    }
    return res
}

func runRule(workdir string, r plan.VerifyRule) RuleResult {
    switch r.Type {
    case "command":
        return runCommand(workdir, r)
    case "file_exists":
        return runFileExists(workdir, r)
    case "file_contains":
        return runFileContains(workdir, r)
    case "http":
        return runHTTP(r)
    default:
        return RuleResult{Rule: r, Ok: false, Detail: "unknown rule type"}
    }
}

func runCommand(workdir string, r plan.VerifyRule) RuleResult {
    cmd := exec.Command("bash", "-lc", r.Cmd)
    cmd.Dir = workdir
    var stdout bytes.Buffer
    var stderr bytes.Buffer
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr
    err := cmd.Run()

    exitCode := 0
    if err != nil {
        if ee, ok := err.(*exec.ExitError); ok {
            exitCode = ee.ExitCode()
        } else {
            return RuleResult{Rule: r, Ok: false, Detail: "exec error: " + err.Error()}
        }
    }

    if r.Expect != nil && r.Expect.ExitCode != nil {
        if exitCode != *r.Expect.ExitCode {
            return RuleResult{Rule: r, Ok: false, Detail: fmt.Sprintf("exit_code=%d expected=%d; stderr=%s", exitCode, *r.Expect.ExitCode, trim(stderr.String()))}
        }
    } else {
        if exitCode != 0 {
            return RuleResult{Rule: r, Ok: false, Detail: fmt.Sprintf("exit_code=%d; stderr=%s", exitCode, trim(stderr.String()))}
        }
    }

    if r.Expect != nil && r.Expect.StdoutRegex != nil {
        re, e := regexp.Compile(*r.Expect.StdoutRegex)
        if e != nil {
            return RuleResult{Rule: r, Ok: false, Detail: "bad stdout_regex: " + e.Error()}
        }
        if !re.MatchString(stdout.String()) {
            return RuleResult{Rule: r, Ok: false, Detail: "stdout_regex did not match"}
        }
    }

    return RuleResult{Rule: r, Ok: true}
}

func runFileExists(workdir string, r plan.VerifyRule) RuleResult {
    p := join(workdir, r.File)
    if _, err := os.Stat(p); err != nil {
        return RuleResult{Rule: r, Ok: false, Detail: err.Error()}
    }
    return RuleResult{Rule: r, Ok: true}
}

func runFileContains(workdir string, r plan.VerifyRule) RuleResult {
    p := join(workdir, r.File)
    b, err := os.ReadFile(p)
    if err != nil {
        return RuleResult{Rule: r, Ok: false, Detail: err.Error()}
    }
    matched, e := regexp.Match(r.Pattern, b)
    if e != nil {
        return RuleResult{Rule: r, Ok: false, Detail: "bad pattern: " + e.Error()}
    }
    if !matched {
        return RuleResult{Rule: r, Ok: false, Detail: "pattern not found"}
    }
    return RuleResult{Rule: r, Ok: true}
}

func runHTTP(r plan.VerifyRule) RuleResult {
    method := r.Method
    if method == "" {
        method = "GET"
    }
    client := &http.Client{Timeout: 8 * time.Second}
    req, err := http.NewRequest(method, r.URL, nil)
    if err != nil {
        return RuleResult{Rule: r, Ok: false, Detail: err.Error()}
    }
    resp, err := client.Do(req)
    if err != nil {
        return RuleResult{Rule: r, Ok: false, Detail: err.Error()}
    }
    defer resp.Body.Close()
    io.Copy(io.Discard, resp.Body)

    expected := 200
    if r.Expect != nil && r.Expect.HTTPStatus != nil {
        expected = *r.Expect.HTTPStatus
    }
    if resp.StatusCode != expected {
        return RuleResult{Rule: r, Ok: false, Detail: fmt.Sprintf("http_status=%d expected=%d", resp.StatusCode, expected)}
    }
    return RuleResult{Rule: r, Ok: true}
}

func trim(s string) string {
    s = strings.TrimSpace(s)
    if len(s) > 400 {
        return s[:400] + "…"
    }
    return s
}

func join(dir, file string) string {
    file = strings.TrimPrefix(file, "./")
    if strings.HasSuffix(dir, "/") {
        return dir + file
    }
    return dir + "/" + file
}
