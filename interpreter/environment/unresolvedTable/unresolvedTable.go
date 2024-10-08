package unresolvedTable

import (
	"misc/nintasm/assemble/errorHandler"
	"misc/nintasm/assemble/fileStack"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	"misc/nintasm/interpreter"
	"misc/nintasm/interpreter/environment"
	"misc/nintasm/interpreter/environment/symbolAsNodeTable"
	"misc/nintasm/interpreter/operandFactory"
	"misc/nintasm/romBuilder"
	"misc/nintasm/romBuilder/nodesToBytes"
)

type Node = operandFactory.Node

type childScopesType []*symbolAsNodeTable.SymbolTableType

type unresolvedEntry struct {
	originalRomSegment int
	originalBank       int
	originalOffset     int
	originalNode       Node
	parentLabel        string
	localLabel         string
	neededBytes        int
	fileName           string
	lineNumber         uint8
	lineContent        string
	isBranch           bool
	isBigEndian        bool
	childScopes        childScopesType
}

var unresolvedSymbolTable = []unresolvedEntry{}
var unresolvedRomTable = []unresolvedEntry{}

func newUnresolvedEntry(node Node, neededBytes int, isBranch bool, isBigEndian bool) unresolvedEntry {
	currentChildScopes := symbolAsNodeTable.DeepCopyLocalBlockScopes()
	fileData := fileStack.GetTopOfFileStack()

	return unresolvedEntry{
		originalRomSegment: romBuilder.GetRomSegmentIndex(),
		originalBank:       romBuilder.GetBankIndex(),
		originalOffset:     romBuilder.GetCurrentInsertionIndex(),
		originalNode:       node,
		parentLabel:        interpreter.GetParentLabelNoError(),
		localLabel:         interpreter.GetLocalLabel(),
		neededBytes:        neededBytes,
		fileName:           fileData.FileName,
		lineNumber:         uint8(fileData.CurrentLineNumber),
		lineContent:        fileData.ProcessedLines[fileData.CurrentLineNumber-1],
		isBranch:           isBranch,
		isBigEndian:        isBigEndian,
		childScopes:        currentChildScopes,
	}
}

func AddUnresolvedSymbol(node Node) {
	unresolvedSymbolTable = append(unresolvedSymbolTable, newUnresolvedEntry(node, 0, false, false))
	return
}

func AddUnresolvedRomEntry(node Node, mustResolveSize int, isBranch bool, isBigEndian bool) {
	unresolvedRomTable = append(unresolvedRomTable, newUnresolvedEntry(node, mustResolveSize, isBranch, isBigEndian))
	return
}

// Called at the end of pass 1
func ResolvedUnresolvedSymbols() error {
	interpreter.ClearParentLabel()
	interpreter.ClearLocalLabel()
	totalUnresolved := len(unresolvedSymbolTable)

	for len(unresolvedSymbolTable) > 0 {
		originalUnresolvedLength := len(unresolvedSymbolTable)
		unresolvedDiff := float64(totalUnresolved - originalUnresolvedLength)
		pass2Progress := (unresolvedDiff / float64(totalUnresolved)) * 50
		romBuilder.DrawPass2Progress(pass2Progress)

		newUnresolvedTable := []unresolvedEntry{}
		for _, entry := range unresolvedSymbolTable {
			errorHandler.OverwriteNoFileDefaults(entry.fileName, uint(entry.lineNumber), entry.lineContent)
			interpreter.OverwriteParentLabel(entry.parentLabel)
			interpreter.SetLocalLabel(entry.localLabel)
			symbolAsNodeTable.SetCurrentLocalBlockScopes(entry.childScopes)

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
			interpreter.ClearLocalLabel()
			interpreter.ClearParentLabel()
			symbolAsNodeTable.ClearCurrentLocalBlockScopes()
		}

		//If equal, then no more resolving can be done
		newUnresolvedLength := len(newUnresolvedTable)
		if originalUnresolvedLength == newUnresolvedLength {
			environment.ClearUnresolvedSilentErrorFlag()
			for _, entry := range newUnresolvedTable {
				errorHandler.OverwriteNoFileDefaults(entry.fileName, uint(entry.lineNumber), entry.lineContent)
				interpreter.EvaluateNode(entry.originalNode)
			}
			return errorHandler.AddNew(enumErrorCodes.ResolveDeadlock)
		}
		unresolvedSymbolTable = newUnresolvedTable
	}
	romBuilder.DrawPass2Progress(50)
	return nil
}

// Called at the end of pass 1
func ResolvedUnresolvedRomEntries() error {
	environment.ClearUnresolvedSilentErrorFlag()
	totalSuccessfullyResolved := 0
	totalUnresolved := len(unresolvedRomTable)
	resolvedAddPercent := 50.0

	for _, entry := range unresolvedRomTable {
		interpreter.OverwriteParentLabel(entry.parentLabel)
		interpreter.SetLocalLabel(entry.localLabel)
		symbolAsNodeTable.SetCurrentLocalBlockScopes(entry.childScopes)

		errorHandler.OverwriteNoFileDefaults(entry.fileName, uint(entry.lineNumber), entry.lineContent)
		evaluatedNode, err := interpreter.EvaluateNode(entry.originalNode)
		if err != nil {
			err := errorHandler.CheckErrorContinuesUpwardPropagation(err, enumErrorCodes.Error)
			if err != nil {
				return err // ❌❌ CONTINUES Failing!
			}
			continue
		}

		asRomData, err := nodesToBytes.ConvertNodeValueToUInts(evaluatedNode, entry.neededBytes, entry.isBranch, entry.isBigEndian, true)
		if err != nil {
			continue
		}
		romBuilder.OverwriteResolvedBytesInRom(entry.originalRomSegment, entry.originalBank, entry.originalOffset, asRomData)
		interpreter.ClearParentLabel()
		interpreter.ClearLocalLabel()
		symbolAsNodeTable.ClearCurrentLocalBlockScopes()
		totalSuccessfullyResolved++
	}

	if totalUnresolved > 0 {
		resolvedAddPercent *= (float64(totalSuccessfullyResolved) / float64(totalUnresolved))
	}

	romBuilder.DrawPass2Progress(50 + resolvedAddPercent)

	return nil
}
