package store

import (
	"cvedb-cli/cmd/list"
	"fmt"
	"math"

	"github.com/spf13/cobra"
)

// storeSearchCmd represents the storeSearch command
var storeSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for workflows and tools in the CVEDB store",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		search := ""
		if len(args) > 0 {
			search = args[0]
		}
		tools := list.GetTools(math.MaxInt, search, "")
		workflows := list.GetWorkflows("", "", search, true)
		if len(tools) > 0 {
			printTools(tools)
		} else {
			fmt.Println("Couldn't find any tool in the store that matches the search!")
		}
		if len(workflows) > 0 {
			printWorkflows(workflows)
		} else {
			fmt.Println("Couldn't find any workflow in the store that matches the search!")
		}
	},
}

func init() {
	StoreCmd.AddCommand(storeSearchCmd)
}
