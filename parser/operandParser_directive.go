package parser

import (
	"errors"
	enumTokenTypes "misc/nintasm/enums/tokenTypes"
	"misc/nintasm/romBuilder"
)

type DirectiveOperandParser struct {
	OperandParser
	parentLabel string
}

func NewDirectiveOperandParser() DirectiveOperandParser {
	return DirectiveOperandParser{}
}

func (p *DirectiveOperandParser) Process(operationType tokenEnum, operationValue string) error {

	switch operationType {

	case enumTokenTypes.DIRECTIVE_data:
		if p.lookaheadType == enumTokenTypes.None {
			return errors.New("Directive is empty")
		}
		operandList, err := p.GetOperandList()
		if err != nil {
			return err
		}
		bytesToInsert := make([]uint8, 0)

		for _, operand := range operandList {
			if operand.Resolved {
				asRomData, err := romBuilder.ConvertNodeValueToUInts(operand, 1)
				if err != nil {
					return err
				}
				bytesToInsert = append(bytesToInsert, asRomData...)
			}
		}
		err = romBuilder.AddBytesToRom(bytesToInsert)
		if err != nil {
			return err
		}

		return nil

	default:
		return errors.New("BAD DIRECTIVE OPERATION TYPE!!!")
	}

}

func (p *DirectiveOperandParser) GetParentLabel() string {
	return p.parentLabel
}
