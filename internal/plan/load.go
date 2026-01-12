package plan

import (
    "encoding/json"
    "fmt"
    "os"
)

func Load(path string) (*Plan, error) {
    b, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }
    var p Plan
    if err := json.Unmarshal(b, &p); err != nil {
        return nil, fmt.Errorf("json parse: %w", err)
    }
    if err := p.Validate(); err != nil {
        return nil, err
    }
    return &p, nil
}

func Save(path string, p *Plan) error {
    if err := p.Validate(); err != nil {
        return err
    }
    out, err := json.MarshalIndent(p, "", "  ")
    if err != nil {
        return err
    }
    out = append(out, '\n')
    return os.WriteFile(path, out, 0644)
}

func (p *Plan) Validate() error {
    if p.Version != 1 {
        return fmt.Errorf("version must be 1")
    }
    if p.CurrentStep == "" {
        return fmt.Errorf("current_step must be set")
    }
    if len(p.Steps) == 0 {
        return fmt.Errorf("steps must not be empty")
    }
    ids := map[string]bool{}
    for i, s := range p.Steps {
        if s.ID == "" {
            return fmt.Errorf("steps[%d].id must be set", i)
        }
        if s.Objective == "" {
            return fmt.Errorf("steps[%d].objective must be set", i)
        }
        switch s.Status {
        case StatusLocked, StatusPending, StatusFailed, StatusVerified:
        default:
            return fmt.Errorf("steps[%d].status invalid: %s", i, s.Status)
        }
        if ids[s.ID] {
            return fmt.Errorf("duplicate step id: %s", s.ID)
        }
        ids[s.ID] = true
        if len(s.Verify) == 0 {
            return fmt.Errorf("steps[%d].verify must not be empty", i)
        }
        for j, r := range s.Verify {
            switch r.Type {
            case "command":
                if r.Cmd == "" {
                    return fmt.Errorf("steps[%d].verify[%d].cmd required for command", i, j)
                }
            case "file_exists":
                if r.File == "" {
                    return fmt.Errorf("steps[%d].verify[%d].file required for file_exists", i, j)
                }
            case "file_contains":
                if r.File == "" || r.Pattern == "" {
                    return fmt.Errorf("steps[%d].verify[%d].file+pattern required for file_contains", i, j)
                }
            case "http":
                if r.URL == "" {
                    return fmt.Errorf("steps[%d].verify[%d].url required for http", i, j)
                }
            default:
                return fmt.Errorf("steps[%d].verify[%d].type invalid: %s", i, j, r.Type)
            }
        }
    }
    if !ids[p.CurrentStep] {
        return fmt.Errorf("current_step %q not found in steps", p.CurrentStep)
    }
    return nil
}
