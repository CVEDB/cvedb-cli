package list

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/xlab/treeprint"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"trickest-cli/types"
	"trickest-cli/util"
)

// listCmd represents the list command
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists objects on the Trickest platform",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			spaces := getSpaces("")

			if spaces != nil && len(spaces) > 0 {
				printSpaces(spaces)
			} else {
				fmt.Println("Couldn't find any spaces!")
			}
			return
		}

		pathSplit := strings.Split(strings.Trim(args[0], "/"), "/")
		if len(pathSplit) > 3 {
			fmt.Println("Invalid object path!")
			return
		}
		space := getSpaceByName(pathSplit[0])
		if space == nil {
			fmt.Println("Couldn't find space with the given name!\n" + pathSplit[0])
			return
		}

		if len(pathSplit) == 1 {
			printSpaceDetailed(*space)
			return
		}

		var workflows []types.WorkflowListResponse
		if space.Projects != nil && len(space.Projects) > 0 {
			for _, proj := range space.Projects {
				if proj.Name == pathSplit[1] {
					workflows = getWorkflows(proj.ID)
					if len(pathSplit) == 2 {
						printProject(proj, workflows)
						return
					} else {
						break
					}
				}
			}
		}

		if space.Workflows != nil && len(space.Workflows) > 0 {
			for _, wf := range space.Workflows {
				if wf.Name == pathSplit[1] {
					printWorkflow(wf)
					return
				}
			}
		}

		if len(pathSplit) == 2 {
			fmt.Println("Couldn't find project or workflow named " + pathSplit[1] + " inside " +
				pathSplit[0] + " space!")
		}

		if workflows != nil && len(workflows) > 0 {
			for _, wf := range workflows {
				if wf.Name == pathSplit[2] {
					fullWorkflow := getWorkflowByID(wf.ID)
					printWorkflow(*fullWorkflow)
					return
				}
			}
		} else {
			fmt.Println("No workflows found in " + pathSplit[0] + "/" + pathSplit[1])
		}

		fmt.Println("Workflow named " + pathSplit[2] + " doesn't exist in " + pathSplit[0] + "/" + pathSplit[1] + "/")
	},
}

func init() {

}

func printWorkflow(workflow types.Workflow) {
	tree := treeprint.New()
	tree.SetValue("\U0001f9be " + workflow.Name) //🦾
	if workflow.Description != "" {
		tree.AddNode("\U0001f4cb \033[3m" + workflow.Description + "\033[0m") //📋
	}
	tree.AddNode("Author: " + workflow.Author)
	if len(workflow.Parameters) > 0 {
		branch := tree.AddBranch("Parameters")
		for _, param := range workflow.Parameters {
			paramType := strings.ToLower(param.Type)
			if paramType == "boolean" {
				branch.AddNode("[" + paramType + "] " + strconv.FormatBool(param.Value.(bool)))
			} else {
				branch.AddNode("[" + paramType + "] " + param.Value.(string))
			}
		}
	}

	fmt.Println(tree.String())
}

func printProject(project types.Project, workflows []types.WorkflowListResponse) {
	tree := treeprint.New()
	tree.SetValue("\U0001f5c2  " + project.Name) //🗂
	if project.Description != "" {
		tree.AddNode("\U0001f4cb \033[3m" + project.Description + "\033[0m") //📋
	}
	if workflows != nil && len(workflows) > 0 {
		wfBranch := tree.AddBranch("Workflows")
		for _, workflow := range workflows {
			wfSubBranch := wfBranch.AddBranch("\U0001f9be " + workflow.Name) //🦾
			if workflow.Description != "" {
				wfSubBranch.AddNode("\U0001f4cb \033[3m" + workflow.Description + "\033[0m") //📋
			}
		}
	}

	fmt.Println(tree.String())
}

func printSpaceDetailed(space types.SpaceDetailed) {
	tree := treeprint.New()
	tree.SetValue("\U0001f4c2 " + space.Name) //📂
	if space.Description != "" {
		tree.AddNode("\U0001f4cb \033[3m" + space.Description + "\033[0m") //📋
	}
	if space.Projects != nil && len(space.Projects) > 0 {
		projBranch := tree.AddBranch("Projects")
		for _, proj := range space.Projects {
			projSubBranch := projBranch.AddBranch("\U0001f5c2  " + proj.Name) //🗂
			if proj.Description != "" {
				projSubBranch.AddNode("\U0001f4cb \033[3m" + proj.Description + "\033[0m") //📋
			}
		}
	}
	if space.Workflows != nil && len(space.Workflows) > 0 {
		wfBranch := tree.AddBranch("Workflows")
		for _, workflow := range space.Workflows {
			wfSubBranch := wfBranch.AddBranch("\U0001f9be " + workflow.Name) //🦾
			if workflow.Description != "" {
				wfSubBranch.AddNode("\U0001f4cb \033[3m" + workflow.Description + "\033[0m") //📋
			}
		}
	}
	fmt.Println(tree.String())
}

func printSpaces(spaces []types.Space) {
	tree := treeprint.New()
	tree.SetValue("Spaces")
	for _, space := range spaces {
		branch := tree.AddBranch("\U0001f4c1 " + space.Name) //📂
		if space.Description != "" {
			branch.AddNode("\U0001f4cb \033[3m" + space.Description + "\033[0m") //📋
		}
	}

	fmt.Println(tree.String())
}

func getSpaces(name string) []types.Space {
	urlReq := util.Cfg.BaseUrl + "v1/spaces/?vault=" + util.GetVault()
	urlReq += "&page_size=" + strconv.Itoa(math.MaxInt)

	if name != "" {
		urlReq += "&name=" + url.QueryEscape(name)
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", urlReq, nil)
	req.Header.Add("Authorization", "Token "+util.GetToken())
	req.Header.Add("Accept", "application/json")

	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		fmt.Println("Error: Couldn't get spaces!")
		os.Exit(0)
	}
	defer resp.Body.Close()

	var bodyBytes []byte
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error: Couldn't read spaces response body!")
		os.Exit(0)
	}

	if resp.StatusCode != http.StatusOK {
		util.ProcessUnexpectedResponse(bodyBytes, resp.StatusCode)
	}

	var spaces types.Spaces
	err = json.Unmarshal(bodyBytes, &spaces)
	if err != nil {
		fmt.Println("Error: Couldn't unmarshal spaces response!")
		os.Exit(0)
	}

	return spaces.Results
}

func getSpaceByName(name string) *types.SpaceDetailed {
	spaces := getSpaces(name)
	if spaces == nil || len(spaces) == 0 {
		fmt.Println("Couldn't find space with the given name!\n" + name)
		os.Exit(0)
	}

	return getSpaceByID(spaces[0].ID)
}

func getSpaceByID(id string) *types.SpaceDetailed {
	client := &http.Client{}

	req, err := http.NewRequest("GET", util.Cfg.BaseUrl+"v1/spaces/"+id+"/", nil)
	req.Header.Add("Authorization", "Token "+util.GetToken())
	req.Header.Add("Accept", "application/json")

	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		fmt.Println("Error: Couldn't get space by ID!")
		os.Exit(0)
	}
	defer resp.Body.Close()

	var bodyBytes []byte
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error: Couldn't read space response!")
		os.Exit(0)
	}

	if resp.StatusCode != http.StatusOK {
		util.ProcessUnexpectedResponse(bodyBytes, resp.StatusCode)
	}

	var space types.SpaceDetailed
	err = json.Unmarshal(bodyBytes, &space)
	if err != nil {
		fmt.Println("Error unmarshalling space response!")
		os.Exit(0)
	}

	return &space
}

func getWorkflows(projectID string) []types.WorkflowListResponse {
	urlReq := util.Cfg.BaseUrl + "v1/store/workflow/?vault=" + util.GetVault()
	urlReq += "&page_size=" + strconv.Itoa(math.MaxInt)

	if projectID != "" {
		urlReq = urlReq + "&project=" + projectID
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", urlReq, nil)
	req.Header.Add("Authorization", "Token "+util.GetToken())
	req.Header.Add("Accept", "application/json")

	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		fmt.Println("Error: Couldn't get workflows!")
		os.Exit(0)
	}
	defer resp.Body.Close()

	var bodyBytes []byte
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error: Couldn't read workflows response!")
		os.Exit(0)
	}

	if resp.StatusCode != http.StatusOK {
		util.ProcessUnexpectedResponse(bodyBytes, resp.StatusCode)
	}

	var workflows types.Workflows
	err = json.Unmarshal(bodyBytes, &workflows)
	if err != nil {
		fmt.Println("Error: Couldn't unmarshal workflows response!")
		os.Exit(0)
	}

	return workflows.Results
}

func getWorkflowByID(id string) *types.Workflow {
	client := &http.Client{}

	req, err := http.NewRequest("GET", util.Cfg.BaseUrl+"v1/store/workflow/"+id+"/", nil)
	req.Header.Add("Authorization", "Token "+util.GetToken())
	req.Header.Add("Accept", "application/json")

	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		fmt.Println("Error: Couldn't get workflow.")
		os.Exit(0)
	}
	defer resp.Body.Close()

	var bodyBytes []byte
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error: Couldn't read workflow.")
		os.Exit(0)
	}

	if resp.StatusCode != http.StatusOK {
		util.ProcessUnexpectedResponse(bodyBytes, resp.StatusCode)
	}

	var workflow types.Workflow
	err = json.Unmarshal(bodyBytes, &workflow)
	if err != nil {
		fmt.Println("Error: Couldn't unmarshal workflow response!")
		os.Exit(0)
	}

	return &workflow
}
