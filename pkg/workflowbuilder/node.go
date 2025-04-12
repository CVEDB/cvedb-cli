package workflowbuilder

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cvedb/cvedb-cli/pkg/cvedb"
)

func GetLabeledNodes(wfVersion *cvedb.WorkflowVersion) ([]*cvedb.Node, error) {
	if wfVersion == nil {
		return nil, fmt.Errorf("workflow version is nil")
	}

	var labeledNodes []*cvedb.Node
	for _, node := range wfVersion.Data.Nodes {
		if !isDefaultLabel(node.Meta.Label, node.Name) {
			labeledNodes = append(labeledNodes, node)
		}
	}

	return labeledNodes, nil
}

func isDefaultLabel(label, name string) bool {
	parts := strings.Split(name, "-")
	if len(parts) < 2 {
		return false
	}

	_, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		return false
	}

	defaultLabel := strings.Join(parts[:len(parts)-1], "-")
	return label == defaultLabel
}
