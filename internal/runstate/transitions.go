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
	contextCommand := "specops context " + runID
	noteCommand := func(stage string) string {
		return "specops note " + runID + " --stage " + stage + " --text <file-or-inline>"
	}
	switch status {
	case StatusCreated:
		return mechanical("ingest", "specops ingest-file <path> --run "+runID, "run has no normalized input yet", contextCommand)
	case StatusIngested:
		return mechanical("intake", "specops intake "+runID, "input is ready for intake", contextCommand)
	case StatusIntakeComplete:
		return semantic("refine", "specops refine "+runID+" --from <file>", "intake artifact is ready to refine", contextCommand, noteCommand("refine"), []string{
			"What should the refinement preserve from the source material?",
			"What ambiguities or missing constraints should be called out before synthesis?",
			"What would make the next artifact useful for review?",
		})
	case StatusRefined:
		return semantic("harden", "specops harden "+runID+" --from <file>", "refined artifact can be challenged or synthesized", contextCommand, noteCommand("harden"), []string{
			"What assumptions in the refined notes need pressure-testing?",
			"What failure modes or interface consequences should the hardening pass examine?",
			"Is the refined artifact ready to synthesize, or should it be challenged first?",
		})
	case StatusHardened:
		return semantic("synthesize", "specops synthesize "+runID+" --from <spec_delta.json>", "hardened artifact can produce a spec delta", contextCommand, noteCommand("synthesize"), []string{
			"What decisions should be explicit before canonical docs can change?",
			"Which docs are likely affected by the synthesized delta?",
			"What acceptance criteria should gate the patch plan?",
		})
	case StatusAwaitingDecisions:
		return semantic("decisions", "specops decisions "+runID, "human decision gate is waiting", contextCommand, noteCommand("decisions"), []string{
			"Which proposed decisions should be accepted, rejected, deferred, or amended?",
			"Does any decision require a new or superseding ADR?",
			"What rationale should be recorded with any deferrals or amendments?",
		})
	case StatusDecisionsAccepted:
		return mechanical("compile", "specops compile "+runID+" --accepted-only", "accepted decisions can be compiled", contextCommand)
	case StatusCompiled:
		return mechanical("plan", "specops plan "+runID, "patch plan is ready for review", contextCommand)
	case StatusPlanned:
		return semantic("apply", "specops apply "+runID, "reviewed plan can be applied", contextCommand, noteCommand("apply"), []string{
			"Has the compiled patch plan been reviewed against the accepted decisions?",
			"Should apply run as a dry run, interactive apply, or direct apply?",
			"Are there local files that need checking before mutation?",
		})
	case StatusApplied:
		return mechanical("audit", "specops audit", "applied changes should be audited", contextCommand)
	case StatusAudited:
		return mechanical("eval", "specops eval --gold <repo> --candidate <repo>", "audited repo can be evaluated", contextCommand)
	default:
		return nil
	}
}

func mechanical(stage, command, reason, contextCommand string) *NextAction {
	return &NextAction{
		Command:               command,
		Reason:                reason,
		Stage:                 stage,
		GateKind:              "mechanical",
		ContextCommand:        contextCommand,
		HumanInputRecommended: false,
	}
}

func semantic(stage, command, reason, contextCommand, noteCommand string, questions []string) *NextAction {
	return &NextAction{
		Command:               command,
		Reason:                reason,
		Stage:                 stage,
		GateKind:              "semantic",
		ContextCommand:        contextCommand,
		NoteCommand:           noteCommand,
		SuggestedQuestions:    questions,
		HumanInputRecommended: true,
	}
}
