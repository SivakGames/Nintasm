package predefSymbols

import (
	"misc/nintasm/assemble/blockStack"
	enumSymbolTableTypes "misc/nintasm/constants/enums/symbolTableTypes"
	"misc/nintasm/interpreter/environment"
	"misc/nintasm/interpreter/environment/macroTable"
	"misc/nintasm/parser"
	"misc/nintasm/util"
)

var builtInMacro_ResetCode = []string{
	" SEI",
	" CLD",
	" LDX #$ff",
	" TXS",
	" INX",
	" STX $2000",
	" STX $2001",
	" STX $4010",
	" STX $2003",
	" LDA #$40",
	" STA $4017",
	" LDA #$0f",
	" STA $4015",
	".vbwait1:",
	" BIT $2002",
	" BPL .vbwait1",
	".clearMemory:",
	" LDA #$00",
	" STA $00, X",
	" STA $0100, X",
	" STA $0300, X",
	" STA $0400, X",
	" STA $0500, X",
	" STA $0600, X",
	" STA $0700, X",
	" LDA #$f8",
	" STA $0200, X",
	" INX",
	" BNE .clearMemory",
	".vbwait2:",
	" BIT $2002",
	" BPL .vbwait2",
}

var preGenMacros = map[string]macroTable.MacroTableType{}

func init() {
	generatedMacro := make(macroTable.MacroTableType, len(builtInMacro_ResetCode))
	lineOperationParser := parser.NewOperationParser()

	for i, bim := range builtInMacro_ResetCode {
		lineOperationParser.Process(bim)
		lineOperationParsedValues := lineOperationParser.GetLineOperationValues()
		generatedMacro[i] = blockStack.CapturedLine{
			OriginalLine: bim,
			LineOperationParsedValues: util.LineOperationParsedValues{
				OperationLabel:       lineOperationParsedValues.OperationLabel,
				OperationTokenEnum:   lineOperationParsedValues.OperationTokenEnum,
				OperationTokenValue:  lineOperationParsedValues.OperationTokenValue,
				OperandStartPosition: lineOperationParsedValues.OperandStartPosition,
				ParentParserEnum:     lineOperationParsedValues.ParentParserEnum,
			},
		}
	}
	preGenMacros["__resetCode__"] = generatedMacro
}

func AddPregensToMacroTable() {
	environment.AddOtherIdentifierToMasterTable("__resetCode__", enumSymbolTableTypes.Macro)
	macroTable.AddCapturedLinesToMacro("__resetCode__", macroTable.Macro, preGenMacros["__resetCode__"])
}
