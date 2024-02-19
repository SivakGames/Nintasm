package parser_test

import (
	testHelper "misc/nintasm/__tests__/helper"
	"misc/nintasm/assemble"
	"misc/nintasm/romBuilder"
	"testing"
)

func TestCharmaps(t *testing.T) {

	t.Run("***Testing charmap", func(t *testing.T) {
		t.Log("--- Doing Charmap tests ---")
		moduleLines := testHelper.BaseLines
		moduleLines = append(moduleLines, "testCharmap .charmap")
		moduleLines = append(moduleLines, " .defchar \"A\", $10")
		moduleLines = append(moduleLines, " .defchar \"B\", $11")
		moduleLines = append(moduleLines, " .defchar \"C\", $12")
		moduleLines = append(moduleLines, " .defchar \"Ä\", $13")
		moduleLines = append(moduleLines, " .defchar \"Ô\", $14")
		moduleLines = append(moduleLines, " .defchar \"Ý\", $15")
		moduleLines = append(moduleLines, " .defchar \"┏\", $16")
		moduleLines = append(moduleLines, " .defchar \"┓\", $17")
		moduleLines = append(moduleLines, " .defchar \"┗\", $18")
		moduleLines = append(moduleLines, " .defchar \"┛\", $19")
		moduleLines = append(moduleLines, " .defchar \"あ\", $1a")
		moduleLines = append(moduleLines, " .defchar \"い\", $1b")
		moduleLines = append(moduleLines, " .defchar \"ガ\", $1c")
		moduleLines = append(moduleLines, " .defchar \"ギ\", $1d")
		moduleLines = append(moduleLines, " .defchar \"邪\", $1e")
		moduleLines = append(moduleLines, " .defchar \"鬼\", $1f")
		moduleLines = append(moduleLines, " .defchar \"王\", $20")
		moduleLines = append(moduleLines, " .defchar \"💚\", $21")
		moduleLines = append(moduleLines, " .defchar \"🧊\", $22")
		moduleLines = append(moduleLines, "testCharmap .endcharmap")
		moduleLines = append(moduleLines, " .db toCharmap(\"ABCÄÔÝ┏┓┗┛あいガギ邪鬼王💚🧊\")")

		wanted := []uint8{
			uint8(0x10), uint8(0x11), uint8(0x12), uint8(0x13),
			uint8(0x14), uint8(0x15), uint8(0x16), uint8(0x17),
			uint8(0x18), uint8(0x19), uint8(0x1a), uint8(0x1b),
			uint8(0x1c), uint8(0x1d), uint8(0x1e), uint8(0x1f),
			uint8(0x20), uint8(0x21), uint8(0x22),
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