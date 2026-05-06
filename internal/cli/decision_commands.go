package cli

import (
	"errors"
	"os"

	"github.com/specops/specops/internal/artifacts"
	"github.com/specops/specops/internal/output"
	"github.com/specops/specops/internal/runstate"
	"github.com/spf13/cobra"
)

func (a *App) newDecisionCommands() []*cobra.Command {
	return []*cobra.Command{
		a.newDecisionsCommand(),
		a.newAcceptCommand(),
		a.newRejectCommand(),
		a.newDeferCommand(),
		a.newAmendCommand(),
	}
}

func (a *App) newDecisionsCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "decisions <run-id>",
		Short: "List proposed decisions for a run",
		Args:  requireArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.repoRoot()
			if err != nil {
				return err
			}
			state, err := runstate.Load(repo, args[0])
			if err != nil {
				return err
			}
			if a.JSON {
				return a.writeJSON(map[string]any{"run_id": state.RunID, "decisions": state.Decisions})
			}
			for _, decision := range state.Decisions {
				a.humanf("%s\t%s\t%s\trecommendation=%s\n", decision.ID, decision.Status, decision.Title, decision.Recommendation)
			}
			return nil
		},
	}
}

func (a *App) newAcceptCommand() *cobra.Command {
	var allRecommended bool
	cmd := &cobra.Command{
		Use:   "accept <run-id> <decision-id>|--all-recommended",
		Short: "Accept a proposed decision",
		Args: func(cmd *cobra.Command, args []string) error {
			if allRecommended && len(args) == 1 {
				return nil
			}
			if len(args) == 2 {
				return nil
			}
			return output.UsageError("accept requires <run-id> and <decision-id>, or <run-id> --all-recommended")
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.repoRoot()
			if err != nil {
				return err
			}
			var state *runstate.RunState
			if allRecommended {
				state, err = artifacts.AcceptRecommended(repo, args[0])
			} else {
				state, err = artifacts.SetDecision(repo, args[0], args[1], "accepted", "")
			}
			if err != nil {
				if errors.Is(err, os.ErrNotExist) {
					return output.OperationalError("decision not found")
				}
				return err
			}
			if a.JSON {
				return a.writeJSON(state)
			}
			a.humanf("decisions updated for %s\n", state.RunID)
			return nil
		},
	}
	cmd.Flags().BoolVar(&allRecommended, "all-recommended", false, "accept all decisions whose recommendation is accept")
	return cmd
}

func (a *App) newRejectCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "reject <run-id> <decision-id>",
		Short: "Reject a proposed decision",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return output.UsageError("reject requires <run-id> and <decision-id>")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.repoRoot()
			if err != nil {
				return err
			}
			state, err := artifacts.SetDecision(repo, args[0], args[1], "rejected", "")
			if err != nil {
				return err
			}
			if a.JSON {
				return a.writeJSON(state)
			}
			a.humanf("rejected %s\n", args[1])
			return nil
		},
	}
}

func (a *App) newDeferCommand() *cobra.Command {
	var reason string
	cmd := &cobra.Command{
		Use:   "defer <run-id> <decision-id> --reason <text>",
		Short: "Defer a proposed decision",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return output.UsageError("defer requires <run-id> and <decision-id>")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.repoRoot()
			if err != nil {
				return err
			}
			state, err := artifacts.SetDecision(repo, args[0], args[1], "deferred", reason)
			if err != nil {
				return err
			}
			if a.JSON {
				return a.writeJSON(state)
			}
			a.humanf("deferred %s\n", args[1])
			return nil
		},
	}
	cmd.Flags().StringVar(&reason, "reason", "", "reason for deferral")
	_ = cmd.MarkFlagRequired("reason")
	return cmd
}

func (a *App) newAmendCommand() *cobra.Command {
	var text string
	cmd := &cobra.Command{
		Use:   "amend <run-id> <decision-id> --text <file-or-inline>",
		Short: "Amend a proposed decision",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return output.UsageError("amend requires <run-id> and <decision-id>")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.repoRoot()
			if err != nil {
				return err
			}
			content := text
			if raw, err := os.ReadFile(text); err == nil {
				content = string(raw)
			}
			state, err := artifacts.SetDecision(repo, args[0], args[1], "amended", content)
			if err != nil {
				return err
			}
			if a.JSON {
				return a.writeJSON(state)
			}
			a.humanf("amended %s\n", args[1])
			return nil
		},
	}
	cmd.Flags().StringVar(&text, "text", "", "amendment text or file path")
	_ = cmd.MarkFlagRequired("text")
	return cmd
}
