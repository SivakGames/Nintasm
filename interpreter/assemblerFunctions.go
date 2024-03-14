package interpreter

import (
	"math"
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	enumNodeTypes "misc/nintasm/constants/enums/nodeTypes"
	"misc/nintasm/interpreter/environment/charmapTable"
	"misc/nintasm/interpreter/environment/symbolAsNodeTable"
	"misc/nintasm/interpreter/operandFactory"
)

type assemblerFunction struct {
	minArgs          int
	maxArgs          int
	argMustResolveTo []enumNodeTypes.Def
}

var assemblerBuiltInFunctions = map[string]assemblerFunction{
	"high":                 {1, 1, []enumNodeTypes.Def{enumNodeTypes.NumericLiteral}},
	"low":                  {1, 1, []enumNodeTypes.Def{enumNodeTypes.NumericLiteral}},
	"ceil":                 {1, 1, []enumNodeTypes.Def{enumNodeTypes.NumericLiteral}},
	"floor":                {1, 1, []enumNodeTypes.Def{enumNodeTypes.NumericLiteral}},
	"round":                {1, 1, []enumNodeTypes.Def{enumNodeTypes.NumericLiteral}},
	"modfDeci":             {1, 1, []enumNodeTypes.Def{enumNodeTypes.NumericLiteral}},
	"modfInt":              {1, 1, []enumNodeTypes.Def{enumNodeTypes.NumericLiteral}},
	"sin":                  {1, 1, []enumNodeTypes.Def{enumNodeTypes.NumericLiteral}},
	"sindeg":               {1, 1, []enumNodeTypes.Def{enumNodeTypes.NumericLiteral}},
	"cos":                  {1, 1, []enumNodeTypes.Def{enumNodeTypes.NumericLiteral}},
	"cosdeg":               {1, 1, []enumNodeTypes.Def{enumNodeTypes.NumericLiteral}},
	"strlen":               {1, 1, []enumNodeTypes.Def{enumNodeTypes.StringLiteral}},
	"substr":               {2, 3, []enumNodeTypes.Def{enumNodeTypes.StringLiteral, enumNodeTypes.NumericLiteral, enumNodeTypes.NumericLiteral}},
	"toCharmap":            {1, 1, []enumNodeTypes.Def{enumNodeTypes.StringLiteral}},
	"reverseStr":           {1, 1, []enumNodeTypes.Def{enumNodeTypes.StringLiteral}},
	"defined":              {1, 1, []enumNodeTypes.Def{enumNodeTypes.Identifier}},
	"bank":                 {1, 1, []enumNodeTypes.Def{enumNodeTypes.Identifier}},
	"namespaceValuesToStr": {1, 1, []enumNodeTypes.Def{enumNodeTypes.Identifier}},
}

// Do built-in function
func processAssemblerFunction(node *Node) (bool, error) {
	funcName := node.NodeValue

	functionData, isAsmFunc := assemblerBuiltInFunctions[funcName]
	if !isAsmFunc {
		return isAsmFunc, nil
	}

	// ------------------------------------------------------------
	//Check number of arguments usable
	numArgs := len(*node.ArgumentList)
	if numArgs < functionData.minArgs {
		return isAsmFunc, errorHandler.AddNew(enumErrorCodes.InterpreterFuncTooFewArgs, funcName)
	}
	if numArgs > functionData.maxArgs {
		return isAsmFunc, errorHandler.AddNew(enumErrorCodes.InterpreterFuncTooManyArgs, funcName)
	}

	// ------------------------------------------------------------
	//Do standard evaluation of node(s)
	switch funcName {
	case "ceil", "floor", "round",
		"high", "low",
		"sin", "sindeg", "cos", "cosdeg",
		"modfDeci", "modfInt",
		"strlen", "substr",
		"toCharmap":
		for i, a := range *node.ArgumentList {
			evaluatedFuncNode, err := EvaluateNode(a)
			if err != nil {
				return isAsmFunc, err
			}
			if evaluatedFuncNode.NodeType != functionData.argMustResolveTo[i] {
				return isAsmFunc, errorHandler.AddNew(enumErrorCodes.InterpreterFuncArgWrongType, funcName)
			}
		}
	}

	// ------------------------------------------------------------
	//Actually process the function...
	switch funcName {
	case "high":
		node.AsNumber = ((*node.ArgumentList)[0].AsNumber & 0x0ff00) >> 8
	case "low":
		node.AsNumber = ((*node.ArgumentList)[0].AsNumber & 0x000ff)
	case "ceil":
		node.AsNumber = int(math.Ceil(float64((*node.ArgumentList)[0].AsNumber)))
	case "floor":
		node.AsNumber = int(math.Floor(float64((*node.ArgumentList)[0].AsNumber)))
	case "round":
		node.AsNumber = int(math.Round(float64((*node.ArgumentList)[0].AsNumber)))
	case "modfDeci":
		result, _ := math.Modf(float64((*node.ArgumentList)[0].AsNumber))
		node.AsNumber = int(result)
	case "modfInt":
		_, result := math.Modf(float64((*node.ArgumentList)[0].AsNumber))
		node.AsNumber = int(result)
	case "sin":
		node.AsNumber = int(math.Sin(float64((*node.ArgumentList)[0].AsNumber)))
	case "cos":
		node.AsNumber = int(math.Cos(float64((*node.ArgumentList)[0].AsNumber)))
	case "sindeg":
		node.AsNumber = int(math.Sin(float64((*node.ArgumentList)[0].AsNumber) * (180 / math.Pi)))
	case "cosdeg":
		node.AsNumber = int(math.Cos(float64((*node.ArgumentList)[0].AsNumber) * (180 / math.Pi)))
	case "strlen":
		node.AsNumber = len((*node.ArgumentList)[0].NodeValue)
	case "namespaceValuesToStr":
		(*node.ArgumentList)[0].NodeValue = "tatata"

	case "toCharmap":
		nodeString := ((*node.ArgumentList)[0].NodeValue)
		replacedStringAsBytes, err := charmapTable.MapStringToCharmap(nodeString)
		if err != nil {
			return isAsmFunc, err
		}
		multiBytes := []Node{}
		for _, r := range replacedStringAsBytes {
			n := operandFactory.CreateNumericLiteralNode(r)
			multiBytes = append(multiBytes, n)
		}
		operandFactory.ConvertNodeToMultiBytes(node, multiBytes)

	case "bank":
		_, err := EvaluateNode((*node.ArgumentList)[0])
		if err != nil {
			return isAsmFunc, err
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

	switch funcName {
	case
		"high", "low",
		"ceil", "floor", "round",
		"modfDeci", "modfInt",
		"sin", "sindeg", "cos", "cosdeg",
		"bank",
		"strlen":
		operandFactory.ConvertNodeToNumericLiteral(node)
	case "namespaceValuesToStr":
		operandFactory.ConvertNodeToStringLiteral(node)
	}

	return isAsmFunc, nil
}
