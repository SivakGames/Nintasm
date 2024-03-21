package interpreter

import (
	"math"
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	enumNodeTypes "misc/nintasm/constants/enums/nodeTypes"
	"misc/nintasm/interpreter/environment"
	"misc/nintasm/interpreter/environment/charmapTable"
	"misc/nintasm/interpreter/environment/namespaceTable"
	"misc/nintasm/interpreter/environment/symbolAsNodeTable"
	"misc/nintasm/interpreter/operandFactory"
)

type assemblerFunction struct {
	minArgs          int
	maxArgs          int
	selfEval         bool
	argMustResolveTo []enumNodeTypes.Def
}

var assemblerBuiltInFunctions = map[string]assemblerFunction{
	"high":     {1, 1, true, []enumNodeTypes.Def{enumNodeTypes.NumericLiteral}},
	"low":      {1, 1, true, []enumNodeTypes.Def{enumNodeTypes.NumericLiteral}},
	"ceil":     {1, 1, true, []enumNodeTypes.Def{enumNodeTypes.NumericLiteral}},
	"floor":    {1, 1, true, []enumNodeTypes.Def{enumNodeTypes.NumericLiteral}},
	"round":    {1, 1, true, []enumNodeTypes.Def{enumNodeTypes.NumericLiteral}},
	"modfDeci": {1, 1, true, []enumNodeTypes.Def{enumNodeTypes.NumericLiteral}},
	"modfInt":  {1, 1, true, []enumNodeTypes.Def{enumNodeTypes.NumericLiteral}},
	"sin":      {1, 1, true, []enumNodeTypes.Def{enumNodeTypes.NumericLiteral}},
	"sindeg":   {1, 1, true, []enumNodeTypes.Def{enumNodeTypes.NumericLiteral}},
	"cos":      {1, 1, true, []enumNodeTypes.Def{enumNodeTypes.NumericLiteral}},
	"cosdeg":   {1, 1, true, []enumNodeTypes.Def{enumNodeTypes.NumericLiteral}},

	"strlen": {1, 1, true, []enumNodeTypes.Def{enumNodeTypes.MultiByte}},
	"substr": {2, 3, true, []enumNodeTypes.Def{enumNodeTypes.MultiByte, enumNodeTypes.NumericLiteral, enumNodeTypes.NumericLiteral}},

	"toCharmap":            {1, 1, true, []enumNodeTypes.Def{enumNodeTypes.StringLiteral}},
	"reverseStr":           {1, 1, false, []enumNodeTypes.Def{enumNodeTypes.StringLiteral}},
	"bank":                 {1, 1, false, []enumNodeTypes.Def{enumNodeTypes.Identifier}},
	"defined":              {1, 1, false, []enumNodeTypes.Def{enumNodeTypes.Identifier}},
	"namespaceValuesToStr": {1, 1, false, []enumNodeTypes.Def{enumNodeTypes.Identifier}},
}

func isAssemblerFunction(node *Node) bool {
	funcName := node.NodeValue
	_, isAsmFunc := assemblerBuiltInFunctions[funcName]
	return isAsmFunc
}

// Do built-in assembler function
func processAssemblerFunction(node *Node) error {
	funcName := node.NodeValue
	functionData, _ := assemblerBuiltInFunctions[funcName]

	// ------------------------------------------------------------
	//Check number of arguments usable
	numArgs := len(*node.ArgumentList)
	if numArgs < functionData.minArgs {
		return errorHandler.AddNew(enumErrorCodes.InterpreterFuncTooFewArgs, funcName)
	}
	if numArgs > functionData.maxArgs {
		return errorHandler.AddNew(enumErrorCodes.InterpreterFuncTooManyArgs, funcName)
	}

	// ------------------------------------------------------------
	//Do standard evaluation of node(s)
	evaluatedArguments := make([]Node, len(*node.ArgumentList))

	if functionData.selfEval {
		for i, a := range *node.ArgumentList {
			evaluatedFuncNode, err := EvaluateNode(a)
			if err != nil {
				return err
			}
			if evaluatedFuncNode.NodeType != functionData.argMustResolveTo[i] {
				return errorHandler.AddNew(enumErrorCodes.InterpreterFuncArgWrongType, funcName)
			}
			evaluatedArguments[i] = evaluatedFuncNode
		}
	}

	// ------------------------------------------------------------
	//Actually process the function...
	switch funcName {

	//+-*/+-*/+-*/+-*/+-*/+-*/+-*/+-*/
	//Math functions
	case "high":
		node.AsNumber = (evaluatedArguments[0].AsNumber & 0x0ff00) >> 8
	case "low":
		node.AsNumber = (evaluatedArguments[0].AsNumber & 0x000ff)
	case "ceil":
		node.AsNumber = int(math.Ceil(float64(evaluatedArguments[0].AsNumber)))
	case "floor":
		node.AsNumber = int(math.Floor(float64(evaluatedArguments[0].AsNumber)))
	case "round":
		node.AsNumber = int(math.Round(float64(evaluatedArguments[0].AsNumber)))
	case "modfDeci":
		result, _ := math.Modf(float64(evaluatedArguments[0].AsNumber))
		node.AsNumber = int(result)
	case "modfInt":
		_, result := math.Modf(float64(evaluatedArguments[0].AsNumber))
		node.AsNumber = int(result)
	case "sin":
		node.AsNumber = int(math.Sin(float64(evaluatedArguments[0].AsNumber)))
	case "cos":
		node.AsNumber = int(math.Cos(float64(evaluatedArguments[0].AsNumber)))
	case "sindeg":
		node.AsNumber = int(math.Sin(float64(evaluatedArguments[0].AsNumber) * (180 / math.Pi)))
	case "cosdeg":
		node.AsNumber = int(math.Cos(float64(evaluatedArguments[0].AsNumber) * (180 / math.Pi)))

	case "strlen":
		if node.NodeType == enumNodeTypes.StringLiteral {
			node.AsNumber = len(node.NodeValue)
		} else {
			node.AsNumber = len(*evaluatedArguments[0].ArgumentList)
		}

	case "substr":
		var limit Node
		var slicedNodes []Node

		target := evaluatedArguments[0]
		offset := evaluatedArguments[1]
		if len(evaluatedArguments) > 2 {
			limit = (evaluatedArguments)[2]
			slicedNodes = (*target.ArgumentList)[offset.AsNumber:limit.AsNumber]
		} else {
			slicedNodes = (*target.ArgumentList)[offset.AsNumber:]
		}
		operandFactory.ConvertNodeToMultiBytes(node, slicedNodes)

	case "toCharmap":
		nodeString := (evaluatedArguments[0].NodeValue)
		replacedStringAsBytes, err := charmapTable.MapStringToCharmap(nodeString)
		if err != nil {
			return err
		}
		multiBytes := []Node{}
		for _, r := range replacedStringAsBytes {
			n := operandFactory.CreateNumericLiteralNode(r)
			multiBytes = append(multiBytes, n)
		}
		operandFactory.ConvertNodeToMultiBytes(node, multiBytes)

	// NON pre-evaluated functions!
	case "namespaceValuesToStr":
		namespaceLabel := (*node.ArgumentList)[0].NodeValue
		nsValues, err := namespaceTable.GetNamespaceValues(namespaceLabel)
		if err != nil {
			return err
		}
		nsValuesAsNode := []Node{}
		for _, nsv := range *nsValues {
			fullName := namespaceLabel + nsv.Key
			if !nsv.Resolved {
				return errorHandler.AddNew(enumErrorCodes.NamespaceToValuesNotResolved, nsv.Key)
			}
			lookedUpNode, _, err := environment.LookupIdentifierInSymbolAsNodeTable(fullName)
			if err != nil {
				return err
			}
			switch lookedUpNode.NodeType {
			case enumNodeTypes.NumericLiteral, enumNodeTypes.StringLiteral:
				nsValuesAsNode = append(nsValuesAsNode, lookedUpNode)
			case enumNodeTypes.MultiByte:
				for _, arg := range *lookedUpNode.ArgumentList {
					nsValuesAsNode = append(nsValuesAsNode, arg)
				}
			default:
				panic("Can't unpack for namespace values to str")
			}
		}

		node.NodeValue = namespaceLabel
		operandFactory.ConvertNodeToMultiBytes(node, nsValuesAsNode)

	case "bank":
		_, err := EvaluateNode((*node.ArgumentList)[0])
		if err != nil {
			return err
		}
		bankValue, _ := symbolAsNodeTable.GetValueFromLabelAsBankTable((*node.ArgumentList)[0].NodeValue)
		node.AsNumber = bankValue

	case "defined":
		baseNode := (*node.ArgumentList)[0]
		if baseNode.Resolved {
			node.AsBool = true
			operandFactory.ConvertNodeToBooleanLiteral(node)
		} else if operandFactory.ValidateNodeIsIdentifier(&baseNode) ||
			operandFactory.ValidateNodeIsSubstitutionID(&baseNode) {
			node.AsBool = false
			operandFactory.ConvertNodeToBooleanLiteral(node)
		} else if baseNode.NodeType == enumNodeTypes.Undefined {
			node.AsBool = false
			operandFactory.ConvertNodeToBooleanLiteral(node)
		}
	}

	//Final resolve
	switch funcName {
	case
		"high", "low",
		"ceil", "floor", "round",
		"modfDeci", "modfInt",
		"sin", "sindeg", "cos", "cosdeg",
		"bank",
		"strlen":
		operandFactory.ConvertNodeToNumericLiteral(node)

	}

	return nil
}