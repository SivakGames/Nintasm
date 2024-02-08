package parser_test

import (
	"misc/nintasm/assemble"
	"misc/nintasm/interpreter/environment"
	"testing"
)

func TestLabelOperandParser(t *testing.T) {
	t.Run("Testing Assignment", func(t *testing.T) {

		testLine1 := "temp1 = 1"
		testLine2 := "temp2 equ 2"
		wanted1 := 1
		wanted2 := 2

		lines := []string{testLine1, testLine2}
		err := assemble.Start(lines)
		if err != nil {
			t.Error(err)
		}

		result1, _ := environment.LookupInEnvironment("temp1")
		if result1.AsNumber != wanted1 {
			t.Error("= Result not as expected")
			t.Log("Wanted:", wanted1, "Got:", result1.AsNumber)
		}
		result2, _ := environment.LookupInEnvironment("temp2")
		if result2.AsNumber != wanted2 {
			t.Error("EQU Result not as expected")
			t.Log("Wanted:", wanted2, "Got:", result2.AsNumber)
		}
	})

}
