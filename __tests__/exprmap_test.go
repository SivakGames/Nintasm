package parser_test

import (
	testHelper "misc/nintasm/__tests__/helper"
	"misc/nintasm/assemble"
	"misc/nintasm/romBuilder"
	"testing"
)

func TestExprmaps(t *testing.T) {

	t.Run("***Testing exprmap", func(t *testing.T) {
		t.Log("--- Doing Exprmap tests ---")
		moduleLines := testHelper.BaseLines
		moduleLines = append(moduleLines, "testExprmap .exprmap")
		moduleLines = append(moduleLines, " .defexpr \"A\", $10")
		moduleLines = append(moduleLines, " .defexpr \"B1\", $11")
		moduleLines = append(moduleLines, " .defexpr \"2C3\", $12")
		moduleLines = append(moduleLines, " .defexpr \"D#4\", $13")
		moduleLines = append(moduleLines, "testExprmap .endexprmap")
		moduleLines = append(moduleLines, "testExprmap2 .exprmap")
		moduleLines = append(moduleLines, " .defexpr \"A\", $20")
		moduleLines = append(moduleLines, " .defexpr \"B1\", $21")
		moduleLines = append(moduleLines, " .defexpr \"2C3\", $22")
		moduleLines = append(moduleLines, " .defexpr \"D#4\", $23")
		moduleLines = append(moduleLines, "testExprmap2 .endexprmap")

		moduleLines = append(moduleLines, " .db `A`, `B1`, `2C3`, `D#4`")
		moduleLines = append(moduleLines, " .setExprMap testExprmap2")
		moduleLines = append(moduleLines, " .db `A`, `B1`, `2C3`, `D#4`")

		wanted := []uint8{
			uint8(0x10), uint8(0x11), uint8(0x12), uint8(0x13),
			uint8(0x20), uint8(0x21), uint8(0x22), uint8(0x23),
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
