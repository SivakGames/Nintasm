package enumTokenTypes

type Def int

const (
	//No token type/skippable
	None Def = iota + 0
	//TODO: Replace???
	INIT_LINE_OTHER_CHARS
	//Whitespace
	WHITESPACE
	//Template strings
	DYNAMIC_LABEL
	DYNAMIC_STRING
	//Strings
	STRING
	BACKTICK_STRING
	//Instructions
	INSTRUCTION
	//++++Helper for directive capture
	DIRECTIVE_RANGE_START
	//Directives
	DIRECTIVE_blockStart
	DIRECTIVE_blockEnd
	DIRECTIVE_defCharMap
	DIRECTIVE_dataBytes
	DIRECTIVE_dataSeries
	DIRECTIVE_defExprMap
	DIRECTIVE_deleteSymbol
	DIRECTIVE_exitMacro
	DIRECTIVE_if
	DIRECTIVE_include
	DIRECTIVE_INES
	DIRECTIVE_invokeKeyVal
	DIRECTIVE_mixedData
	DIRECTIVE_reassign
	DIRECTIVE_romBuilding
	DIRECTIVE_setting
	DIRECTIVE_settingReset
	DIRECTIVE_throw
	//++++Helper for directive capture
	DIRECTIVE_RANGE_END
	//Labeled directives - put here as they're treated slightly differently
	DIRECTIVE_labeled
	DIRECTIVE_labeledBlockStart
	DIRECTIVE_labeledBlockEnd
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
	OPERATOR_questionMark
)
