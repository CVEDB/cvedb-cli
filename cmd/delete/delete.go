package delete

import (
	"fmt"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/cvedb/cvedb-cli/client/request"
	"github.com/cvedb/cvedb-cli/util"

	"github.com/spf13/cobra"
)

// DeleteCmd represents the delete command
var DeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes an object on the Cvedb platform",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		space, project, workflow, found := util.GetObjects(args)
		if !found {
			return
		}

		if workflow != nil {
			deleteWorkflow(workflow.ID)
		} else if project != nil {
			DeleteProject(project.ID)
		} else if space != nil {
			deleteSpace("", space.ID)
		}
	},
}

func deleteSpace(name string, id uuid.UUID) {
	if id == uuid.Nil {
		space := util.GetSpaceByName(name)
		if space == nil {
			fmt.Println("Couldn't find space named " + name + "!")
			os.Exit(0)
		}
		id = space.ID
	}

	resp := request.Cvedb.Delete().DoF("spaces/%s/", id.String())
	if resp == nil {
		fmt.Println("Couldn't delete space with ID: " + id.String())
		os.Exit(0)
	}

	if resp.Status() != http.StatusNoContent {
		request.ProcessUnexpectedResponse(resp)
	} else {
		fmt.Println("Space deleted successfully!")
	}
}

func DeleteProject(id uuid.UUID) {
	resp := request.Cvedb.Delete().DoF("projects/%s/", id.String())
	if resp == nil {
		fmt.Println("Couldn't delete project with ID: " + id.String())
		os.Exit(0)
	}

	if resp.Status() != http.StatusNoContent {
		request.ProcessUnexpectedResponse(resp)
	} else {
		fmt.Println("Project deleted successfully!")
	}
}

func deleteWorkflow(id uuid.UUID) {
	resp := request.Cvedb.Delete().DoF("workflow/%s/", id.String())
	if resp == nil {
		fmt.Println("Couldn't delete workflow with ID: " + id.String())
		os.Exit(0)
	}

	if resp.Status() != http.StatusNoContent {
		request.ProcessUnexpectedResponse(resp)
	} else {
		fmt.Println("Workflow deleted successfully!")
	}
}
