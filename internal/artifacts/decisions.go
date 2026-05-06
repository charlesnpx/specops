package artifacts

import (
	"os"
	"strings"

	"github.com/specops/specops/internal/runstate"
)

func SetDecision(repo, runID, decisionID, status, text string) (*runstate.RunState, error) {
	state, err := runstate.Load(repo, runID)
	if err != nil {
		return nil, err
	}
	decision, ok := state.Decisions[decisionID]
	if !ok {
		return nil, os.ErrNotExist
	}
	decision.Status = status
	if text != "" {
		decision.Text = text
	}
	state.Decisions[decisionID] = decision
	if allDecisionsSettled(state) {
		state.Status = runstate.StatusDecisionsAccepted
	}
	if err := runstate.Save(repo, state); err != nil {
		return nil, err
	}
	return state, nil
}

func AcceptRecommended(repo, runID string) (*runstate.RunState, error) {
	state, err := runstate.Load(repo, runID)
	if err != nil {
		return nil, err
	}
	for id, decision := range state.Decisions {
		if strings.EqualFold(decision.Recommendation, "accept") && decision.Status == "proposed" {
			decision.Status = "accepted"
			state.Decisions[id] = decision
		}
	}
	if allDecisionsSettled(state) {
		state.Status = runstate.StatusDecisionsAccepted
	}
	if err := runstate.Save(repo, state); err != nil {
		return nil, err
	}
	return state, nil
}

func allDecisionsSettled(state *runstate.RunState) bool {
	if len(state.Decisions) == 0 {
		return false
	}
	for _, decision := range state.Decisions {
		if decision.Status == "proposed" || decision.Status == "amended" {
			return false
		}
	}
	return true
}
