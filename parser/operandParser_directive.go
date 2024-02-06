package parser

import (
	"errors"
	"log"
	enumNodeTypes "misc/nintasm/enums/nodeTypes"
	enumTokenTypes "misc/nintasm/enums/tokenTypes"
	"misc/nintasm/romBuilder"
	"unicode/utf8"
)

type DirectiveOperandParser struct {
	OperandParser
	parentLabel string
}

func NewDirectiveOperandParser() DirectiveOperandParser {
	return DirectiveOperandParser{}
}

func (p *DirectiveOperandParser) Process(operationType tokenEnum, operationValue string) error {
	var evalDirectiveOperandFunc func(Node, int) ([]uint8, error)
	var evalOperandSize int

	operandList, err := p.GetOperandList()
	bytesToInsert := make([]uint8, 0)

	if err != nil {
		return err
	}

	switch operationType {

	case enumTokenTypes.DIRECTIVE_data:
		if len(operandList) == 0 {
			return errors.New("Directive is empty!")
		}
		evalOperandSize = 1
		evalDirectiveOperandFunc = evalDBOperand

	default:
		return errors.New("BAD DIRECTIVE OPERATION TYPE!!!")
	}

	for _, operand := range operandList {
		asRomData, err := evalDirectiveOperandFunc(operand, evalOperandSize)

		if err != nil {
			return err
		}
		bytesToInsert = append(bytesToInsert, asRomData...)
	}

	err = romBuilder.AddBytesToRom(bytesToInsert)
	if err != nil {
		return err
	}

	return nil

}

func evalDBOperand(operand Node, operandSize int) ([]uint8, error) {
	var err error
	asRomData := make([]uint8, 0)

	if operand.Resolved {
		switch operand.NodeType {
		case enumNodeTypes.NumericLiteral, enumNodeTypes.BooleanLiteral:
			asRomData, err = romBuilder.ConvertNodeValueToUInts(operand, operandSize)
			return asRomData, err

		case enumNodeTypes.StringLiteral:
			asRomData := stringToUint8Array(operand.NodeValue)
			return asRomData, nil
		}
		panic("BAD resolved operand for directive eval!")

	} else {
		return []uint8{0}, nil
	}

}

func stringToUint8Array(s string) []uint8 {
	result := make([]uint8, 0, len(s))

	for _, c := range s {
		runeLen := utf8.RuneLen(c)
		if runeLen > 1 {
			log.Println("Rune bigger than 1 byte", c)
			for i := 0; i < runeLen; i++ {
				writeRune := (rune(c) >> (i * 8)) & 0x0ff
				result = append(result, uint8(writeRune))
			}
		} else {
			result = append(result, uint8(rune(c)))
		}

	}

	return result
}
