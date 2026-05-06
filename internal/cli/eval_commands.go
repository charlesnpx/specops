package cli

import (
	evalpkg "github.com/specops/specops/internal/eval"
	"github.com/specops/specops/internal/input"
	"github.com/spf13/cobra"
)

func (a *App) newEvalCommands() []*cobra.Command {
	return []*cobra.Command{
		a.newReproduceCommand(),
		a.newEvalCommand(),
		a.newDiffCommand(),
		a.newScoreCommand(),
	}
}

func (a *App) newReproduceCommand() *cobra.Command {
	var fixture string
	var out string
	cmd := &cobra.Command{
		Use:   "reproduce --fixture <dir> --out <dir>",
		Short: "Reproduce a fixture into an output directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := input.BuildFixture(fixture, out)
			if err != nil {
				return err
			}
			if a.JSON {
				return a.writeJSON(result)
			}
			a.humanf("reproduced fixture to %s\n", result.Out)
			return nil
		},
	}
	cmd.Flags().StringVar(&fixture, "fixture", "", "fixture directory")
	cmd.Flags().StringVar(&out, "out", "", "output directory")
	_ = cmd.MarkFlagRequired("fixture")
	_ = cmd.MarkFlagRequired("out")
	return cmd
}

func (a *App) newEvalCommand() *cobra.Command {
	var gold string
	var candidate string
	cmd := &cobra.Command{
		Use:   "eval --gold <repo> --candidate <repo>",
		Short: "Evaluate a candidate spec repository against gold",
		RunE: func(cmd *cobra.Command, args []string) error {
			report, err := evalpkg.Run(gold, candidate)
			if err != nil {
				return err
			}
			if a.JSON {
				return a.writeJSON(report)
			}
			a.humanf("eval %s\n", report.EvalID)
			for key, value := range report.Scores {
				a.humanf("%s=%.2f\n", key, value)
			}
			for _, finding := range report.Findings {
				a.humanf("finding: %s\n", finding)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&gold, "gold", "", "gold repository")
	cmd.Flags().StringVar(&candidate, "candidate", "", "candidate repository")
	_ = cmd.MarkFlagRequired("gold")
	_ = cmd.MarkFlagRequired("candidate")
	return cmd
}

func (a *App) newDiffCommand() *cobra.Command {
	var gold string
	var candidate string
	cmd := &cobra.Command{
		Use:   "diff --gold <repo> --candidate <repo>",
		Short: "Compare gold and candidate repository file sets",
		RunE: func(cmd *cobra.Command, args []string) error {
			report, err := evalpkg.Diff(gold, candidate)
			if err != nil {
				return err
			}
			if a.JSON {
				return a.writeJSON(report)
			}
			a.humanf("gold-only: %d\ncandidate-only: %d\ncommon: %d\n", len(report.GoldOnly), len(report.CandidateOnly), len(report.Common))
			return nil
		},
	}
	cmd.Flags().StringVar(&gold, "gold", "", "gold repository")
	cmd.Flags().StringVar(&candidate, "candidate", "", "candidate repository")
	_ = cmd.MarkFlagRequired("gold")
	_ = cmd.MarkFlagRequired("candidate")
	return cmd
}

func (a *App) newScoreCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "score <eval-report>",
		Short: "Score an eval report",
		Args:  requireArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			scores, err := evalpkg.Score(args[0])
			if err != nil {
				return err
			}
			if a.JSON {
				return a.writeJSON(scores)
			}
			for key, value := range scores {
				a.humanf("%s=%.2f\n", key, value)
			}
			return nil
		},
	}
}
