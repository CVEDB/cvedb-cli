package cmd

import (
	"github.com/spf13/cobra"
	"trickest-cli/cmd/create"
	"trickest-cli/cmd/list"
	"trickest-cli/util"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "trickest",
	Short: "Trickest client for platform access from your local machine",
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
	RootCmd.PersistentFlags().StringVar(&util.Cfg.User.Token, "token", "", "Trickest authentication token")

	cobra.OnInitialize(initVaultID)

	RootCmd.AddCommand(list.ListCmd)
	RootCmd.AddCommand(create.CreateCmd)
}

func initVaultID() {
	util.GetVault()
}