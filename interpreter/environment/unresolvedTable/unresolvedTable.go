package unresolvedTable

import (
	"fmt"
	"misc/nintasm/assemble/fileStack"
	"misc/nintasm/interpreter"
	"misc/nintasm/interpreter/operandFactory"
	"misc/nintasm/romBuilder"
)

type Node = operandFactory.Node

type unresolvedRomEntry struct {
	originalRomSegment int
	originalBank       int
	originalPrgCounter int
	originalNode       Node
	mustResolveSize    int
	fileName           string
	lineNumber         uint8
	lineContent        string
}

var unresolvedRomTable = []unresolvedRomEntry{}

func newUnresolvedEntry(node Node, mustResolveSize int) unresolvedRomEntry {
	fileData := fileStack.GetTopOfFileStack()
	return unresolvedRomEntry{
		originalRomSegment: romBuilder.GetRomSegmentIndex(),
		originalBank:       romBuilder.GetBankIndex(),
		originalPrgCounter: romBuilder.GetCurrentInsertionIndex(),
		originalNode:       node,
		mustResolveSize:    mustResolveSize,
		fileName:           fileData.FileName,
		lineNumber:         uint8(fileData.CurrentLineNumber),
		lineContent:        fileData.ProcessedLines[fileData.CurrentLineNumber-1],
	}
}

func AddUnresolvedRomEntry(node Node, mustResolveSize int) {
	unresolvedRomTable = append(unresolvedRomTable, newUnresolvedEntry(node, mustResolveSize))
	return
}

// Called at the end of pass 1
func ResolvedUnresolved() {
	for len(unresolvedRomTable) > 0 {
		originalUnresolvedLength := len(unresolvedRomTable)
		newUnresolvedTable := []unresolvedRomEntry{}
		for _, entry := range unresolvedRomTable {
			evaluatedNode, err := interpreter.EvaluateNode(entry.originalNode)
			if err != nil {
				panic("Reeval fail")
			}
			if evaluatedNode.Resolved {
				fmt.Println("Reeval OK", entry, evaluatedNode)
			}
		}
		newUnresolvedLength := len(newUnresolvedTable)
		fmt.Println(originalUnresolvedLength == newUnresolvedLength)
		unresolvedRomTable = newUnresolvedTable
	}

	//len(unresolvedTable)

}
