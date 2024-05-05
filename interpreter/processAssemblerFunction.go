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
	"misc/nintasm/romBuilder/nodesToBytes"
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

	//String-related
	"strlen":     {1, 1, true, []enumNodeTypes.Def{enumNodeTypes.StringLiteral}},
	"substr":     {2, 3, true, []enumNodeTypes.Def{enumNodeTypes.StringLiteral, enumNodeTypes.NumericLiteral, enumNodeTypes.NumericLiteral}},
	"reverseStr": {1, 1, true, []enumNodeTypes.Def{enumNodeTypes.StringLiteral}},
	"toCharmap":  {1, 1, true, []enumNodeTypes.Def{enumNodeTypes.StringLiteral}},

	//Array-related
	"itemlen":  {1, 1, true, []enumNodeTypes.Def{enumNodeTypes.MultiByte}},
	"subitem":  {2, 3, true, []enumNodeTypes.Def{enumNodeTypes.MultiByte, enumNodeTypes.NumericLiteral, enumNodeTypes.NumericLiteral}},
	"bytelen":  {1, 1, true, []enumNodeTypes.Def{enumNodeTypes.MultiByte}},
	"contains": {2, 2, true, []enumNodeTypes.Def{enumNodeTypes.MultiByte, enumNodeTypes.NumericLiteral}},
	"index":    {2, 2, true, []enumNodeTypes.Def{enumNodeTypes.MultiByte, enumNodeTypes.NumericLiteral}},

	//Misc
	"bytesInCurrentLabel":  {0, 0, false, []enumNodeTypes.Def{}},
	"bytesInLabel":         {1, 1, false, []enumNodeTypes.Def{enumNodeTypes.Identifier}},
	"bytesInCurrentLocal":  {0, 0, false, []enumNodeTypes.Def{}},
	"bytesInLocal":         {1, 1, false, []enumNodeTypes.Def{enumNodeTypes.Identifier}},
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
		node.AsNumber = float64((int(evaluatedArguments[0].AsNumber) & 0x0ff00) >> 8)
	case "low":
		node.AsNumber = float64(int(evaluatedArguments[0].AsNumber) & 0x000ff)
	case "ceil":
		node.AsNumber = math.Ceil(evaluatedArguments[0].AsNumber)
	case "floor":
		node.AsNumber = math.Floor(evaluatedArguments[0].AsNumber)
	case "round":
		node.AsNumber = math.Round(evaluatedArguments[0].AsNumber)
	case "modfDeci":
		result, _ := math.Modf(evaluatedArguments[0].AsNumber)
		node.AsNumber = result
	case "modfInt":
		_, result := math.Modf(evaluatedArguments[0].AsNumber)
		node.AsNumber = result
	case "sin":
		node.AsNumber = math.Sin(evaluatedArguments[0].AsNumber)
	case "cos":
		node.AsNumber = math.Cos(evaluatedArguments[0].AsNumber)
	case "sindeg":
		node.AsNumber = math.Sin(evaluatedArguments[0].AsNumber * math.Pi / 180)
	case "cosdeg":
		node.AsNumber = math.Cos(evaluatedArguments[0].AsNumber * math.Pi / 180)

	case "strlen":
		node.AsNumber = float64(len(*&evaluatedArguments[0].NodeValue))
	case "bytelen":
		convertedBytes := make([]uint8, 0)
		for _, c := range *evaluatedArguments[0].ArgumentList {
			cConverted, err := nodesToBytes.ConvertNodeValueToUInts(c, 1, false, false, false)
			if err != nil {
				return err
			}
			convertedBytes = append(convertedBytes, cConverted...)
		}
		node.AsNumber = float64(len(convertedBytes))

	case "itemlen":
		node.AsNumber = float64(len(*evaluatedArguments[0].ArgumentList))

	case "reverseStr":
		original := evaluatedArguments[0].NodeValue
		chars := []rune(original)
		for i, j := 0, len(chars)-1; i < j; i, j = i+1, j-1 {
			chars[i], chars[j] = chars[j], chars[i]
		}
		node.NodeValue = string(chars)

	case "substr":
		var slicedString string
		target := evaluatedArguments[0]
		startIndex := int(evaluatedArguments[1].AsNumber)
		if startIndex >= len(target.NodeValue) {
			return errorHandler.AddNew(enumErrorCodes.SubFuncStartTooBig, len(target.NodeValue)-1)
		}

		if len(evaluatedArguments) > 2 {
			endIndex := int(evaluatedArguments[2].AsNumber)
			if endIndex > len(target.NodeValue) {
				return errorHandler.AddNew(enumErrorCodes.SubFuncEndTooBig, len(target.NodeValue))
			} else if endIndex <= startIndex {
				return errorHandler.AddNew(enumErrorCodes.SubFuncEndBiggerThanStart)
			}
			slicedString = (*&target.NodeValue)[startIndex:endIndex]
		} else {
			slicedString = (*&target.NodeValue)[startIndex:]
		}
		node.NodeValue = slicedString

	case "subitem":
		var slicedItem []Node
		target := evaluatedArguments[0]
		startIndex := int(evaluatedArguments[1].AsNumber)
		if startIndex >= len(*target.ArgumentList) {
			return errorHandler.AddNew(enumErrorCodes.SubFuncStartTooBig, len(*target.ArgumentList)-1)
		}

		if len(evaluatedArguments) > 2 {
			endIndex := int(evaluatedArguments[2].AsNumber)
			if endIndex > len(*target.ArgumentList) {
				return errorHandler.AddNew(enumErrorCodes.SubFuncEndTooBig, len(*target.ArgumentList))
			} else if endIndex <= startIndex {
				return errorHandler.AddNew(enumErrorCodes.SubFuncEndBiggerThanStart)
			}
			slicedItem = (*target.ArgumentList)[startIndex:endIndex]
		} else {
			slicedItem = (*target.ArgumentList)[startIndex:]
		}
		node.ArgumentList = &slicedItem

	case "contains":
		target := evaluatedArguments[0].ArgumentList
		containValue := int(evaluatedArguments[1].AsNumber)
		foundContain := false

		for _, v := range *target {
			if containValue == int(v.AsNumber) {
				foundContain = true
				break
			}
		}
		node.AsBool = foundContain

	case "index":
		target := evaluatedArguments[0].ArgumentList
		containValue := int(evaluatedArguments[1].AsNumber)
		foundIndex := -1

		for i, v := range *target {
			if containValue == int(v.AsNumber) {
				foundIndex = i
				break
			}
		}
		if foundIndex < 0 {
			return errorHandler.AddNew(enumErrorCodes.Other, "index function failed and didn't find value")
		}

		node.AsNumber = float64(foundIndex)

	case "toCharmap":
		nodeString := (evaluatedArguments[0].NodeValue)
		replacedStringAsBytes, err := charmapTable.MapStringToCharmap(nodeString)
		if err != nil {
			return err
		}
		multiBytes := []Node{}
		for _, r := range replacedStringAsBytes {
			n := operandFactory.CreateNumericLiteralNode(float64(r))
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
		node.AsNumber = float64(bankValue)

	case "bytesInCurrentLabel",
		"bytesInLabel",
		"bytesInCurrentLocal",
		"bytesInLocal":

		prevLabel := ""
		nextLabel := ""
		exists := false

		switch funcName {
		case "bytesInCurrentLabel":
			parentLabel, err := GetParentLabel()
			if err != nil {
				return err
			}
			prevLabel = parentLabel
		case "bytesInLabel":
			baseNode := (*node.ArgumentList)[0]
			prevLabel = baseNode.NodeValue
		case "bytesInCurrentLocal":
			parentLabel, err := GetParentLabel()
			if err != nil {
				return err
			}
			currentLocalLabel := GetLocalLabel()
			if currentLocalLabel == "" {
				return errorHandler.AddNew(enumErrorCodes.InterpreterNoLocalLabel)
			}
			prevLabel = parentLabel + currentLocalLabel

		case "bytesInLocal":
			return errorHandler.AddNew(enumErrorCodes.Other, "Bytes in local not yet created")
		}

		if funcName == "bytesInCurrentLabel" || funcName == "bytesInLabel" {
			nextLabel, exists = symbolAsNodeTable.GetValueFromPrevLabelNextLabelTable(prevLabel)
		} else {
			nextLabel, exists = symbolAsNodeTable.GetValueFromPrevLocalLabelNextLocalLabelTable(prevLabel)
		}

		// ------------------------------

		prevNode, resolved, err := environment.LookupIdentifierInSymbolAsNodeTable(prevLabel)
		if !resolved {
			return err
		}
		if !exists {
			if environment.GetUnresolvedSilentErrorFlag() {
				return errorHandler.AddUnresolved()
			}
			return errorHandler.AddNew(enumErrorCodes.BytesWithinLabelNoEnd)
		}
		nextNode, resolved, err := environment.LookupIdentifierInSymbolAsNodeTable(nextLabel)
		if !resolved {
			return err
		}
		operationNode := operandFactory.CreateBinaryExpressionNode("-", nextNode, prevNode)
		finalValueNode, err := EvaluateNode(operationNode)
		if err != nil {
			return err
		}
		node.AsNumber = finalValueNode.AsNumber
		operandFactory.ConvertNodeToNumericLiteral(node)
		return nil

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
		"strlen", "bytelen", "itemlen",
		"index":
		operandFactory.ConvertNodeToNumericLiteral(node)
	case "substr", "reverseStr":
		operandFactory.ConvertNodeToStringLiteral(node)
	case "contains":
		operandFactory.ConvertNodeToBooleanLiteral(node)
	case "subitem":
		operandFactory.ConvertNodeToMultiBytes(node, *node.ArgumentList)
	}

	return nil
}
