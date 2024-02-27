package parser

import (
	"errors"
	"fmt"
	enumTokenTypes "misc/nintasm/constants/enums/tokenTypes"
	"misc/nintasm/tokenizer"
)

type tokenEnum = enumTokenTypes.Def

const OPERAND_TARGET_TOKENIZER = "operand"

type Parser struct {
	hasMore        bool
	lookaheadType  tokenEnum
	lookaheadValue string
	tokenizer      tokenizer.Tokenizer
}

func New() Parser {
	return Parser{
		hasMore:        false,
		lookaheadType:  enumTokenTypes.None,
		lookaheadValue: "",
		tokenizer:      tokenizer.New(),
	}
}

// Should be called upon first starting parsing
// Will setup the tokenizer and get the tokenizer in position for the next token
func (p *Parser) startAndAdvanceToNext(line string, tokenizerSpecType string) error {
	p.tokenizer.Start(line, tokenizerSpecType)
	err := p.advanceToNext()
	return err
}

// Go to the next token
func (p *Parser) advanceToNext() error {
	hasMore, lookaheadType, lookaheadValue, err := p.tokenizer.GetNextToken()
	if err != nil {
		return err
	}

	p.hasMore = hasMore
	p.lookaheadType = lookaheadType
	p.lookaheadValue = lookaheadValue
	return nil
}

// See if the next token type is the desired token to follow
func (p *Parser) eat(desiredTokenType tokenEnum) error {
	if p.lookaheadType == enumTokenTypes.None {
		return errors.New("UNEXPECTED END OF INPUT")
	}
	if p.lookaheadType != desiredTokenType {
		errMsg := fmt.Sprintf("UNEXPECTED TOKEN: \x1b[33m%v\x1b[0m", p.lookaheadValue)
		return errors.New(errMsg)
	}
	return nil
}

func (p *Parser) eatAndAdvance(desiredTokenType tokenEnum) error {
	err := p.eat(desiredTokenType)
	if err != nil {
		return err
	}
	err = p.advanceToNext()
	if err != nil {
		return err
	}
	return nil
}

// The token is guaranteed
func (p *Parser) eatFreelyAndAdvance(desiredTokenType tokenEnum) error {
	_ = p.eat(desiredTokenType)
	err := p.advanceToNext()
	if err != nil {
		return err
	}
	return nil
}

// Should be called before parsing any operands.
// Will set the tokenizer spec to operand, put the cursor where operands are supposed to start,
// and advance the lookahead where it needs to be.
func (p *OperandParser) SetupOperandParser(line string, operandListStringStartPosition int) error {
	p.operandListStringStartPosition = operandListStringStartPosition
	p.tokenizer.Start(line, OPERAND_TARGET_TOKENIZER)
	p.tokenizer.RepositionCursor(operandListStringStartPosition)
	err := p.advanceToNext()
	if err != nil {
		return err
	}
	return nil
}
