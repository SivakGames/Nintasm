package parser_test

import (
	"fmt"
	testHelper "misc/nintasm/__tests__/helper"
	"misc/nintasm/assemble"
	"misc/nintasm/romBuilder"
	"testing"
)

func TestDirectiveOperandParser(t *testing.T) {
	t.Run("Testing .ines*** and declaring segments", func(t *testing.T) {
		err := assemble.Start(testHelper.BaseLines)
		if err != nil {
			t.Error(err)
		}

		t.Log("--- Doing INES Header tests ---")

		if romBuilder.GetInesMap() != 0 {
			t.Error("Bad INES Mapper result")
		}
		if romBuilder.GetInesPrgHeaderValue() != 2 {
			t.Error("Bad INES PRG header value result")
		}
		if romBuilder.GetInesPrgSizeInKb() != 0x08000 {
			errMsg := fmt.Sprintf("Bad INES PRG size in KB result: %v", romBuilder.GetInesPrgSizeInKb())
			t.Error(errMsg)
		}
		if romBuilder.GetInesChrHeaderValue() != 1 {
			t.Error("Bad INES CHR header value result")
		}
		if romBuilder.GetInesChrSizeInKb() != 0x02000 {
			errMsg := fmt.Sprintf("Bad INES CHR size in KB result: Got %v / Want %v", romBuilder.GetInesChrSizeInKb(), 0x020000)
			t.Error(errMsg)
		}
		if romBuilder.GetTotalRomSegmentsInRom() != 1 {
			errMsg := fmt.Sprintf("Unexpected ROM Segment length: Got %v / Want %v", romBuilder.GetTotalRomSegmentsInRom(), 1)
			t.Error(errMsg)
		}
		if romBuilder.GetTotalBanksInCurrentRomSegment() != 2 {
			errMsg := fmt.Sprintf("Unexpected ROM Segment bank quantity length: Got %v / Want %v", romBuilder.GetTotalBanksInCurrentRomSegment(), 2)
			t.Error(errMsg)
		}
		if len(*romBuilder.GetCurrentBankSegmentBytes()) != 0x004000 {
			errMsg := fmt.Sprintf("Wrong number of bytes for bank. Got %v / Want %v", len(*romBuilder.GetCurrentBankSegmentBytes()), 0x004000)
			t.Error(errMsg)
		}
		if romBuilder.GetOrg() != 0x008000 {
			errMsg := fmt.Sprintf("Unexpected ORG: Got %v / Want %v", romBuilder.GetOrg(), 0x008000)
			t.Error(errMsg)
		}
		romBuilder.ClearINES()
		romBuilder.ClearRom()
	})

}
