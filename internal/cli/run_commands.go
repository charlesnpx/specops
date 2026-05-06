package cli

import (
	"github.com/specops/specops/internal/artifacts"
	"github.com/specops/specops/internal/runstate"
	"github.com/spf13/cobra"
)

func (a *App) newRunCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "run", Short: "Manage SpecOps runs"}
	cmd.AddCommand(a.newRunNewCommand())
	cmd.AddCommand(a.newRunListCommand())
	cmd.AddCommand(a.newRunShowCommand())
	cmd.AddCommand(a.newRunStatusCommand())
	return cmd
}

func (a *App) newRunNewCommand() *cobra.Command {
	var name string
	cmd := &cobra.Command{
		Use:   "new --name <name>",
		Short: "Create a new run",
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.repoRoot()
			if err != nil {
				return err
			}
			state, err := runstate.NewRun(repo, name)
			if err != nil {
				return err
			}
			if a.JSON {
				return a.writeJSON(state)
			}
			a.humanf("%s\n", state.RunID)
			return nil
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "human-readable run name")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}

func (a *App) newRunListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List runs",
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.repoRoot()
			if err != nil {
				return err
			}
			runs, err := runstate.List(repo)
			if err != nil {
				return err
			}
			if a.JSON {
				return a.writeJSON(runs)
			}
			for _, run := range runs {
				a.humanf("%s\t%s\t%s\n", run.RunID, run.Status, run.Name)
			}
			return nil
		},
	}
}

func (a *App) newRunShowCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "show <run-id>",
		Short: "Show run state",
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
				return a.writeJSON(state)
			}
			a.humanf("run: %s\nstatus: %s\nname: %s\n", state.RunID, state.Status, state.Name)
			return nil
		},
	}
}

func (a *App) newRunStatusCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "status <run-id>",
		Short: "Show a run status",
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
				return a.writeJSON(map[string]any{"run_id": state.RunID, "status": state.Status, "next": state.Next})
			}
			a.humanf("%s\n", state.Status)
			return nil
		},
	}
}

func (a *App) newNextCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "next <run-id>",
		Short: "Recommend the next legal step for a run",
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
			next := runstate.NextForStatus(state.Status, state.RunID)
			if a.JSON {
				return a.writeJSON(map[string]any{"run_id": state.RunID, "status": state.Status, "next": next})
			}
			if next == nil {
				a.humanf("no next step for %s\n", state.Status)
				return nil
			}
			a.humanf("%s\n%s\n", next.Command, next.Reason)
			if next.GateKind == "semantic" && next.NoteCommand != "" {
				a.humanf("semantic commands require a stage note before execution\n%s\n", next.NoteCommand)
			}
			return nil
		},
	}
}

func (a *App) newContextCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "context <run-id>",
		Short: "Show compiled run context without mutating state",
		Args:  requireArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.repoRoot()
			if err != nil {
				return err
			}
			context, err := artifacts.Context(repo, args[0])
			if err != nil {
				return err
			}
			if a.JSON {
				return a.writeJSON(context)
			}
			a.humanf("run: %s\nstatus: %s\n", context.RunID, context.Status)
			if context.SourceSummary != "" {
				a.humanf("\nsource summary:\n%s\n", context.SourceSummary)
			}
			if context.NextGate != nil {
				a.humanf("\nnext gate: %s (%s)\n%s\n%s\n", context.NextGate.Stage, context.NextGate.GateKind, context.NextGate.Command, context.NextGate.Reason)
				if context.NextGate.GateKind == "semantic" && context.NextGate.NoteCommand != "" {
					a.humanf("semantic commands require a stage note before execution\nnote command: %s\n", context.NextGate.NoteCommand)
				}
			}
			if len(context.OperatorGuidance.SuggestedQuestions) > 0 {
				a.humanf("\nsuggested questions:\n")
				for _, question := range context.OperatorGuidance.SuggestedQuestions {
					a.humanf("- %s\n", question)
				}
				a.humanf("- %s\n", context.OperatorGuidance.ControlQuestion)
			}
			if len(context.Artifacts) > 0 {
				a.humanf("\nartifacts:\n")
				for _, artifact := range context.Artifacts {
					a.humanf("- %s\t%s\n", artifact.Type, artifact.Path)
				}
			}
			if len(context.Decisions) > 0 {
				a.humanf("\ndecisions:\n")
				for _, decision := range context.Decisions {
					a.humanf("- %s\t%s\t%s\n", decision.ID, decision.Status, decision.Title)
				}
			}
			if context.PatchPlan != nil {
				a.humanf("\npatch plan: %d item(s)\n", len(context.PatchPlan.Items))
			}
			return nil
		},
	}
}

func (a *App) newNoteCommand() *cobra.Command {
	var stage string
	var text string
	cmd := &cobra.Command{
		Use:   "note <run-id> --stage <stage> --text <file-or-inline>",
		Short: "Record operator guidance for a run without advancing state",
		Args:  requireArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.repoRoot()
			if err != nil {
				return err
			}
			result, err := artifacts.Note(repo, args[0], stage, text)
			if err != nil {
				return err
			}
			if a.JSON {
				return a.writeJSON(result)
			}
			a.humanf("noted %s for %s\n", stage, result.RunID)
			return nil
		},
	}
	cmd.Flags().StringVar(&stage, "stage", "", "stage this note applies to")
	cmd.Flags().StringVar(&text, "text", "", "note text or file path")
	_ = cmd.MarkFlagRequired("stage")
	_ = cmd.MarkFlagRequired("text")
	return cmd
}
