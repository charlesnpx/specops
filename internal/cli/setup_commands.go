package cli

import (
	"fmt"

	"github.com/specops/specops/internal/audit"
	"github.com/specops/specops/internal/config"
	"github.com/specops/specops/internal/output"
	"github.com/specops/specops/internal/scaffold"
	"github.com/spf13/cobra"
)

func (a *App) newInitCommand() *cobra.Command {
	var mode string
	var agent string
	var force bool
	var backup bool
	cmd := &cobra.Command{
		Use:   "init [path]",
		Short: "Install a .specops scaffold into a target repository",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) > 1 {
				return output.UsageError("init accepts at most one path")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			path := "."
			if len(args) == 1 {
				path = args[0]
			}
			result, err := scaffold.Init(scaffold.Options{Path: path, Mode: mode, Agent: agent, Force: force, Backup: backup})
			if err != nil {
				return err
			}
			if a.JSON {
				return a.writeJSON(result)
			}
			a.humanf("initialized SpecOps scaffold at %s\n", result.Path)
			for _, warning := range result.Warnings {
				a.humanf("warning: %s\n", warning)
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&mode, "mode", "minimal", "scaffold mode: minimal, vendor, or linked")
	cmd.Flags().StringVar(&agent, "agent", "both", "agent files to write: claude, codex, or both")
	cmd.Flags().BoolVar(&force, "force", false, "overwrite divergent scaffold files")
	cmd.Flags().BoolVar(&backup, "backup", false, "backup divergent scaffold files before overwrite")
	cmd.Flags().Bool("minimal", false, "use minimal scaffold mode")
	cmd.Flags().Bool("vendor", false, "use vendor scaffold mode")
	cmd.Flags().Bool("linked", false, "use linked scaffold mode")
	cmd.PreRun = func(cmd *cobra.Command, args []string) {
		for _, candidate := range []string{"minimal", "vendor", "linked"} {
			if value, _ := cmd.Flags().GetBool(candidate); value {
				mode = candidate
			}
		}
	}
	return cmd
}

func (a *App) newDoctorCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Check the current SpecOps repository",
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
			if report.Passed {
				a.humanf("doctor passed for %s\n", repo)
			} else {
				a.humanf("doctor found issues for %s\n", repo)
			}
			for _, check := range report.Checks {
				a.humanf("- %s: %v\n", check.Name, check.Passed)
			}
			for _, warning := range report.Warnings {
				a.humanf("warning: %s\n", warning)
			}
			return nil
		},
	}
}

func (a *App) newUpgradeCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade the local scaffold if needed",
		RunE: func(cmd *cobra.Command, args []string) error {
			result := map[string]string{"status": "noop", "version": a.Version}
			if a.JSON {
				return a.writeJSON(result)
			}
			a.humanf("scaffold is already compatible with specops %s\n", a.Version)
			return nil
		},
	}
}

func (a *App) newConfigCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "config", Short: "Get, set, or list user configuration"}
	cmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := config.Load(a.Config)
			if err != nil {
				return err
			}
			if a.JSON {
				return a.writeJSON(map[string]any{"path": store.Path, "values": store.Data})
			}
			for key, value := range store.Data {
				a.humanf("%s=%s\n", key, value)
			}
			return nil
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "get <key>",
		Short: "Get a configuration value",
		Args:  requireArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := config.Load(a.Config)
			if err != nil {
				return err
			}
			value, ok := store.Data[args[0]]
			if !ok {
				return output.OperationalError(fmt.Sprintf("config key %q not found", args[0]))
			}
			if a.JSON {
				return a.writeJSON(map[string]string{"key": args[0], "value": value})
			}
			a.humanf("%s\n", value)
			return nil
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set a configuration value",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return output.UsageError("set requires <key> and <value>")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := config.Load(a.Config)
			if err != nil {
				return err
			}
			store.Data[args[0]] = args[1]
			if err := store.Save(); err != nil {
				return err
			}
			if a.JSON {
				return a.writeJSON(map[string]string{"path": store.Path, "key": args[0], "value": args[1]})
			}
			a.humanf("set %s\n", args[0])
			return nil
		},
	})
	return cmd
}
