package parser_test

import (
	"misc/nintasm/assemble"
	"misc/nintasm/interpreter/environment"
	"misc/nintasm/romBuilder"
	"testing"
)

func TestLabelOperandParser(t *testing.T) {
	t.Run("***Testing Assignment", func(t *testing.T) {
		var lineCounter uint = 0
		var testLines []string
		var wanted []int

		testLines = append(testLines, "temp1 = 1")
		testLines = append(testLines, "temp2 equ 2")
		wanted = append(wanted, 1)
		wanted = append(wanted, 2)

		err := assemble.ReadLines(&testLines, &lineCounter)
		if err != nil {
			t.Error(err)
		}

		result1, _, _ := environment.LookupIdentifierInSymbolAsNodeTable("temp1")
		if result1.AsNumber != wanted[0] {
			t.Error("= Result not as expected")
			t.Log("Wanted:", wanted[0], "Got:", result1.AsNumber)
		}
		result2, _, _ := environment.LookupIdentifierInSymbolAsNodeTable("temp2")
		if result2.AsNumber != wanted[1] {
			t.Error("EQU Result not as expected")
			t.Log("Wanted:", wanted[1], "Got:", result2.AsNumber)
		}
		romBuilder.ClearINES()
		romBuilder.ClearRom()
	})

	t.Run("***Testing Namespace", func(t *testing.T) {
		var lineCounter uint = 0
		var testLines []string
		var wanted []int

		testLines = append(testLines, "testNS .namespace")
		testLines = append(testLines, ".var1 = 3")
		testLines = append(testLines, ".var2 = 44+2")
		testLines = append(testLines, ".var3 = .var2 - 2")
		testLines = append(testLines, "testNS .endNamespace")
		wanted = append(wanted, 3)
		wanted = append(wanted, 46)
		wanted = append(wanted, 44)

		err := assemble.ReadLines(&testLines, &lineCounter)
		if err != nil {
			t.Error(err)
		}

		result1, _, _ := environment.LookupIdentifierInSymbolAsNodeTable("testNS.var1")
		if result1.AsNumber != wanted[0] {
			t.Error("Namespace 0 result not as expected")
			t.Log("Wanted:", wanted[0], "Got:", result1.AsNumber)
		}
		result2, _, _ := environment.LookupIdentifierInSymbolAsNodeTable("testNS.var2")
		if result2.AsNumber != wanted[1] {
			t.Error("Namespace 1 result not as expected")
			t.Log("Wanted:", wanted[1], "Got:", result2.AsNumber)
		}
		result3, _, _ := environment.LookupIdentifierInSymbolAsNodeTable("testNS.var3")
		if result3.AsNumber != wanted[2] {
			t.Error("Namespace 2 result not as expected")
			t.Log("Wanted:", wanted[2], "Got:", result3.AsNumber)
		}
		romBuilder.ClearINES()
		romBuilder.ClearRom()
	})
}
