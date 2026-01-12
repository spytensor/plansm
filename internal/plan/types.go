package plan

type Plan struct {
    Version     int      `json:"version"`
    CurrentStep string   `json:"current_step"`
    Invariants  []string `json:"invariants,omitempty"`
    Steps       []Step   `json:"steps"`
}

type StepStatus string

const (
    StatusLocked   StepStatus = "LOCKED"
    StatusPending  StepStatus = "PENDING"
    StatusFailed   StepStatus = "FAILED"
    StatusVerified StepStatus = "VERIFIED"
)

type Step struct {
    ID         string       `json:"id"`
    Objective  string       `json:"objective"`
    Status     StepStatus   `json:"status"`
    DependsOn  []string     `json:"depends_on,omitempty"`
    AllowPaths []string     `json:"allow_paths,omitempty"`
    Verify     []VerifyRule `json:"verify"`
}

type VerifyRule struct {
    Type    string  `json:"type"`
    Cmd     string  `json:"cmd,omitempty"`
    File    string  `json:"file,omitempty"`
    Pattern string  `json:"pattern,omitempty"`
    URL     string  `json:"url,omitempty"`
    Method  string  `json:"method,omitempty"`
    Expect  *Expect `json:"expect,omitempty"`
}

type Expect struct {
    ExitCode    *int    `json:"exit_code,omitempty"`
    StdoutRegex *string `json:"stdout_regex,omitempty"`
    HTTPStatus  *int    `json:"http_status,omitempty"`
}
