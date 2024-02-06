package parser_test

import (
	"misc/nintasm/assemble"
	"misc/nintasm/romBuilder"
	"testing"
)

func TestDirectiveOperandParser(t *testing.T) {
	t.Run("Testing .db", func(t *testing.T) {

		//testLine1 := " .db 6,5,\"0\",2, \"あ\""
		testLine1 := " .db \"aた🍁\""
		//wanted := []uint8{6, 5, uint8(rune('0')), 2, uint8(rune('あ') & 0x000ff), uint8(rune('あ') & 0x0ff00 >> 8)}
		wanted := []uint8{uint8(rune('a')),
			uint8(rune('た') & 0x00000ff),
			uint8(rune('た') & 0x000ff00 >> 8),
			uint8(rune('た') & 0x0ff0000 >> 16),
			uint8(rune('🍁') & 0x00000ff),
			uint8(rune('🍁') & 0x000ff00 >> 8),
			uint8(rune('🍁') & 0x0ff0000 >> 16),
		}

		lines := []string{testLine1}
		err := assemble.Start(lines)
		if err != nil {
			t.Error(err)
		}

		bankSegment := romBuilder.GetCurrentBankSegment()
		result := (*bankSegment)
		for i := range wanted {
			if result[i] != wanted[i] {
				t.Error("= Result at index", i, "not as expected")
				t.Log("Wanted:", wanted[i], "Got:", result[i])
			}
		}

	})
	t.Run("Testing .ds", func(t *testing.T) {

		testLine1 := " .ds 2, 5 "
		wanted := []uint8{5, 5}

		lines := []string{testLine1}
		err := assemble.Start(lines)
		if err != nil {
			t.Error(err)
		}

		bankSegment := romBuilder.GetCurrentBankSegment()
		result := (*bankSegment)
		for i := range wanted {
			if result[i] != wanted[i] {
				t.Error("= Result at index", i, "not as expected")
				t.Log("Wanted:", wanted[i], "Got:", result[i])
			}
		}

	})

}
