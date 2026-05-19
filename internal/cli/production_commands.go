package cli

import (
	"github.com/specops/specops/internal/artifacts"
	"github.com/spf13/cobra"
)

func (a *App) newProductionCommands() []*cobra.Command {
	return []*cobra.Command{
		a.newIntakeCommand(),
		a.newRefineCommand(),
		a.newHardenCommand(),
		a.newSynthesizeCommand(),
		a.newSupersedeSynthesisCommand(),
		a.newDeepenCommand(),
	}
}

func (a *App) newIntakeCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "intake <run-id>",
		Short: "Produce an intake artifact for a run",
		Args:  requireArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.repoRoot()
			if err != nil {
				return err
			}
			result, err := artifacts.Intake(repo, args[0])
			if err != nil {
				return err
			}
			if a.JSON {
				return a.writeJSON(result)
			}
			a.humanf("intake complete for %s\n", result.RunID)
			return nil
		},
	}
}

func (a *App) newRefineCommand() *cobra.Command {
	var from string
	cmd := &cobra.Command{
		Use:   "refine <run-id> --from <file>",
		Short: "Copy an authored refined artifact into the run",
		Args:  requireArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.repoRoot()
			if err != nil {
				return err
			}
			result, err := artifacts.RefineFrom(repo, args[0], from)
			if err != nil {
				return err
			}
			if a.JSON {
				return a.writeJSON(result)
			}
			a.humanf("refined %s\n", result.RunID)
			return nil
		},
	}
	cmd.Flags().StringVar(&from, "from", "", "required agent or human-authored refined artifact")
	return cmd
}

func (a *App) newHardenCommand() *cobra.Command {
	var backend string
	var from string
	cmd := &cobra.Command{
		Use:   "harden <run-id> --from <file> [--backend convo-relay]",
		Short: "Copy an authored hardening artifact into the run",
		Args:  requireArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if backend == "" {
				backend = a.Backend
			}
			repo, err := a.repoRoot()
			if err != nil {
				return err
			}
			result, err := artifacts.HardenFrom(repo, args[0], backend, from)
			if err != nil {
				return err
			}
			if a.JSON {
				return a.writeJSON(result)
			}
			a.humanf("hardened %s\n", result.RunID)
			return nil
		},
	}
	cmd.Flags().StringVar(&backend, "backend", "", "backend used for hardening")
	cmd.Flags().StringVar(&from, "from", "", "required agent or human-authored hardened artifact")
	return cmd
}

func (a *App) newSynthesizeCommand() *cobra.Command {
	var from string
	cmd := &cobra.Command{
		Use:   "synthesize <run-id> --from <spec_delta.json>",
		Short: "Copy an authored typed spec delta into the run",
		Args:  requireArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.repoRoot()
			if err != nil {
				return err
			}
			result, err := artifacts.SynthesizeFrom(repo, args[0], from)
			if err != nil {
				return err
			}
			if a.JSON {
				return a.writeJSON(result)
			}
			a.humanf("synthesized spec delta for %s\n", result.RunID)
			return nil
		},
	}
	cmd.Flags().StringVar(&from, "from", "", "required agent or human-authored spec_delta.json")
	return cmd
}

func (a *App) newSupersedeSynthesisCommand() *cobra.Command {
	var from string
	var reopenDecisions bool
	cmd := &cobra.Command{
		Use:   "supersede-synthesis <run-id> --from <spec_delta.json> [--reopen-decisions]",
		Short: "Replace a planned run's synthesized spec delta before apply",
		Args:  requireArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.repoRoot()
			if err != nil {
				return err
			}
			result, err := artifacts.SupersedeSynthesisFrom(repo, args[0], from, reopenDecisions)
			if err != nil {
				return err
			}
			if a.JSON {
				return a.writeJSON(result)
			}
			a.humanf("superseded synthesis for %s; status=%s archived=%d\n", result.RunID, result.Status, len(result.ArchivedArtifacts))
			return nil
		},
	}
	cmd.Flags().StringVar(&from, "from", "", "required replacement spec_delta.json")
	cmd.Flags().BoolVar(&reopenDecisions, "reopen-decisions", false, "allow new or changed decisions and return to the decision gate")
	return cmd
}

func (a *App) newDeepenCommand() *cobra.Command {
	var target string
	cmd := &cobra.Command{
		Use:   "deepen <run-id> --target <concept-or-doc>",
		Short: "Create a focused deepening artifact",
		Args:  requireArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.repoRoot()
			if err != nil {
				return err
			}
			result, err := artifacts.Deepen(repo, args[0], target)
			if err != nil {
				return err
			}
			if a.JSON {
				return a.writeJSON(result)
			}
			a.humanf("deepened %s for %s\n", target, result.RunID)
			return nil
		},
	}
	cmd.Flags().StringVar(&target, "target", "", "concept or document to deepen")
	_ = cmd.MarkFlagRequired("target")
	return cmd
}
