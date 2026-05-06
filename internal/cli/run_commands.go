package cli

import (
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
			return nil
		},
	}
}
