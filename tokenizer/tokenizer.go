package tokenizer

import (
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	enumTokenTypes "misc/nintasm/constants/enums/tokenTypes"
	"regexp"
)

type tokenEnum = enumTokenTypes.Def

// ***********************************************
// Main struct
type Tokenizer struct {
	text          string
	cursor        int
	prevCursor    int
	tokenizerSpec []specRegexEnum
}

// Used for seeing if a value is label-like (used mainly for local labels to avoid token overlap)
var identifierLikeRegex *regexp.Regexp = regexp.MustCompile("^[A-Za-z_][0-9A-Za-z_]*$")
var identifierWithParentLikeRegex *regexp.Regexp = regexp.MustCompile("^[A-Za-z_][0-9A-Za-z_]*\\.[A-Za-z_][0-9A-Za-z_]*$")

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
	t.tokenizerSpec = GenerateSpec(tokenizerSpecType)
}

// ================================================
func (t *Tokenizer) hasMoreTokens() bool {
	return t.cursor < len(t.text)
}

// ================================================
func (t *Tokenizer) GetNextToken() (hasMore bool, tokenType tokenEnum, tokenValue string, err error) {
	if !t.hasMoreTokens() {
		return false, enumTokenTypes.None, "", nil
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
		if specTuple.OperationType == enumTokenTypes.None {
			return t.GetNextToken()
		}
		return true, specTuple.OperationType, tokenValue, nil

	}

	//If reached then it's an unknown/illegal token
	return false, enumTokenTypes.None, "",
		errorHandler.AddNew(enumErrorCodes.TokenizerUnknownIllegalToken, stringSegment)

}

func (t *Tokenizer) GetCursor() int {
	return t.cursor
}

func (t *Tokenizer) GetPrevCursor() int {
	return t.prevCursor
}

// Used for local labels. If a token is LIKE an identifier
// (i.e. starts with letter/underscore and followed by letters, underscores, numbers) then it's OK
func (t *Tokenizer) IsTokenIdentifierLike(identifierLikeToken string) bool {
	return identifierLikeRegex.MatchString(identifierLikeToken)
}
func (t *Tokenizer) IsTokenIdentifierLikeWithParent(identifierLikeToken string) bool {
	return identifierWithParentLikeRegex.MatchString(identifierLikeToken) ||
		identifierLikeRegex.MatchString(identifierLikeToken)
}

func (t *Tokenizer) RepositionCursor(newPosition int) {
	t.cursor = newPosition
	return
}
