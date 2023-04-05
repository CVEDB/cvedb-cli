package cmd

import (
	"cvedb-cli/cmd/create"
	"cvedb-cli/cmd/delete"
	"cvedb-cli/cmd/execute"
	"cvedb-cli/cmd/export"
	"cvedb-cli/cmd/get"
	"cvedb-cli/cmd/list"
	"cvedb-cli/cmd/output"
	"cvedb-cli/cmd/store"
	"cvedb-cli/util"

	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "cvedb",
	Short: "CVEDB client for platform access from your local machine",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(RootCmd.Execute())
}

func init() {
	RootCmd.PersistentFlags().StringVar(&util.Cfg.User.Token, "token", "", "CVEDB authentication token")
	RootCmd.PersistentFlags().StringVar(&util.SpaceName, "space", "", "Space name")
	RootCmd.PersistentFlags().StringVar(&util.ProjectName, "project", "", "Project name")
	RootCmd.PersistentFlags().StringVar(&util.WorkflowName, "workflow", "", "Workflow name")

	cobra.OnInitialize(initVaultID)

	RootCmd.AddCommand(list.ListCmd)
	RootCmd.AddCommand(store.StoreCmd)
	RootCmd.AddCommand(create.CreateCmd)
	RootCmd.AddCommand(delete.DeleteCmd)
	RootCmd.AddCommand(output.OutputCmd)
	RootCmd.AddCommand(execute.ExecuteCmd)
	RootCmd.AddCommand(get.GetCmd)
	RootCmd.AddCommand(export.ExportCmd)
}

func initVaultID() {
	util.GetVault()
}
