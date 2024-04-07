package library

import (
	"fmt"
	"math"

	"github.com/spf13/cobra"
	"github.com/cvedb/cvedb-cli/cmd/list"
)

// libraryListToolsCmd represents the libraryListTools command
var libraryListToolsCmd = &cobra.Command{
	Use:   "tools",
	Short: "List tools from the Cvedb library",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		tools := list.GetTools(math.MaxInt, "", "")
		if len(tools) > 0 {
			PrintTools(tools, jsonOutput)
		} else {
			fmt.Println("Couldn't find any tool in the library!")
		}
	},
}

func init() {
	libraryListCmd.AddCommand(libraryListToolsCmd)
	libraryListToolsCmd.Flags().BoolVar(&jsonOutput, "json", false, "Display output in JSON format")
}
