package parser_test

import (
	"misc/nintasm/assemble"
	"misc/nintasm/romBuilder"
	"testing"
)

func TestDirectiveOperandParser(t *testing.T) {
	t.Run("Testing .db", func(t *testing.T) {

		testLine1 := " .db 54"
		wanted1 := uint8(55)

		lines := []string{testLine1}
		err := assemble.Start(lines)
		if err != nil {
			t.Error(err)
		}

		bankSegment := romBuilder.GetCurrentBankSegment()
		result1 := (*bankSegment)[0]

		if result1 != wanted1 {
			t.Error("= Result not as expected")
			t.Log("Wanted:", wanted1, "Got:", result1)
		}

	})

}
