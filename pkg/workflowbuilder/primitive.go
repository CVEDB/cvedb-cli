package workflowbuilder

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cvedb/cvedb-cli/pkg/cvedb"
)

func GetLabeledPrimitiveNodes(wfVersion *cvedb.WorkflowVersion) ([]*cvedb.PrimitiveNode, error) {
	if wfVersion == nil {
		return nil, fmt.Errorf("workflow version is nil")
	}

	var labeledNodes []*cvedb.PrimitiveNode
	for _, node := range wfVersion.Data.PrimitiveNodes {
		if node.Label != node.Value {
			labeledNodes = append(labeledNodes, node)
		}
	}

	return labeledNodes, nil
}

// addPrimitiveNode adds a new primitive node or updates an existing one
func addPrimitiveNode(wfVersion *cvedb.WorkflowVersion, nodeType string, value any) (*cvedb.PrimitiveNode, error) {
	if wfVersion == nil {
		return nil, fmt.Errorf("workflow version is nil")
	}

	node, err := createPrimitiveNode(wfVersion, nodeType, value)
	if err != nil {
		return nil, fmt.Errorf("failed to create primitive node: %w", err)
	}

	if wfVersion.Data.PrimitiveNodes == nil {
		wfVersion.Data.PrimitiveNodes = make(map[string]*cvedb.PrimitiveNode)
	}

	wfVersion.Data.PrimitiveNodes[node.Name] = node
	return node, nil
}

// createPrimitiveNode creates a new primitive node based on its type and value
func createPrimitiveNode(wfVersion *cvedb.WorkflowVersion, nodeType string, value any) (*cvedb.PrimitiveNode, error) {
	var name string
	var typeName string
	switch nodeType {
	case "STRING":
		name = fmt.Sprintf("string-input-%d", getAvailablePrimitiveNodeID(nodeType, wfVersion))
		typeName = "STRING"
	case "BOOLEAN":
		name = fmt.Sprintf("boolean-input-%d", getAvailablePrimitiveNodeID(nodeType, wfVersion))
		typeName = "BOOLEAN"
	case "FILE":
		name = fmt.Sprintf("http-input-%d", getAvailablePrimitiveNodeID(nodeType, wfVersion))
		typeName = "URL"
	case "FOLDER":
		name = fmt.Sprintf("git-input-%d", getAvailablePrimitiveNodeID(nodeType, wfVersion))
		typeName = "GIT"
	default:
		return nil, fmt.Errorf("unsupported node type: %s", nodeType)
	}

	normalizedValue, label, err := processPrimitiveNodeValue(nodeType, value)
	if err != nil {
		return nil, err
	}

	pNode := &cvedb.PrimitiveNode{
		Name:     name,
		Type:     nodeType,
		TypeName: typeName,
		Value:    normalizedValue,
		Label:    label,
	}

	return pNode, nil
}

// setPrimitiveNodeValue updates the value of a primitive node through the provided pointer
func setPrimitiveNodeValue(pNode *cvedb.PrimitiveNode, value any) error {
	if pNode == nil {
		return fmt.Errorf("primitive node cannot be nil")
	}

	normalizedValue, impliedLabel, err := processPrimitiveNodeValue(pNode.Type, value)
	if err != nil {
		return err
	}

	// Only update the label if it matches the current value, indicating it's using the default behavior
	// This preserves custom labels that were manually set to user-friendly names
	if pNode.Label == pNode.Value {
		pNode.Label = impliedLabel
	}
	pNode.Value = normalizedValue

	return nil
}

// processPrimitiveNodeValue normalizes and validates primitive node values based on type and returns the normalized value, a display label, and any validation errors.
func processPrimitiveNodeValue(nodeType string, value any) (any, string, error) {
	var normalizedValue any
	var label string

	switch val := value.(type) {
	case string:
		switch nodeType {
		case "STRING":
			normalizedValue = val
		case "FILE":
			if !strings.HasPrefix(val, "https://") && !strings.HasPrefix(val, "http://") && !strings.HasPrefix(val, "cvedb://file/") && !strings.HasPrefix(val, "cvedb://output/") {
				return nil, "", fmt.Errorf("file input must be a valid URL (http:// or https://) for a remote file, cvedb://file/path for stored files, or cvedb://output/id for workflow outputs")
			}
			normalizedValue = val
		case "FOLDER":
			if !strings.HasPrefix(val, "http://") && !strings.HasPrefix(val, "https://") {
				return nil, "", fmt.Errorf("folder input must be a valid git repository URL")
			}
			normalizedValue = val
		case "BOOLEAN":
			normalizedValue = val
		}
	case int:
		normalizedValue = strconv.Itoa(val)
	case bool:
		normalizedValue = val
	default:
		return nil, "", fmt.Errorf("unsupported value type: %T; only string, int, and bool are valid for primitive nodes", value)
	}

	if nodeType == "BOOLEAN" {
		label = strconv.FormatBool(normalizedValue.(bool))
	} else {
		label = normalizedValue.(string)
	}

	return normalizedValue, label, nil
}

// getAvailablePrimitiveNodeID determines the next available ID for a primitive node type
func getAvailablePrimitiveNodeID(nodeType string, wfVersion *cvedb.WorkflowVersion) int {
	availableID := 1

	if wfVersion == nil || wfVersion.Data.PrimitiveNodes == nil {
		return availableID
	}

	var prefix string
	switch nodeType {
	case "STRING":
		prefix = "string-input-"
	case "BOOLEAN":
		prefix = "boolean-input-"
	case "FILE":
		prefix = "http-input-"
	case "FOLDER":
		prefix = "git-input-"
	}
	for nodeName := range wfVersion.Data.PrimitiveNodes {
		if strings.HasPrefix(nodeName, prefix) {
			currentID, _ := strconv.Atoi(strings.TrimPrefix(nodeName, prefix))
			if currentID >= availableID {
				availableID = currentID + 1
			}
		}
	}
	return availableID
}

// removePrimitiveNode removes a primitive node from the workflow version
func removePrimitiveNode(wfVersion *cvedb.WorkflowVersion, pNodeID string) error {
	if wfVersion == nil {
		return fmt.Errorf("workflow version is nil")
	}

	delete(wfVersion.Data.PrimitiveNodes, pNodeID)
	return nil
}
