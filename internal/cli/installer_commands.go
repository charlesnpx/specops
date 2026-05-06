package cli

import (
	"fmt"

	"github.com/specops/specops/internal/install"
	"github.com/specops/specops/internal/output"
	"github.com/spf13/cobra"
)

func (a *App) newInstallSkillCommand() *cobra.Command {
	var plan bool
	var doInstall bool
	var uninstall bool
	var target string
	var installRoot string
	cmd := &cobra.Command{
		Use:   "install-skill [plan|install|uninstall] --target claude|codex|tools|all --json",
		Short: "Run the mise-en-place delegated installer contract",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) > 1 {
				return output.UsageError("install-skill accepts at most one operation argument")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			op, err := installerOperation(args, plan, doInstall, uninstall)
			if err != nil {
				return err
			}
			if target == "" {
				target = string(install.TargetAll)
			}
			opts := install.Options{Operation: op, Target: install.Target(target), InstallRoot: installRoot, Version: a.Version}
			report, err := install.Execute(opts)
			if err != nil {
				return output.ContractError(err.Error())
			}
			if a.JSON {
				return a.writeJSON(report)
			}
			a.humanf("%s %s\n", report.Name, report.Operation)
			for name, target := range report.Targets {
				a.humanf("%s: %d file(s)\n", name, len(target.Files))
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&plan, "plan", false, "plan installation")
	cmd.Flags().BoolVar(&doInstall, "install", false, "install files")
	cmd.Flags().BoolVar(&uninstall, "uninstall", false, "uninstall files")
	cmd.Flags().StringVar(&target, "target", "all", "target: claude, codex, tools, or all")
	cmd.Flags().StringVar(&installRoot, "install-root", "", "absolute staging root")
	return cmd
}

func installerOperation(args []string, plan, doInstall, uninstall bool) (install.Operation, error) {
	var ops []install.Operation
	if plan {
		ops = append(ops, install.OperationPlan)
	}
	if doInstall {
		ops = append(ops, install.OperationInstall)
	}
	if uninstall {
		ops = append(ops, install.OperationUninstall)
	}
	if len(args) == 1 {
		switch args[0] {
		case "plan":
			ops = append(ops, install.OperationPlan)
		case "install":
			ops = append(ops, install.OperationInstall)
		case "uninstall":
			ops = append(ops, install.OperationUninstall)
		default:
			return "", output.UsageError(fmt.Sprintf("unknown installer operation %q", args[0]))
		}
	}
	if len(ops) != 1 {
		return "", output.UsageError("exactly one of --plan, --install, --uninstall, or operation argument is required")
	}
	return ops[0], nil
}
