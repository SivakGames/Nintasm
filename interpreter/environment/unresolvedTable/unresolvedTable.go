package unresolvedTable

import (
	"fmt"
	"misc/nintasm/assemble/errorHandler"
	"misc/nintasm/assemble/fileStack"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	"misc/nintasm/interpreter"
	"misc/nintasm/interpreter/environment"
	"misc/nintasm/interpreter/operandFactory"
	"misc/nintasm/romBuilder"
	"misc/nintasm/romBuilder/nodesToBytes"
)

type Node = operandFactory.Node

type unresolvedEntry struct {
	originalRomSegment int
	originalBank       int
	originalOffset     int
	originalNode       Node
	parentLabel        string
	neededBytes        int
	fileName           string
	lineNumber         uint8
	lineContent        string
}

var unresolvedSymbolTable = []unresolvedEntry{}
var unresolvedRomTable = []unresolvedEntry{}

func newUnresolvedEntry(node Node, neededBytes int) unresolvedEntry {
	fileData := fileStack.GetTopOfFileStack()
	return unresolvedEntry{
		originalRomSegment: romBuilder.GetRomSegmentIndex(),
		originalBank:       romBuilder.GetBankIndex(),
		originalOffset:     romBuilder.GetCurrentInsertionIndex(),
		originalNode:       node,
		parentLabel:        interpreter.GetParentLabelNoError(),
		neededBytes:        neededBytes,
		fileName:           fileData.FileName,
		lineNumber:         uint8(fileData.CurrentLineNumber),
		lineContent:        fileData.ProcessedLines[fileData.CurrentLineNumber-1],
	}
}

func AddUnresolvedSymbol(node Node) {
	unresolvedSymbolTable = append(unresolvedSymbolTable, newUnresolvedEntry(node, 0))
	return
}

func AddUnresolvedRomEntry(node Node, mustResolveSize int) {
	unresolvedRomTable = append(unresolvedRomTable, newUnresolvedEntry(node, mustResolveSize))
	return
}

// Called at the end of pass 1
func ResolvedUnresolvedSymbols() error {
	interpreter.ClearParentLabel()

	for len(unresolvedSymbolTable) > 0 {
		originalUnresolvedLength := len(unresolvedSymbolTable)
		newUnresolvedTable := []unresolvedEntry{}
		for _, entry := range unresolvedSymbolTable {
			interpreter.OverwriteParentLabel(entry.parentLabel)
			evaluatedNode, err := interpreter.EvaluateNode(entry.originalNode)
			if err != nil {
				err := errorHandler.CheckErrorContinuesUpwardPropagation(err, enumErrorCodes.UnresolvedIdentifier)
				if err != nil {
					err := errorHandler.CheckErrorContinuesUpwardPropagation(err, enumErrorCodes.Error)
					if err != nil {
						return err // ❌❌ CONTINUES Failing!
					}
				} else {
					entry.originalNode = evaluatedNode
					newUnresolvedTable = append(newUnresolvedTable, entry)
				}
			}
			interpreter.ClearParentLabel()
		}
		newUnresolvedLength := len(newUnresolvedTable)
		if originalUnresolvedLength == newUnresolvedLength {
			fmt.Println(unresolvedSymbolTable)
			panic("🛑 DEADLOCK")
		}
		unresolvedSymbolTable = newUnresolvedTable
	}
	return nil
}

// Called at the end of pass 1
func ResolvedUnresolvedRomEntries() error {
	environment.ClearUnresolvedSilentErrorFlag()
	for _, entry := range unresolvedRomTable {
		errorHandler.OverwriteNoFileDefaults(entry.fileName, uint(entry.lineNumber), entry.lineContent)
		evaluatedNode, err := interpreter.EvaluateNode(entry.originalNode)
		if err != nil {
			err := errorHandler.CheckErrorContinuesUpwardPropagation(err, enumErrorCodes.Error)
			if err != nil {
				return err // ❌❌ CONTINUES Failing!
			}
			continue
		}

		asRomData, err := nodesToBytes.ConvertNodeValueToUInts(evaluatedNode, entry.neededBytes, false, true)
		if err != nil {
			continue
		}
		romBuilder.OverwriteResolvedBytesInRom(entry.originalRomSegment, entry.originalBank, entry.originalOffset, asRomData)
	}

	return nil
}
