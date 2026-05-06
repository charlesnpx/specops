package cli

import (
	"github.com/specops/specops/internal/artifacts"
	"github.com/specops/specops/internal/audit"
	"github.com/spf13/cobra"
)

func (a *App) newCompileCommands() []*cobra.Command {
	return []*cobra.Command{
		a.newCompileCommand(),
		a.newPlanCommand(),
		a.newApplyCommand(),
		a.newAuditCommand(),
	}
}

func (a *App) newCompileCommand() *cobra.Command {
	var acceptedOnly bool
	cmd := &cobra.Command{
		Use:   "compile <run-id> --accepted-only",
		Short: "Compile accepted decisions into a patch plan artifact",
		Args:  requireArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.repoRoot()
			if err != nil {
				return err
			}
			plan, err := artifacts.Compile(repo, args[0], acceptedOnly)
			if err != nil {
				return err
			}
			if a.JSON {
				return a.writeJSON(plan)
			}
			a.humanf("compiled %d patch item(s) for %s\n", len(plan.Items), plan.RunID)
			return nil
		},
	}
	cmd.Flags().BoolVar(&acceptedOnly, "accepted-only", false, "compile only accepted decisions")
	return cmd
}

func (a *App) newPlanCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "plan <run-id>",
		Short: "Show exact file-level patch intent",
		Args:  requireArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.repoRoot()
			if err != nil {
				return err
			}
			plan, err := artifacts.LoadPlan(repo, args[0])
			if err != nil {
				return err
			}
			if _, err := artifacts.MarkPlanned(repo, args[0]); err != nil {
				return err
			}
			if a.JSON {
				return a.writeJSON(plan)
			}
			for _, item := range plan.Items {
				a.humanf("%s\t%s\t%s\n", item.ID, item.Action, item.Path)
			}
			return nil
		},
	}
}

func (a *App) newApplyCommand() *cobra.Command {
	var interactive bool
	var dryRun bool
	var commit bool
	cmd := &cobra.Command{
		Use:   "apply <run-id> [--interactive] [--dry-run] [--commit]",
		Short: "Apply a reviewed patch plan",
		Args:  requireArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_ = interactive
			repo, err := a.repoRoot()
			if err != nil {
				return err
			}
			result, err := artifacts.Apply(repo, args[0], dryRun, commit)
			if err != nil {
				return err
			}
			if a.JSON {
				return a.writeJSON(result)
			}
			for _, file := range result.Files {
				a.humanf("%s\t%s\n", file.Status, file.Path)
			}
			for _, warning := range result.Warnings {
				a.humanf("warning: %s\n", warning)
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&interactive, "interactive", false, "prompt for patch decisions where supported")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "show what would be written")
	cmd.Flags().BoolVar(&commit, "commit", false, "commit applied files when the repo is a git worktree")
	return cmd
}

func (a *App) newAuditCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "audit",
		Short: "Run deterministic SpecOps audit checks",
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.repoRoot()
			if err != nil {
				return err
			}
			report, err := audit.Run(repo)
			if err != nil {
				return err
			}
			if a.JSON {
				return a.writeJSON(report)
			}
			for _, check := range report.Checks {
				a.humanf("%s\t%v\n", check.Name, check.Passed)
			}
			for _, warning := range report.Warnings {
				a.humanf("warning: %s\n", warning)
			}
			return nil
		},
	}
}
