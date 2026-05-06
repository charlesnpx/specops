package runstate

import "fmt"

var transitions = map[Status][]Status{
	StatusCreated:           {StatusIngested},
	StatusIngested:          {StatusIngested, StatusIntakeComplete},
	StatusIntakeComplete:    {StatusRefined},
	StatusRefined:           {StatusHardened, StatusAwaitingDecisions},
	StatusHardened:          {StatusAwaitingDecisions},
	StatusSynthesized:       {StatusAwaitingDecisions},
	StatusAwaitingDecisions: {StatusDecisionsAccepted},
	StatusDecisionsAccepted: {StatusCompiled},
	StatusCompiled:          {StatusPlanned},
	StatusPlanned:           {StatusApplied},
	StatusApplied:           {StatusAudited},
	StatusAudited:           {StatusEvaluated},
}

func Transition(state *RunState, next Status) error {
	allowed := transitions[state.Status]
	for _, candidate := range allowed {
		if candidate == next {
			state.Status = next
			return nil
		}
	}
	return fmt.Errorf("illegal transition from %s to %s", state.Status, next)
}

func SetStatus(state *RunState, status Status) {
	state.Status = status
}

func NextForStatus(status Status, runID string) *NextAction {
	switch status {
	case StatusCreated:
		return &NextAction{Command: "specops ingest-file <path> --run " + runID, Reason: "run has no normalized input yet"}
	case StatusIngested:
		return &NextAction{Command: "specops intake " + runID, Reason: "input is ready for intake"}
	case StatusIntakeComplete:
		return &NextAction{Command: "specops refine " + runID, Reason: "intake artifact is ready to refine"}
	case StatusRefined:
		return &NextAction{Command: "specops harden " + runID, Reason: "refined artifact can be challenged or synthesized"}
	case StatusHardened:
		return &NextAction{Command: "specops synthesize " + runID, Reason: "hardened artifact can produce a spec delta"}
	case StatusAwaitingDecisions:
		return &NextAction{Command: "specops decisions " + runID, Reason: "human decision gate is waiting"}
	case StatusDecisionsAccepted:
		return &NextAction{Command: "specops compile " + runID + " --accepted-only", Reason: "accepted decisions can be compiled"}
	case StatusCompiled:
		return &NextAction{Command: "specops plan " + runID, Reason: "patch plan is ready for review"}
	case StatusPlanned:
		return &NextAction{Command: "specops apply " + runID, Reason: "reviewed plan can be applied"}
	case StatusApplied:
		return &NextAction{Command: "specops audit", Reason: "applied changes should be audited"}
	case StatusAudited:
		return &NextAction{Command: "specops eval --gold <repo> --candidate <repo>", Reason: "audited repo can be evaluated"}
	default:
		return nil
	}
}
