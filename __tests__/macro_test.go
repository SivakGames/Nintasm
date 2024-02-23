package parser_test

import (
	testHelper "misc/nintasm/__tests__/helper"
	"misc/nintasm/assemble"
	"misc/nintasm/romBuilder"
	"testing"
)

func TestMacros(t *testing.T) {

	t.Run("***Testing macro", func(t *testing.T) {
		t.Log("--- Doing Macro tests ---")
		moduleLines := testHelper.BaseLines
		moduleLines = append(moduleLines, "testMacro .macro", " .db \\1, \\2, \\3", "testMacro .endm")
		moduleLines = append(moduleLines, "testMacro2 .macro", " lda \\1", "testMacro2 .endm")
		moduleLines = append(moduleLines, "testMacro3 .kvmacro", " .db \\aa, \\bb", "testMacro3 .endkvm")
		moduleLines = append(moduleLines, " testMacro 1,2,3")
		moduleLines = append(moduleLines, " testMacro2 #11")
		moduleLines = append(moduleLines, " testMacro2 {22, x}")
		moduleLines = append(moduleLines, " testMacro2 {44, y}")
		moduleLines = append(moduleLines, " testMacro2 {[33], y}")
		moduleLines = append(moduleLines, " .ikv testMacro3")
		moduleLines = append(moduleLines, " .kv \\aa, 55")
		moduleLines = append(moduleLines, " .kv \\bb, 66")
		moduleLines = append(moduleLines, " .endikv")

		wanted := []uint8{
			uint8(1), uint8(2), uint8(3),
			uint8(0xa9), uint8(11),
			uint8(0xb5), uint8(22),
			uint8(0xb9), uint8(44), uint8(0),
			uint8(0xb1), uint8(33),
			uint8(55), uint8(66),
		}

		err := assemble.Start(moduleLines)
		if err != nil {
			t.Error(err)
		}

		iterateOverResultWanted(t, wanted)
		romBuilder.ClearINES()
		romBuilder.ClearRom()
	})

}
