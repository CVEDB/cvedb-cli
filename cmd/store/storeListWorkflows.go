package store

import (
	"cvedb-cli/cmd/list"
	"cvedb-cli/types"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xlab/treeprint"
)

// storeListWorkflowsCmd represents the storeListWorkflows command
var storeListWorkflowsCmd = &cobra.Command{
	Use:   "workflows",
	Short: "List workflows from the CVEDB store",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		workflows := list.GetWorkflows("", "", "", true)
		if len(workflows) > 0 {
			printWorkflows(workflows)
		} else {
			fmt.Println("Couldn't find any workflow in the store!")
		}
	},
}

func init() {
	storeListCmd.AddCommand(storeListWorkflowsCmd)
}

func printWorkflows(workflows []types.WorkflowListResponse) {
	tree := treeprint.New()
	tree.SetValue("Workflows")
	for _, workflow := range workflows {
		wfSubBranch := tree.AddBranch("\U0001f9be " + workflow.Name) //🦾
		if workflow.Description != "" {
			wfSubBranch.AddNode("\U0001f4cb \033[3m" + workflow.Description + "\033[0m") //📋
		}
	}

	fmt.Println(tree.String())
}
