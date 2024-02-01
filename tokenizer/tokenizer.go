package tokenizer

import (
	"errors"
	"fmt"
	"regexp"

	"misc/nintasm/tokenizer/tokenizerSpec"
)

// ***********************************************
// Main struct
type Tokenizer struct {
	text          string
	cursor        int
	prevCursor    int
	tokenizerSpec []tokenizerSpec.TokenizerSpec
}

// Used for seeing if a value is label-like (used mainly for local labels to avoid token overlap)
var identifierLikeRegex *regexp.Regexp = regexp.MustCompile("^[A-Za-z_][0-9A-Za-z_]*")

func New() Tokenizer {
	return Tokenizer{
		text:       "",
		cursor:     -1,
		prevCursor: -1,
	}
}

// ================================================
func (t *Tokenizer) Start(text string, tokenizerSpecType string) {
	t.text = text
	t.cursor = 0
	t.prevCursor = 0
	t.tokenizerSpec = tokenizerSpec.GenerateSpec(tokenizerSpecType)
}

// ================================================
func (t *Tokenizer) hasMoreTokens() bool {
	return t.cursor < len(t.text)
}

// ================================================
func (t *Tokenizer) GetNextToken() (hasMore bool, tokenType tokenizerSpec.TokenType, tokenValue string, err error) {
	if !t.hasMoreTokens() {
		return false, tokenizerSpec.None, "", nil
	}
	stringSegment := t.text[t.cursor:]

	for _, specTuple := range t.tokenizerSpec {
		tokenValue := specTuple.Regex.FindString(stringSegment)
		if tokenValue == "" {
			continue
		}

		//Advance the cursor
		t.prevCursor = t.cursor
		t.cursor += len(tokenValue)
		if specTuple.OperationType == tokenizerSpec.None {
			return t.GetNextToken()
		}
		return true, specTuple.OperationType, tokenValue, nil

	}

	//If reached then it's an unknown/illegal token
	illegalTokenMessage := fmt.Sprintf("UNKNOWN/ILLEGAL TOKEN: \x1b[31m%v\x1b[0m", stringSegment)

	return false, tokenizerSpec.None, "", errors.New(illegalTokenMessage)
}

func (t *Tokenizer) GoBackToPrev() {
	t.cursor = t.prevCursor
	t.prevCursor = -1
}

func (t *Tokenizer) GetCursor() int {
	return t.cursor
}

// Used for local labels. If a token is LIKE an identifier
// (i.e. starts with letter/underscore and followed by letters, underscores, numbers) then it's OK
func (t *Tokenizer) IsTokenIdentifierLike(identifierLikeToken string) bool {
	return identifierLikeRegex.MatchString(identifierLikeToken)
}

func (t *Tokenizer) RepositionCursor(newPosition int) {
	t.cursor = newPosition
	return
}
