package plan

import "testing"

func TestUnlockReady(t *testing.T) {
    p := &Plan{
        Version:     1,
        CurrentStep: "A",
        Steps: []Step{
            {ID: "A", Objective: "a", Status: StatusVerified, Verify: []VerifyRule{{Type: "command", Cmd: "true"}}},
            {ID: "B", Objective: "b", Status: StatusLocked, DependsOn: []string{"A"}, Verify: []VerifyRule{{Type: "command", Cmd: "true"}}},
        },
    }
    p.UnlockReady()
    if p.Steps[1].Status != StatusPending {
        t.Fatalf("expected B pending, got %s", p.Steps[1].Status)
    }
}
