package parser_test

import (
	"fmt"
	"misc/nintasm/assemble"
	"misc/nintasm/romBuilder"
	"testing"
)

func iterateOverResultWanted(t *testing.T, wanted []uint8) {
	bankSegment := romBuilder.GetCurrentBankSegment()
	result := (*bankSegment)
	for i := range wanted {
		if result[i] != wanted[i] {
			t.Error("= Result at index", i, "not as expected")
			t.Log("Wanted:", wanted[i], "Got:", result[i])
		}
	}
	return
}

func TestDirectiveOperandParser(t *testing.T) {
	t.Run("Testing .db", func(t *testing.T) {
		lines := []string{" .db 6,5,0,2,\"aた🍁\""}
		wanted := []uint8{
			uint8(6), uint8(5), uint8(0), uint8(2),
			uint8(rune('a')),
			uint8(rune('た') & 0x00000ff), uint8(rune('た') & 0x000ff00 >> 8), uint8(rune('た') & 0x0ff0000 >> 16),
			uint8(rune('🍁') & 0x00000ff), uint8(rune('🍁') & 0x000ff00 >> 8), uint8(rune('🍁') & 0x0ff0000 >> 16),
		}

		err := assemble.Start(lines)
		if err != nil {
			t.Error(err)
		}

		iterateOverResultWanted(t, wanted)
	})

	t.Run("Testing .dw/.dwbe", func(t *testing.T) {
		lines := []string{" .dw $1234, $5678", " .dwbe $1234, $5678"}
		wanted := []uint8{0x34, 0x12, 0x78, 0x56, 0x12, 0x34, 0x56, 0x78}

		err := assemble.Start(lines)
		if err != nil {
			t.Error(err)
		}

		iterateOverResultWanted(t, wanted)
	})

	t.Run("Testing .d_***", func(t *testing.T) {
		lines := []string{" .d_bw $01, $4523, $67, $0089", " .d_eb_ $1234, 5, 6, 7"}
		wanted := []uint8{0x01, 0x23, 0x45, 0x67, 0x89, 0x00, 0x12, 0x34, 5, 6, 7}

		err := assemble.Start(lines)
		if err != nil {
			t.Error(err)
		}

		iterateOverResultWanted(t, wanted)
	})

	t.Run("Testing .ds", func(t *testing.T) {
		lines := []string{" .ds 2, 5 "}
		wanted := []uint8{5, 5}

		err := assemble.Start(lines)
		if err != nil {
			t.Error(err)
		}

		iterateOverResultWanted(t, wanted)
	})

	t.Run("Testing .ines***", func(t *testing.T) {
		lines := []string{" .inesMap 5", " .inesPrg $8000", " .inesChr \"128kb\"", " .inesMir 1"}

		err := assemble.Start(lines)
		if err != nil {
			t.Error(err)
		}
		if romBuilder.GetInesMap() != 5 {
			t.Error("Bad INES Mapper result")
		}
		if romBuilder.GetInesPrgHeaderValue() != 2 {
			t.Error("Bad INES PRG header value result")
		}
		if romBuilder.GetInesPrgSizeInKb() != 0x08000 {
			errMsg := fmt.Sprintf("Bad INES PRG size in KB result: %v", romBuilder.GetInesPrgSizeInKb())
			t.Error(errMsg)
		}
		if romBuilder.GetInesChrHeaderValue() != 16 {
			t.Error("Bad INES CHR result")
		}
		if romBuilder.GetInesChrSizeInKb() != 0x020000 {
			errMsg := fmt.Sprintf("Bad INES CHR size in KB result: Got %v / Want %v", romBuilder.GetInesChrSizeInKb(), 0x020000)
			t.Error(errMsg)
		}
	})

}
