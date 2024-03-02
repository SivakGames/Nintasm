package tokenizer

import (
	"fmt"
	enumTokenTypes "misc/nintasm/constants/enums/tokenTypes"
	"regexp"
)

type specRegexEnum struct {
	Regex         *regexp.Regexp
	OperationType tokenEnum
}

// All the general specs
var (
	initialLineAnythingSpec = []specRegexEnum{
		{regexp.MustCompile(`^[^\;\s\"\']+`), enumTokenTypes.INIT_LINE_OTHER_CHARS},
	}
	indirectCapturingDelimiterspec = []specRegexEnum{
		{regexp.MustCompile(`^[^\;\s\"\'()\[\]{}]+`), enumTokenTypes.INIT_LINE_OTHER_CHARS},
	}
	whitespaceAcknowledgeSpec = []specRegexEnum{
		{regexp.MustCompile(`^\s+`), enumTokenTypes.WHITESPACE},
	}
	whitespaceIgnoreSpec = []specRegexEnum{
		{regexp.MustCompile(`^\s+`), enumTokenTypes.None},
	}
	commentSpec = []specRegexEnum{
		{regexp.MustCompile(`^;.*`), enumTokenTypes.None},
	}
	stringSpec = []specRegexEnum{
		{regexp.MustCompile(`^\"[^\"]*\"`), enumTokenTypes.STRING},
		{regexp.MustCompile(`^\'[^\']*\'`), enumTokenTypes.STRING},
		{regexp.MustCompile("^\\`[^\\`]*\\`"), enumTokenTypes.BACKTICK_STRING},
	}
	instructionSpec = []specRegexEnum{
		{regexp.MustCompile(generateBoundaries("brk|clc|cld|cli|clv|dex|dey|inx|iny|nop|pha|php|pla|plp|rti|rts|sec|sed|sei|tax|tay|tsx|txa|txs|tya|adc|and|cmp|eor|lda|ora|sbc|asl|lsr|rol|ror|cpx|cpy|dec|inc|ldx|ldy|sta|stx|sty|bit|jmp|jsr|bpl|bmi|bvc|bvs|bcc|bcs|bne|beq")), enumTokenTypes.INSTRUCTION},
	}
	directiveSpec = []specRegexEnum{
		{regexp.MustCompile(generateBoundaries("db|byte|dw|word|dwBe|wordBe|rdb|reverseByte")), enumTokenTypes.DIRECTIVE_dataBytes},
		{regexp.MustCompile(generateBoundaries("ds|pad")), enumTokenTypes.DIRECTIVE_dataSeries},
		{regexp.MustCompile(generateBoundaries("d_[bwe]+_?")), enumTokenTypes.DIRECTIVE_mixedData},
		{regexp.MustCompile(generateBoundaries("inesPrg|inesChr|inesMap|inesMir|inesBat")), enumTokenTypes.DIRECTIVE_INES},

		{regexp.MustCompile(generateBoundaries("romSegment|bank|org")), enumTokenTypes.DIRECTIVE_romBuilding},
		{regexp.MustCompile(generateBoundaries("include|incbin")), enumTokenTypes.DIRECTIVE_include},
		{regexp.MustCompile(generateBoundaries("autoZP|autoZeroPage|emptyRomFill|rsset|setCharMap|setExprMap")), enumTokenTypes.DIRECTIVE_setting},
		{regexp.MustCompile(generateBoundaries("resetCharMap|resetExprMap")), enumTokenTypes.DIRECTIVE_settingReset},
		{regexp.MustCompile(generateBoundaries("throw")), enumTokenTypes.DIRECTIVE_throw},

		{regexp.MustCompile(generateBoundaries("if|ikv|repeat|switch|elseif|ifdef|ifndef|else|case|default")), enumTokenTypes.DIRECTIVE_blockStart},
		{regexp.MustCompile(generateBoundaries("defChar|defCharRange")), enumTokenTypes.DIRECTIVE_defCharMap},
		{regexp.MustCompile(generateBoundaries("defExpr")), enumTokenTypes.DIRECTIVE_defExprMap},
		{regexp.MustCompile(generateBoundaries("kv|keyvalue")), enumTokenTypes.DIRECTIVE_invokeKeyVal},
		{regexp.MustCompile(generateBoundaries("end(?:if|ikv|repeat|switch)")), enumTokenTypes.DIRECTIVE_blockEnd},

		{regexp.MustCompile(generateBoundaries("func|rs")), enumTokenTypes.DIRECTIVE_labeled},
		{regexp.MustCompile(generateBoundaries("nameSpace|macro|kvMacro|charMap|exprMap")), enumTokenTypes.DIRECTIVE_labeledBlockStart},
		{regexp.MustCompile(generateBoundaries("end(?:nameSpace|m|kvm|charMap|exprMap)")), enumTokenTypes.DIRECTIVE_labeledBlockEnd},
	}
	delimiterSpec = []specRegexEnum{
		{regexp.MustCompile(`^,`), enumTokenTypes.DELIMITER_comma},
		{regexp.MustCompile(`^\.`), enumTokenTypes.DELIMITER_period},
		{regexp.MustCompile(`^\(`), enumTokenTypes.DELIMITER_leftParenthesis},
		{regexp.MustCompile(`^\)`), enumTokenTypes.DELIMITER_rightParenthesis},
		{regexp.MustCompile(`^\[`), enumTokenTypes.DELIMITER_leftSquareBracket},
		{regexp.MustCompile(`^\]`), enumTokenTypes.DELIMITER_rightSquareBracket},
		{regexp.MustCompile(`^#`), enumTokenTypes.DELIMITER_hash},
		{regexp.MustCompile(`^:`), enumTokenTypes.DELIMITER_colon},
		{regexp.MustCompile(`^{`), enumTokenTypes.DELIMITER_leftCurlyBrace},
		{regexp.MustCompile(`^}`), enumTokenTypes.DELIMITER_rightCurlyBrace},
	}
	remainingSpec = []specRegexEnum{
		//Register letters (Used for indexes or shifting)
		{regexp.MustCompile(generateBoundaries("a")), enumTokenTypes.REGISTER_A},
		{regexp.MustCompile(generateBoundaries("x")), enumTokenTypes.REGISTER_X},
		{regexp.MustCompile(generateBoundaries("y")), enumTokenTypes.REGISTER_Y},
		//EQU assign
		{regexp.MustCompile(generateBoundaries("equ")), enumTokenTypes.ASSIGN_EQU},
		//Substitutions
		{regexp.MustCompile(`^\\[0-9]+`), enumTokenTypes.SUBSTITUTION_numericID},
		{regexp.MustCompile(`^\\\b[A-Za-z_][0-9A-Za-z_]*\b`), enumTokenTypes.SUBSTITUTION_stringID},
		{regexp.MustCompile(`^\\#`), enumTokenTypes.SUBSTITUTION_numMacroArgs},
		//Numeric literals
		{regexp.MustCompile(`^\$[0-9A-Fa-f]+`), enumTokenTypes.NUMBER_hex},
		{regexp.MustCompile(`^\%[0-1]+`), enumTokenTypes.NUMBER_binary},
		{regexp.MustCompile(`^\d+`), enumTokenTypes.NUMBER_decimal},
		//Identifiers
		{regexp.MustCompile(`^\w+`), enumTokenTypes.IDENTIFIER},
		//Equality operator
		{regexp.MustCompile(`^[=!]=`), enumTokenTypes.OPERATOR_equality},
		//Assignment operator: =
		{regexp.MustCompile(`^=`), enumTokenTypes.ASSIGN_simple},
		//Math operators +, -, *, /
		{regexp.MustCompile(`^[+\-]`), enumTokenTypes.OPERATOR_additive},
		{regexp.MustCompile(`^[*\/\%]`), enumTokenTypes.OPERATOR_multiplicative},

		//Left and right shifts <<, >>
		{regexp.MustCompile(`^(?:<<|>>)`), enumTokenTypes.OPERATOR_shift},

		//Relational operator <, <=, >, >=
		{regexp.MustCompile(`^[><]=?`), enumTokenTypes.OPERATOR_relational},

		//Logical operators &&, ||, !
		{regexp.MustCompile(`^\&\&`), enumTokenTypes.OPERATOR_logicalAnd},
		{regexp.MustCompile(`^\|\|`), enumTokenTypes.OPERATOR_logicalOr},
		{regexp.MustCompile(`^!`), enumTokenTypes.OPERATOR_logicalNot},

		//Boolean operators &, |, ^, ~
		{regexp.MustCompile(`^\&`), enumTokenTypes.OPERATOR_bitwiseAnd},
		{regexp.MustCompile(`^\|`), enumTokenTypes.OPERATOR_bitwiseOr},
		{regexp.MustCompile(`^\^`), enumTokenTypes.OPERATOR_bitwiseXor},
		{regexp.MustCompile(`^\~`), enumTokenTypes.OPERATOR_negate},
	}
)

// All the combined general specs
var (
	CombinedInitialSpec   []specRegexEnum
	CombinedStartLineSpec []specRegexEnum
	CombinedOperandSpec   []specRegexEnum
	IndirectCapturingSpec []specRegexEnum
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

	// IndirectCapturingSpec = append(commentSpec, whitespaceAcknowledgeSpec...)
	// IndirectCapturingSpec = append(IndirectCapturingSpec, delimiterSpec...)
	// IndirectCapturingSpec = append(IndirectCapturingSpec, stringSpec...)
	// IndirectCapturingSpec = append(IndirectCapturingSpec, indirectCapturingDelimiterspec...)

}

// ================================================
func GenerateSpec(specType string) []specRegexEnum {

	switch specType {
	case "initial":
		return CombinedInitialSpec
	case "startLine":
		return CombinedStartLineSpec
	case "operand":
		return CombinedOperandSpec
		//case "indirectCapturing":
		//	return IndirectCapturingSpec
	}
	panic("🛑 UNKNOWN TOKENIZER SPEC NAME")
}

// ++++++++++++++++++++++++++++++++++++++++++++++++
func generateBoundaries(s string) string {
	return fmt.Sprintf("(?i)^\\b(?:%s)\\b", s)
}
