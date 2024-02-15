package parser

import (
	"errors"
	enumTokenTypes "misc/nintasm/enums/tokenTypes"
	handlerDirective "misc/nintasm/handlers/directives"

	"strings"
)

type DirectiveOperandParser struct {
	OperandParser
	parentLabel string
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
}

var directiveManuallyEvaluatesOperands = map[string]bool{
	"REPEAT": true,
}

// Main directive parser
func (p *DirectiveOperandParser) Process(operationType tokenEnum, operationValue string, operationLabel string) error {
	var err error

	directiveName := strings.ToUpper(operationValue)
	unaliasedDirectiveName = directiveName

	aliasValue, aliasExists := directiveAliases[directiveName]
	if aliasExists {
		directiveName = aliasValue
	}

	minOperands, maxOperands, err := getMinMaxOperandsForDirective(operationType, directiveName)
	if err != nil {
		return err // ❌ Fails
	}

	_, manuallyEvaluatesOperands := directiveManuallyEvaluatesOperands[directiveName]

	operandList, err := p.GetOperandList(minOperands, maxOperands, manuallyEvaluatesOperands)
	if err != nil {
		return err // ❌ Fails
	}

	err = handlerDirective.Process(operationType, directiveName, &operandList)
	return err

}

// +++++++++++++++++++++++++

var directiveMinMaxOperands = map[enumTokenTypes.Def][2]int{
	enumTokenTypes.DIRECTIVE_dataBytes:       {1, 128},
	enumTokenTypes.DIRECTIVE_dataSeries:      {1, 2},
	enumTokenTypes.DIRECTIVE_mixedData:       {1, 128},
	enumTokenTypes.DIRECTIVE_blockEnd:        {0, 0},
	enumTokenTypes.DIRECTIVE_labeledBlockEnd: {0, 0},
	enumTokenTypes.DIRECTIVE_INES:            {1, 1},
	enumTokenTypes.DIRECTIVE_setting:         {1, 1},
	enumTokenTypes.DIRECTIVE_settingReset:    {0, 0},
	enumTokenTypes.DIRECTIVE_throw:           {1, 1},
}

var directiveNameMinMaxOperands = map[string][2]int{
	"BANK":       {1, 1},
	"INCBIN":     {1, 3},
	"INCLUDE":    {1, 1},
	"ORG":        {1, 1},
	"ROMSEGMENT": {1, 3},
	"REPEAT":     {1, 2},
	"IF":         {1, 1},
	"ELSEIF":     {1, 1},
	"ELSE":       {0, 0},
	"MACRO":      {0, 0},
}

func getMinMaxOperandsForDirective(directiveEnum tokenEnum, directiveName string) (int, int, error) {
	var minMaxOperands [2]int
	var checkOk bool

	minMaxOperands, checkOk = directiveMinMaxOperands[directiveEnum]
	if checkOk {
		return minMaxOperands[0], minMaxOperands[1], nil
	}
	minMaxOperands, checkOk = directiveNameMinMaxOperands[directiveName]
	if !checkOk {
		return 0, 0, errors.New("Unable to determine min/max operands for directive!")
	}
	return minMaxOperands[0], minMaxOperands[1], nil
}
