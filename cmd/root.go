package cmd

import (
	"log"

	"github.com/cvedb/cvedb-cli/cmd/create"
	"github.com/cvedb/cvedb-cli/cmd/delete"
	"github.com/cvedb/cvedb-cli/cmd/execute"
	"github.com/cvedb/cvedb-cli/cmd/files"
	"github.com/cvedb/cvedb-cli/cmd/get"
	"github.com/cvedb/cvedb-cli/cmd/help"
	"github.com/cvedb/cvedb-cli/cmd/investigate"
	"github.com/cvedb/cvedb-cli/cmd/library"
	"github.com/cvedb/cvedb-cli/cmd/list"
	"github.com/cvedb/cvedb-cli/cmd/output"
	"github.com/cvedb/cvedb-cli/cmd/scripts"
	"github.com/cvedb/cvedb-cli/cmd/stop"
	"github.com/cvedb/cvedb-cli/cmd/tools"
	"github.com/cvedb/cvedb-cli/pkg/version"
	"github.com/cvedb/cvedb-cli/util"

	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "cvedb",
	Short: "Cvedb client for platform access from your local machine",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
	Version: version.Version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	log.SetFlags(0)
	cobra.CheckErr(RootCmd.Execute())
}

func init() {
	RootCmd.PersistentFlags().StringVar(&util.Cfg.User.Token, "token", "", "Cvedb authentication token")
	RootCmd.PersistentFlags().StringVar(&util.Cfg.User.TokenFilePath, "token-file", "", "Cvedb authentication token file")
	RootCmd.PersistentFlags().StringVar(&util.SpaceName, "space", "", "Space name")
	RootCmd.PersistentFlags().StringVar(&util.ProjectName, "project", "", "Project name")
	RootCmd.PersistentFlags().StringVar(&util.WorkflowName, "workflow", "", "Workflow name")
	RootCmd.PersistentFlags().StringVar(&util.URL, "url", "", "URL for referencing a workflow, project, or space")
	RootCmd.PersistentFlags().StringVar(&util.Cfg.Dependency, "node-dependency", "", "This flag doesn't affect the execution logic of the CLI in any way and is intended for controlling node execution order on the Cvedb platform only.")
	RootCmd.PersistentFlags().StringVar(&util.Cfg.BaseUrl, "api-endpoint", "https://cvedb.github.io/api", "The base Cvedb platform API endpoint.")

	RootCmd.AddCommand(list.ListCmd)
	RootCmd.AddCommand(library.LibraryCmd)
	RootCmd.AddCommand(create.CreateCmd)
	RootCmd.AddCommand(delete.DeleteCmd)
	RootCmd.AddCommand(output.OutputCmd)
	RootCmd.AddCommand(execute.ExecuteCmd)
	RootCmd.AddCommand(get.GetCmd)
	RootCmd.AddCommand(files.FilesCmd)
	RootCmd.AddCommand(tools.ToolsCmd)
	RootCmd.AddCommand(scripts.ScriptsCmd)
	RootCmd.AddCommand(stop.StopCmd)
	RootCmd.AddCommand(help.HelpCmd)
	RootCmd.AddCommand(investigate.InvestigateCmd)

	RootCmd.SetVersionTemplate(`{{printf "Cvedb CLI %s\n" .Version}}`)
}
