package tokenizerSpec

import (
	"fmt"
	"regexp"
)

const DIRECTIVE_SUFFIX = "_DIRECTIVE"

type TokenType int

const (
	//No token type/skippable
	None TokenType = iota + 0
	//TODO: Replace???
	INIT_LINE_OTHER_CHARS
	//Whitespace
	WHITESPACE
	//Strings
	STRING
	BACKTICK_STRING
	//Instructions
	INSTRUCTION
	//++++Helper for directive capture
	DIRECTIVE_RANGE_START
	//Directives
	DIRECTIVE_charmap
	DIRECTIVE_data
	DIRECTIVE_exprMap
	DIRECTIVE_if
	DIRECTIVE_include
	DIRECTIVE_INES
	DIRECTIVE_invokeKeyVal
	DIRECTIVE_labeled
	DIRECTIVE_labeledBlockStart
	DIRECTIVE_labeledBlockEnd
	DIRECTIVE_mixedData
	DIRECTIVE_repeat
	DIRECTIVE_romBuilding
	DIRECTIVE_setting
	DIRECTIVE_throw
	//++++Helper for directive capture
	DIRECTIVE_RANGE_END
	//Delimiters
	DELIMITER_comma
	DELIMITER_period
	DELIMITER_leftParenthesis
	DELIMITER_rightParenthesis
	DELIMITER_leftSquareBracket
	DELIMITER_rightSquareBracket
	DELIMITER_hash
	DELIMITER_colon
	DELIMITER_leftCurlyBrace
	DELIMITER_rightCurlyBrace
	//Misc - registers
	REGISTER_A
	REGISTER_X
	REGISTER_Y
	//Misc - equality and assign
	ASSIGN_EQU
	ASSIGN_simple
	//Misc - substitutions
	SUBSTITUTION_numericID
	SUBSTITUTION_stringID
	SUBSTITUTION_numMacroArgs
	//Misc - numbers
	NUMBER_hex
	NUMBER_binary
	NUMBER_decimal
	//Misc - Identifiers
	IDENTIFIER
	//Operators
	OPERATOR_equality
	OPERATOR_additive
	OPERATOR_multiplicative
	OPERATOR_shift
	OPERATOR_relational
	OPERATOR_logicalAnd
	OPERATOR_logicalOr
	OPERATOR_logicalNot
	OPERATOR_bitwiseAnd
	OPERATOR_bitwiseOr
	OPERATOR_bitwiseXor
	OPERATOR_negate
)

// +++++++++++++++++++++++++++++++++++++++++++++++
// Helpers
type TokenizerSpec struct {
	Regex         *regexp.Regexp
	OperationType TokenType
}

// All the general specs
var (
	initialLineAnythingSpec = []TokenizerSpec{
		{regexp.MustCompile(`^[^\;\s\"\']+`), INIT_LINE_OTHER_CHARS},
	}
	indirectCapturingDelimiterspec = []TokenizerSpec{
		{regexp.MustCompile(`^[^\;\s\"\'()\[\]{}]+`), INIT_LINE_OTHER_CHARS},
	}
	whitespaceAcknowledgeSpec = []TokenizerSpec{
		{regexp.MustCompile(`^\s+`), WHITESPACE},
	}
	whitespaceIgnoreSpec = []TokenizerSpec{
		{regexp.MustCompile(`^\s+`), None},
	}
	commentSpec = []TokenizerSpec{
		{regexp.MustCompile(`^;.*`), None},
	}
	stringSpec = []TokenizerSpec{
		{regexp.MustCompile(`^\"[^\"]*\"`), STRING},
		{regexp.MustCompile(`^\'[^\']*\'`), STRING},
		{regexp.MustCompile("^\\`[^\\`]*\\`"), BACKTICK_STRING},
	}
	instructionSpec = []TokenizerSpec{
		{regexp.MustCompile(generateBoundaries("brk|clc|cld|cli|clv|dex|dey|inx|iny|nop|pha|php|pla|plp|rti|rts|sec|sed|sei|tax|tay|tsx|txa|txs|tya|adc|and|cmp|eor|lda|ora|sbc|asl|lsr|rol|ror|cpx|cpy|dec|inc|ldx|ldy|sta|stx|sty|bit|jmp|jsr|bpl|bmi|bvc|bvs|bcc|bcs|bne|beq")), INSTRUCTION},
	}
	directiveSpec = []TokenizerSpec{
		{regexp.MustCompile(generateBoundaries("db|byte|dw|word|dwBe|wordBe|rdb|reverseByte|ds|pad|ed|exprData")), DIRECTIVE_data},
		{regexp.MustCompile(generateBoundaries("d_[bwe]+_?")), DIRECTIVE_mixedData},
		{regexp.MustCompile(generateBoundaries("repeat|endrepeat")), DIRECTIVE_repeat},
		{regexp.MustCompile(generateBoundaries("if|elseif|ifdef|ifndef|else|endif|switch|case|default|endswitch")), DIRECTIVE_if},
		{regexp.MustCompile(generateBoundaries("ikv|kv|keyvalue|endIkv")), DIRECTIVE_invokeKeyVal},
		{regexp.MustCompile(generateBoundaries("func|rs")), DIRECTIVE_labeled},
		{regexp.MustCompile(generateBoundaries("nameSpace|macro|kvMacro|charMap|exprMap")), DIRECTIVE_labeledBlockStart},
		{regexp.MustCompile(generateBoundaries("endNameSpace|endm|endKVM|endCharMap|endExprMap")), DIRECTIVE_labeledBlockEnd},
		{regexp.MustCompile(generateBoundaries("inesPrg|inesChr|inesMap|inesMir|inesBat")), DIRECTIVE_INES},
		{regexp.MustCompile(generateBoundaries("romSegment|bank|org")), DIRECTIVE_romBuilding},
		{regexp.MustCompile(generateBoundaries("include|incbin")), DIRECTIVE_include},
		{regexp.MustCompile(generateBoundaries("defChar|defCharRange")), DIRECTIVE_charmap},
		{regexp.MustCompile(generateBoundaries("defExpr")), DIRECTIVE_exprMap},
		{regexp.MustCompile(generateBoundaries("autoZP|autoZeroPage|emptyRomFill|rsset|resetCharMap|setCharMap|resetExprMap|setExprMap")), DIRECTIVE_setting},
		{regexp.MustCompile(generateBoundaries("throw")), DIRECTIVE_throw},
	}
	delimiterSpec = []TokenizerSpec{
		{regexp.MustCompile(`^,`), DELIMITER_comma},
		{regexp.MustCompile(`^\.`), DELIMITER_period},
		{regexp.MustCompile(`^\(`), DELIMITER_leftParenthesis},
		{regexp.MustCompile(`^\)`), DELIMITER_rightParenthesis},
		{regexp.MustCompile(`^\[`), DELIMITER_leftSquareBracket},
		{regexp.MustCompile(`^\]`), DELIMITER_rightSquareBracket},
		{regexp.MustCompile(`^#`), DELIMITER_hash},
		{regexp.MustCompile(`^:`), DELIMITER_colon},
		{regexp.MustCompile(`^{`), DELIMITER_leftCurlyBrace},
		{regexp.MustCompile(`^}`), DELIMITER_rightCurlyBrace},
	}
	remainingSpec = []TokenizerSpec{
		//Register letters (Used for indexes or shifting)
		{regexp.MustCompile(generateBoundaries("a")), REGISTER_A},
		{regexp.MustCompile(generateBoundaries("x")), REGISTER_X},
		{regexp.MustCompile(generateBoundaries("y")), REGISTER_Y},
		//EQU assign
		{regexp.MustCompile(generateBoundaries("equ")), ASSIGN_EQU},
		//Substitutions
		{regexp.MustCompile(`^\\[0-9]+`), SUBSTITUTION_numericID},
		{regexp.MustCompile(`^\\\b[A-Za-z_][0-9A-Za-z_]*\b`), SUBSTITUTION_stringID},
		{regexp.MustCompile(`^\\#`), SUBSTITUTION_numMacroArgs},
		//Numeric literals
		{regexp.MustCompile(`^\$[0-9A-Fa-f]+`), NUMBER_hex},
		{regexp.MustCompile(`^\%[0-1]+`), NUMBER_binary},
		{regexp.MustCompile(`^\d+`), NUMBER_decimal},
		//Identifiers
		{regexp.MustCompile(`^\w+`), IDENTIFIER},
		//Equality operator
		{regexp.MustCompile(`^[=!]=`), OPERATOR_equality},
		//Assignment operator: =
		{regexp.MustCompile(`^=`), ASSIGN_simple},
		//Math operators +, -, *, /
		{regexp.MustCompile(`^[+\-]`), OPERATOR_additive},
		{regexp.MustCompile(`^[*\/\%]`), OPERATOR_multiplicative},

		//Left and right shifts <<, >>
		{regexp.MustCompile(`^(?:<<|>>)`), OPERATOR_shift},

		//Relational operator <, <=, >, >=
		{regexp.MustCompile(`^[><]=?`), OPERATOR_relational},

		//Logical operators &&, ||, !
		{regexp.MustCompile(`^\&\&`), OPERATOR_logicalAnd},
		{regexp.MustCompile(`^\|\|`), OPERATOR_logicalOr},
		{regexp.MustCompile(`^!`), OPERATOR_logicalNot},

		//Boolean operators &, |, ^, ~
		{regexp.MustCompile(`^\&`), OPERATOR_bitwiseAnd},
		{regexp.MustCompile(`^\|`), OPERATOR_bitwiseOr},
		{regexp.MustCompile(`^\^`), OPERATOR_bitwiseXor},
		{regexp.MustCompile(`^\~`), OPERATOR_negate},
	}
)

// All the combined general specs
var (
	CombinedInitialSpec   []TokenizerSpec
	CombinedStartLineSpec []TokenizerSpec
	CombinedOperandSpec   []TokenizerSpec
	IndirectCapturingSpec []TokenizerSpec
)

// Generate the combined specs
func init() {
	CombinedInitialSpec = append(commentSpec, whitespaceAcknowledgeSpec...)
	CombinedInitialSpec = append(CombinedInitialSpec, commentSpec...)
	CombinedInitialSpec = append(CombinedInitialSpec, stringSpec...)
	CombinedInitialSpec = append(CombinedInitialSpec, initialLineAnythingSpec...)

	CombinedStartLineSpec = append(commentSpec, whitespaceAcknowledgeSpec...)
	CombinedStartLineSpec = append(CombinedStartLineSpec, delimiterSpec...)
	CombinedStartLineSpec = append(CombinedStartLineSpec, instructionSpec...)
	CombinedStartLineSpec = append(CombinedStartLineSpec, directiveSpec...)
	CombinedStartLineSpec = append(CombinedStartLineSpec, remainingSpec...)
	CombinedStartLineSpec = append(CombinedStartLineSpec, stringSpec...)

	CombinedOperandSpec = append(commentSpec, whitespaceIgnoreSpec...)
	CombinedOperandSpec = append(CombinedOperandSpec, delimiterSpec...)
	CombinedOperandSpec = append(CombinedOperandSpec, instructionSpec...)
	CombinedOperandSpec = append(CombinedOperandSpec, remainingSpec...)
	CombinedOperandSpec = append(CombinedOperandSpec, stringSpec...)

	IndirectCapturingSpec = append(commentSpec, whitespaceAcknowledgeSpec...)
	IndirectCapturingSpec = append(IndirectCapturingSpec, delimiterSpec...)
	IndirectCapturingSpec = append(IndirectCapturingSpec, stringSpec...)
	IndirectCapturingSpec = append(IndirectCapturingSpec, indirectCapturingDelimiterspec...)

}

// ================================================
func GenerateSpec(specType string) []TokenizerSpec {

	switch specType {
	case "initial":
		return CombinedInitialSpec
	case "startLine":
		return CombinedStartLineSpec
	case "operand":
		return CombinedOperandSpec
	case "indirectCapturing":
		return IndirectCapturingSpec
	default:
		fmt.Println("UNKNOWN TOKENIZER SPEC NAME")
	}
	return nil
}

// ++++++++++++++++++++++++++++++++++++++++++++++++
func generateBoundaries(s string) string {
	return fmt.Sprintf("(?i)^\\b(?:%s)\\b", s)
}
