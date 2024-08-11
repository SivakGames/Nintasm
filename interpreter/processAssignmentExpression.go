package interpreter

import (
	enumNodeTypes "misc/nintasm/constants/enums/nodeTypes"
	enumSymbolTableTypes "misc/nintasm/constants/enums/symbolTableTypes"
	"misc/nintasm/interpreter/environment"
	"misc/nintasm/interpreter/environment/namespaceTable"
	"strings"
)

func processAssignmentExpression(node Node) (Node, error) {
	var symbolType enumSymbolTableTypes.Def

	assignmentTypeIsLabel := node.NodeType == enumNodeTypes.AssignLabelExpression
	nodeHasResolved := false

	//Left node is the label itself
	symbolName := (*node.Left).NodeValue
	originalSymbolName := symbolName
	isLocal := strings.HasPrefix(symbolName, ".")

	if isLocal {
		parentLabel, err := GetParentLabel()
		if err != nil {
			return node, err
		}
		symbolName = parentLabel + symbolName
		if namespaceTable.IsDefiningNamespace {
			defer func() {
				namespaceTable.AddKeyToCurrentNamespace(parentLabel, originalSymbolName, nodeHasResolved)
			}()
		}
	}

	// Right node is the expression to set the label to
	evaluatedLabelNode, err := EvaluateNode(*node.Right)
	if err != nil {
		return node, err
	}

	if assignmentTypeIsLabel {
		symbolType = enumSymbolTableTypes.Label
	} else {
		symbolType = enumSymbolTableTypes.SymbolAsNode
	}

	err = environment.AddIdentifierToSymbolAsNodeTable(symbolName, evaluatedLabelNode, symbolType)
	if err != nil {
		return node, err
	}

	//Labels also add an entry to the bank ID table
	if assignmentTypeIsLabel {
		environment.AddToLabelAsBankTable(symbolName)
	}

	nodeHasResolved = true
	node.Left = nil
	node.Right = nil
	return node, nil
}
