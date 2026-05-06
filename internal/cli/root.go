package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/specops/specops/internal/output"
	"github.com/spf13/cobra"
)

type App struct {
	Repo    string
	Config  string
	JSON    bool
	Quiet   bool
	Verbose bool
	Backend string
	Version string
	Out     io.Writer
	Err     io.Writer
}

func NewRoot(out, err io.Writer, version string) *cobra.Command {
	app := &App{Out: out, Err: err, Version: version}
	root := &cobra.Command{
		Use:           "specops",
		Short:         "SpecOps specification-production CLI",
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.SetOut(out)
	root.SetErr(err)
	root.PersistentFlags().StringVar(&app.Repo, "repo", "", "target repository path")
	root.PersistentFlags().StringVar(&app.Config, "config", "", "config file path")
	root.PersistentFlags().BoolVar(&app.JSON, "json", false, "emit machine-readable JSON")
	root.PersistentFlags().BoolVar(&app.Quiet, "quiet", false, "suppress human output")
	root.PersistentFlags().BoolVar(&app.Verbose, "verbose", false, "emit verbose diagnostics to stderr")
	root.PersistentFlags().StringVar(&app.Backend, "backend", "manual", "agent backend")
	setupCommands := []*cobra.Command{
		app.newInitCommand(),
		app.newDoctorCommand(),
		app.newUpgradeCommand(),
		app.newConfigCommand(),
		app.newRunCommand(),
		app.newContextCommand(),
		app.newNoteCommand(),
		app.newNextCommand(),
	}
	setupCommands = append(setupCommands, app.newInputCommands()...)
	root.AddCommand(setupCommands...)
	root.AddCommand(
		app.newProductionCommands()...,
	)
	root.AddCommand(
		app.newDecisionCommands()...,
	)
	root.AddCommand(
		app.newCompileCommands()...,
	)
	root.AddCommand(
		app.newEvalCommands()...,
	)
	root.AddCommand(app.newInstallSkillCommand())
	return root
}

func (a *App) repoRoot() (string, error) {
	if a.Repo != "" {
		return filepath.Abs(a.Repo)
	}
	return os.Getwd()
}

func (a *App) writeJSON(value any) error {
	return output.WriteJSON(a.Out, value)
}

func (a *App) humanf(format string, args ...any) {
	if a.JSON || a.Quiet {
		return
	}
	fmt.Fprintf(a.Out, format, args...)
}

func requireArgs(count int) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) != count {
			return output.UsageError(fmt.Sprintf("expected %d argument(s), got %d", count, len(args)))
		}
		return nil
	}
}
