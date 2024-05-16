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

var builtInMacro_SetPPU = []string{
	" .if \\# == 2",
	" LDA \\1",
	" STA $2006",
	" LDA \\2",
	" STA $2006",
	" .elseif \\# == 1",
	" LDA #high(\\1)",
	" STA $2006",
	" LDA #low(\\1)",
	" STA $2006",
	" .else",
	" .throw \"Too many arguments for predefined __setPPU__ macro!\" ",
	" .endif",
}

type predefMacro struct {
	name  string
	lines *[]string
}

func newPredefMacro(name string, lines *[]string) predefMacro {
	return predefMacro{
		name:  name,
		lines: lines,
	}
}

var macrozzz = []predefMacro{
	newPredefMacro("__resetCode__", &builtInMacro_ResetCode),
	newPredefMacro("__setPPU__", &builtInMacro_SetPPU),
}

var preGenMacros = map[string]macroTable.MacroTableType{}

// -------------------------------------------------

func init() {
	for _, m := range macrozzz {
		generatedMacro := make(macroTable.MacroTableType, len(*m.lines))
		lineOperationParser := parser.NewOperationParser()

		for i, bim := range *m.lines {
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
		preGenMacros[m.name] = generatedMacro
	}
}

// ===================================================

func AddPregensToMacroTable() {
	for _, m := range macrozzz {
		environment.AddOtherIdentifierToMasterTable(m.name, enumSymbolTableTypes.Macro)
		macroTable.AddCapturedLinesToMacro(m.name, preGenMacros[m.name])
	}
}
