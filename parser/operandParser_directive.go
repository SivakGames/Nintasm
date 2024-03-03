package parser

import (
	"fmt"
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	enumTokenTypes "misc/nintasm/constants/enums/tokenTypes"
	"misc/nintasm/parser/directiveHandler"

	"strings"
)

type DirectiveOperandParser struct {
	OperandParser
}

func NewDirectiveOperandParser() DirectiveOperandParser {
	return DirectiveOperandParser{}
}

// Copy of name before checking aliases for error output relating to directive
var unaliasedDirectiveName string

var directiveAliases = map[string]string{
	"BYTE":         "DB",
	"WORD":         "DW",
	"WORDBE":       "DWBE",
	"REVERSEBYTE":  "RDB",
	"PAD":          "DS",
	"AUTOZEROPAGE": "AUTOZP",
	"KEYVALUE":     "KV",
}

var directiveManuallyEvaluatesOperands = map[string]bool{
	"REPEAT":     true,
	"KV":         true,
	"IKV":        true,
	"SETCHARMAP": true,
	"SETEXPRMAP": true,
	"FUNC":       true,
}

var directiveEvaluatesLikeMacroOperands = map[string][]string{
	"KV": {"", "macro"},
}

// Main directive parser
func (p *DirectiveOperandParser) Process(operationTokenEnum tokenEnum, operationValue string, operationLabel string) error {
	var err error
	var captureMasks []string

	directiveName := strings.ToUpper(operationValue)
	unaliasedDirectiveName = directiveName

	aliasValue, aliasExists := directiveAliases[directiveName]
	if aliasExists {
		directiveName = aliasValue
	}

	minOperands, maxOperands := getMinMaxOperandsForDirective(operationTokenEnum, directiveName)
	_, manuallyEvaluatesOperands := directiveManuallyEvaluatesOperands[directiveName]

	//VERY special exception for KV
	evalLikeMacro, ok := directiveEvaluatesLikeMacroOperands[directiveName]
	if ok {
		captureMasks = evalLikeMacro
	}

	operandList, err := p.GetOperandList(minOperands, maxOperands, manuallyEvaluatesOperands, captureMasks)
	if err != nil {
		return err // ❌ Fails
	}

	err = directiveHandler.EvaluateDirective(operationTokenEnum, directiveName, operationLabel, &operandList)
	if err != nil {
		err := errorHandler.CheckErrorContinuesUpwardPropagation(err, enumErrorCodes.Error)
		if err != nil {
			return err // ❌❌ CONTINUES Failing!
		}
	}

	return nil
}

// +++++++++++++++++++++++++

var directiveMinMaxOperands = map[enumTokenTypes.Def][2]int{
	enumTokenTypes.DIRECTIVE_dataBytes:       {1, 128},
	enumTokenTypes.DIRECTIVE_dataSeries:      {1, 2},
	enumTokenTypes.DIRECTIVE_mixedData:       {1, 128},
	enumTokenTypes.DIRECTIVE_blockEnd:        {0, 0},
	enumTokenTypes.DIRECTIVE_labeled:         {1, 1},
	enumTokenTypes.DIRECTIVE_labeledBlockEnd: {0, 0},
	enumTokenTypes.DIRECTIVE_INES:            {1, 1},
	enumTokenTypes.DIRECTIVE_invokeKeyVal:    {2, 2},
	enumTokenTypes.DIRECTIVE_setting:         {1, 1},
	enumTokenTypes.DIRECTIVE_settingReset:    {0, 0},
	enumTokenTypes.DIRECTIVE_throw:           {1, 1},
}

var directiveNameMinMaxOperands = map[string][2]int{
	"BANK":         {1, 1},
	"CHARMAP":      {0, 0},
	"DEFCHAR":      {2, 3},
	"DEFCHARRANGE": {3, 3},
	"DEFEXPR":      {2, 2},
	"ELSE":         {0, 0},
	"ELSEIF":       {1, 1},
	"EXPRMAP":      {0, 0},
	"IF":           {1, 1},
	"IKV":          {1, 1},
	"INCBIN":       {1, 3},
	"INCLUDE":      {1, 1},
	"KVMACRO":      {0, 0},
	"MACRO":        {0, 0},
	"NAMESPACE":    {0, 0},
	"ORG":          {1, 1},
	"REPEAT":       {1, 2},
	"ROMSEGMENT":   {1, 3},
}

func getMinMaxOperandsForDirective(directiveEnum tokenEnum, directiveName string) (int, int) {
	var minMaxOperands [2]int
	var checkOk bool

	minMaxOperands, checkOk = directiveMinMaxOperands[directiveEnum]
	if checkOk {
		return minMaxOperands[0], minMaxOperands[1]
	}
	minMaxOperands, checkOk = directiveNameMinMaxOperands[directiveName]
	if checkOk {
		return minMaxOperands[0], minMaxOperands[1]
	}

	errMsg := fmt.Sprintf("Unable to determine min/max operands for %v directive!", directiveName)
	panic(errMsg)
}
