package plan

import "fmt"

func (p *Plan) StepByID(id string) (*Step, int) {
    for i := range p.Steps {
        if p.Steps[i].ID == id {
            return &p.Steps[i], i
        }
    }
    return nil, -1
}

func (p *Plan) Current() (*Step, int, error) {
    st, idx := p.StepByID(p.CurrentStep)
    if st == nil {
        return nil, -1, fmt.Errorf("current_step %q not found in steps", p.CurrentStep)
    }
    return st, idx, nil
}

func (p *Plan) UnlockReady() {
    for i := range p.Steps {
        if p.Steps[i].Status != StatusLocked {
            continue
        }
        ok := true
        for _, dep := range p.Steps[i].DependsOn {
            s, _ := p.StepByID(dep)
            if s == nil || s.Status != StatusVerified {
                ok = false
                break
            }
        }
        if ok {
            p.Steps[i].Status = StatusPending
        }
    }
}

func (p *Plan) NextPendingAfter(idx int) (string, bool) {
    for j := idx + 1; j < len(p.Steps); j++ {
        if p.Steps[j].Status == StatusPending || p.Steps[j].Status == StatusFailed {
            return p.Steps[j].ID, true
        }
    }
    for j := 0; j < len(p.Steps); j++ {
        if p.Steps[j].Status == StatusPending || p.Steps[j].Status == StatusFailed {
            return p.Steps[j].ID, true
        }
    }
    return "", false
}
