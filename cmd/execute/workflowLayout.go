package execute

import (
	"sort"
	"trickest-cli/types"
)

func treeHeight(root *types.TreeNode) int {
	if root.Children == nil || len(root.Children) == 0 {
		return 0
	}

	maxHeight := 0
	for _, child := range root.Children {
		newHeight := treeHeight(child)
		if newHeight > maxHeight {
			maxHeight = newHeight
		}
	}
	return maxHeight + 1
}

func adjustChildrenHeight(root *types.TreeNode, nodesPerHeight *map[int][]*types.TreeNode) {
	if root.Parent == nil {
		(*nodesPerHeight)[root.Height] = append((*nodesPerHeight)[root.Height], root)
	}
	if root.Children == nil || len(root.Children) == 0 {
		return
	}
	for _, child := range root.Children {
		child.Height = root.Height - 1
		found := false
		for _, node := range (*nodesPerHeight)[child.Height] {
			if node.NodeName == child.NodeName {
				found = true
				break
			}
		}
		if !found {
			(*nodesPerHeight)[child.Height] = append((*nodesPerHeight)[child.Height], child)
		}
		adjustChildrenHeight(child, nodesPerHeight)
	}
}

func generateNodesCoordinates(version *types.WorkflowVersionDetailed) {
	treesNodes, rootNodes := CreateTrees(version, true)
	for _, node := range treesNodes {
		node.Height = treeHeight(node)
	}

	nodesPerHeight := make(map[int][]*types.TreeNode, 0)
	maxRootHeight := 0
	for _, node := range rootNodes {
		if node.Height > maxRootHeight {
			maxRootHeight = node.Height
		}
	}
	for _, node := range rootNodes {
		adjustChildrenHeight(node, &nodesPerHeight)
	}
	for _, node := range rootNodes {
		if node.Height == maxRootHeight {
			adjustChildrenHeight(node, &nodesPerHeight)
		}
	}
	for _, node := range treesNodes {
		if node.Parent != nil && node.Parent.Height >= node.Height {
			adjustChildrenHeight(node.Parent, &nodesPerHeight)
		}
	}
	for _, root := range rootNodes {
		for _, child := range root.Children {
			if root.Height == child.Height {
				root.Height = child.Height + 1
			}
		}
	}
	nodesPerHeight = make(map[int][]*types.TreeNode, 0)
	for _, node := range treesNodes {
		nodesPerHeight[node.Height] = append(nodesPerHeight[node.Height], node)
	}

	maxInputsPerHeight := make(map[int]int, 0)
	for height, nodes := range nodesPerHeight {
		maxInputs := 0
		for _, node := range nodes {
			if version.Data.Nodes[node.NodeName] != nil && len(version.Data.Nodes[node.NodeName].Inputs) > maxInputs {
				maxInputs = len(version.Data.Nodes[node.NodeName].Inputs)
			}
		}
		maxInputsPerHeight[height] = maxInputs
	}

	distance := 400
	X := float64(0)
	for height := 0; height < len(nodesPerHeight); height++ {
		nodes := nodesPerHeight[height]
		sort.SliceStable(nodes, func(i, j int) bool {
			return nodes[i].NodeName < nodes[j].NodeName
		})
		total := (len(nodes) - 1) * distance
		start := -total / 2
		nodeSizeIndent := float64(distance * (maxInputsPerHeight[height] / 15))
		previousHeightNodeSizeIndent := float64(0)
		if height-1 >= 0 {
			previousHeightNodeSizeIndent = float64(distance * (maxInputsPerHeight[height-1] / 10))
		}
		for i, node := range nodes {
			if version.Data.Nodes[node.NodeName] != nil {
				version.Data.Nodes[node.NodeName].Meta.Coordinates.X = X
				if i == 0 && height > 0 {
					version.Data.Nodes[node.NodeName].Meta.Coordinates.X += nodeSizeIndent
				}
				version.Data.Nodes[node.NodeName].Meta.Coordinates.X += previousHeightNodeSizeIndent
				version.Data.Nodes[node.NodeName].Meta.Coordinates.Y = float64(start)
				start += distance
				if i+1 < len(nodes) && version.Data.Nodes[nodes[i+1].NodeName] != nil &&
					len(version.Data.Nodes[nodes[i+1].NodeName].Inputs) == maxInputsPerHeight[height] {
					start += int(nodeSizeIndent)
				}
				if len(version.Data.Nodes[node.NodeName].Inputs) == maxInputsPerHeight[height] {
					start += int(nodeSizeIndent)
				}
			} else if version.Data.PrimitiveNodes[node.NodeName] != nil {
				version.Data.PrimitiveNodes[node.NodeName].Coordinates.X = X
				if i == 0 && height > 0 {
					version.Data.PrimitiveNodes[node.NodeName].Coordinates.X += nodeSizeIndent
				}
				version.Data.PrimitiveNodes[node.NodeName].Coordinates.X += previousHeightNodeSizeIndent
				version.Data.PrimitiveNodes[node.NodeName].Coordinates.Y = float64(start)
				start += distance
				if i+1 < len(nodes) && version.Data.Nodes[nodes[i+1].NodeName] != nil &&
					len(version.Data.Nodes[nodes[i+1].NodeName].Inputs) == maxInputsPerHeight[height] {
					start += int(nodeSizeIndent)
				}
			}
			if i == len(nodes)-1 {
				X += float64(distance * 2)
				X += previousHeightNodeSizeIndent
			}
		}
	}
}
