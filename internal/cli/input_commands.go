package cli

import (
	"os"
	"path/filepath"

	"github.com/specops/specops/internal/input"
	"github.com/spf13/cobra"
)

func (a *App) newInputCommands() []*cobra.Command {
	return []*cobra.Command{
		a.newIngestFileCommand(),
		a.newIngestChatCommand(),
		a.newIngestRelayCommand(),
		a.newMineTraceCommand(),
		a.newFixtureBuildCommand(),
	}
}

func (a *App) newIngestFileCommand() *cobra.Command {
	var runID string
	cmd := &cobra.Command{
		Use:   "ingest-file <path>",
		Short: "Ingest a Markdown or text file into a run",
		Args:  requireArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.repoRoot()
			if err != nil {
				return err
			}
			result, err := input.IngestFile(input.IngestOptions{Repo: repo, RunID: runID, Path: args[0], Type: "raw_markdown"})
			if err != nil {
				return err
			}
			if a.JSON {
				return a.writeJSON(result)
			}
			a.humanf("ingested %s into %s\n", args[0], result.RunID)
			return nil
		},
	}
	cmd.Flags().StringVar(&runID, "run", "", "run id")
	_ = cmd.MarkFlagRequired("run")
	return cmd
}

func (a *App) newIngestChatCommand() *cobra.Command {
	var runID string
	var slice bool
	cmd := &cobra.Command{
		Use:   "ingest-chat <path>",
		Short: "Ingest a conversation transcript into a run",
		Args:  requireArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.repoRoot()
			if err != nil {
				return err
			}
			result, err := input.IngestFile(input.IngestOptions{Repo: repo, RunID: runID, Path: args[0], Type: "conversation_transcript", Slice: slice})
			if err != nil {
				return err
			}
			if a.JSON {
				return a.writeJSON(result)
			}
			a.humanf("ingested chat %s into %s\n", args[0], result.RunID)
			return nil
		},
	}
	cmd.Flags().StringVar(&runID, "run", "", "run id")
	cmd.Flags().BoolVar(&slice, "slice", false, "slice transcript into conversation segments")
	_ = cmd.MarkFlagRequired("run")
	return cmd
}

func (a *App) newIngestRelayCommand() *cobra.Command {
	var runID string
	cmd := &cobra.Command{
		Use:   "ingest-relay <path>",
		Short: "Ingest a relay transcript into a run",
		Args:  requireArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			repo, err := a.repoRoot()
			if err != nil {
				return err
			}
			result, err := input.IngestFile(input.IngestOptions{Repo: repo, RunID: runID, Path: args[0], Type: "relay_transcript"})
			if err != nil {
				return err
			}
			if a.JSON {
				return a.writeJSON(result)
			}
			a.humanf("ingested relay %s into %s\n", args[0], result.RunID)
			return nil
		},
	}
	cmd.Flags().StringVar(&runID, "run", "", "run id")
	_ = cmd.MarkFlagRequired("run")
	return cmd
}

func (a *App) newMineTraceCommand() *cobra.Command {
	var inputPath string
	var gold string
	cmd := &cobra.Command{
		Use:   "mine-trace --input <path> [--gold <repo>]",
		Short: "Create a lightweight trace inventory",
		RunE: func(cmd *cobra.Command, args []string) error {
			abs, err := filepath.Abs(inputPath)
			if err != nil {
				return err
			}
			info, err := os.Stat(abs)
			if err != nil {
				return err
			}
			result := map[string]any{"input": abs, "size": info.Size(), "gold": gold}
			if a.JSON {
				return a.writeJSON(result)
			}
			a.humanf("trace input: %s (%d bytes)\n", abs, info.Size())
			if gold != "" {
				a.humanf("gold repo: %s\n", gold)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&inputPath, "input", "", "input trace path")
	cmd.Flags().StringVar(&gold, "gold", "", "gold repository path")
	_ = cmd.MarkFlagRequired("input")
	return cmd
}

func (a *App) newFixtureBuildCommand() *cobra.Command {
	var from string
	var out string
	cmd := &cobra.Command{
		Use:   "fixture-build --from <path> --out <dir>",
		Short: "Build a clean reproduction fixture",
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := input.BuildFixture(from, out)
			if err != nil {
				return err
			}
			if a.JSON {
				return a.writeJSON(result)
			}
			a.humanf("fixture written to %s (%d files)\n", result.Out, len(result.Files))
			return nil
		},
	}
	cmd.Flags().StringVar(&from, "from", "", "source path")
	cmd.Flags().StringVar(&out, "out", "", "output directory")
	_ = cmd.MarkFlagRequired("from")
	_ = cmd.MarkFlagRequired("out")
	return cmd
}
