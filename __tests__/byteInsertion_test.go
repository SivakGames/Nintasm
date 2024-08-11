package parser_test

import (
	testHelper "misc/nintasm/__tests__/helper"
	"misc/nintasm/assemble"
	"misc/nintasm/romBuilder"
	"testing"
)

func iterateOverResultWanted(t *testing.T, wanted []uint8) {
	bankSegment := romBuilder.GetCurrentBankSegmentBytes()
	result := (*bankSegment)
	for i := range wanted {
		if result[i] != wanted[i] {
			t.Error("= Result at index", i, "not as expected")
			t.Log("Wanted:", wanted[i], "Got:", result[i])
		}
	}
	return
}

func TestBytes(t *testing.T) {

	t.Run("***Testing .db", func(t *testing.T) {
		t.Log("--- Doing DB tests ---")
		var lineCounter uint = 0
		moduleLines := testHelper.BaseLines
		moduleLines = append(moduleLines, " .db 6,5,0,2,\"aた🍁\"")

		wanted := []uint8{
			uint8(6), uint8(5), uint8(0), uint8(2),
			uint8(rune('a')),
			uint8(rune('た') & 0x00000ff), uint8(rune('た') & 0x000ff00 >> 8), uint8(rune('た') & 0x0ff0000 >> 16),
			uint8(rune('🍁') & 0x00000ff), uint8(rune('🍁') & 0x000ff00 >> 8), uint8(rune('🍁') & 0x0ff0000 >> 16),
		}

		err := assemble.ReadLines(&moduleLines, &lineCounter)
		if err != nil {
			t.Error(err)
		}

		iterateOverResultWanted(t, wanted)
		romBuilder.ClearINES()
		romBuilder.ClearRom()
	})

	t.Run("Testing .dw/.dwbe", func(t *testing.T) {
		t.Log("--- Doing DW/DWBE tests ---")
		var lineCounter uint = 0
		moduleLines := testHelper.BaseLines
		moduleLines = append(moduleLines, " .dw $1234, $5678", " .dwbe $1234, $5678")
		wanted := []uint8{0x34, 0x12, 0x78, 0x56, 0x12, 0x34, 0x56, 0x78}

		err := assemble.ReadLines(&moduleLines, &lineCounter)
		if err != nil {
			t.Error(err)
		}

		iterateOverResultWanted(t, wanted)
		romBuilder.ClearINES()
		romBuilder.ClearRom()
	})

	t.Run("Testing .d_***", func(t *testing.T) {
		t.Log("--- Doing D_*** tests ---")
		var lineCounter uint = 0
		moduleLines := testHelper.BaseLines
		moduleLines = append(moduleLines, " .d_bw $01, $4523, $67, $0089", " .d_eb_ $1234, 5, 6, 7")
		wanted := []uint8{0x01, 0x23, 0x45, 0x67, 0x89, 0x00, 0x12, 0x34, 5, 6, 7}

		err := assemble.ReadLines(&moduleLines, &lineCounter)
		if err != nil {
			t.Error(err)
		}

		iterateOverResultWanted(t, wanted)
		romBuilder.ClearINES()
		romBuilder.ClearRom()
	})

	t.Run("Testing .ds", func(t *testing.T) {
		t.Log("--- Doing DS tests ---")
		var lineCounter uint = 0
		moduleLines := testHelper.BaseLines
		moduleLines = append(moduleLines, " .ds 3, 5 ")
		wanted := []uint8{5, 5, 5}

		err := assemble.ReadLines(&moduleLines, &lineCounter)
		if err != nil {
			t.Error(err)
		}

		iterateOverResultWanted(t, wanted)
		romBuilder.ClearINES()
		romBuilder.ClearRom()
	})

}
