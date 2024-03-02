package environment

import (
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	enumSymbolTableTypes "misc/nintasm/constants/enums/symbolTableTypes"
	"misc/nintasm/interpreter/environment/charmapTable"
	"misc/nintasm/interpreter/environment/exprmapTable"
	"misc/nintasm/interpreter/environment/macroTable"
	"misc/nintasm/interpreter/environment/namespaceTable"
	"misc/nintasm/interpreter/environment/symbolAsNodeTable"
	"misc/nintasm/interpreter/operandFactory"
	"misc/nintasm/romBuilder"
)

type Node = operandFactory.Node

var masterLookupTable = map[string]enumSymbolTableTypes.Def{}

// ----------------------------------

func CheckIfAlreadyDefinedInMasterTable(symbolName string) error {
	var exists bool
	_, exists = masterLookupTable[symbolName]
	if exists {
		return errorHandler.AddNew(enumErrorCodes.InterpreterAlreadyDefined, symbolName, "TODO TYPE DESC")
	}
	return nil
}

// ----------------------------------

func addToMasterTable(symbolName string, symbolEnum enumSymbolTableTypes.Def) error {
	masterLookupTable[symbolName] = symbolEnum
	return nil
}

// ============================================================================
// ============================================================================

// Add identifier to the symbol table
func AddIdentifierToSymbolAsNodeTable(symbolName string, node Node) error {
	err := CheckIfAlreadyDefinedInMasterTable(symbolName)
	if err != nil {
		return err
	}
	addToMasterTable(symbolName, enumSymbolTableTypes.SymbolAsNode)
	symbolAsNodeTable.AddIdentifierKeyToSymbolAsNodeTable(symbolName, node)
	return nil
}

// See if identifier has a value in the symbol table
func LookupIdentifierInSymbolAsNodeTable(symbolName string) (Node, error) {
	node, exists := symbolAsNodeTable.GetNodeFromSymbolAsNodeTable(symbolName)
	if !exists {
		_, otherExists := masterLookupTable[symbolName]
		if otherExists {
			return node, errorHandler.AddNew(enumErrorCodes.InterpreterAlreadyDefined, symbolName, "TODO TYPE DESC")
		}
	}
	return node, nil
}

// -----------------------------------------------------------------------------

// In the case of labels, the current bank they are in gets added
func AddToLabelAsBankTable(symbolName string) error {
	bankId := romBuilder.GetBankIndex()
	symbolAsNodeTable.AddIdentifierKeyToLabelAsBankTable(symbolName, bankId)
	return nil
}

// ============================================================================
// ============================================================================

// Add things that do not hold any numeric value(s) to the master env. table and their own respective tables
func AddOtherIdentifierToMasterTable(symbolName string, symbolEnum enumSymbolTableTypes.Def) error {
	err := CheckIfAlreadyDefinedInMasterTable(symbolName)
	if err != nil {
		return err
	}
	addToMasterTable(symbolName, symbolEnum)
	switch symbolEnum {
	case enumSymbolTableTypes.CharMap:
		charmapTable.AddIdentifierKeyToCharmapTable(symbolName)
	case enumSymbolTableTypes.ExprMap:
		exprmapTable.AddIdentifierKeyToExprmapTable(symbolName)
	case enumSymbolTableTypes.KVMacro:
		macroTable.AddIdentifierKeyToKVMacroTable(symbolName)
	case enumSymbolTableTypes.Macro:
		macroTable.AddIdentifierKeyToMacroTable(symbolName)
	case enumSymbolTableTypes.Namespace:
		namespaceTable.AddIdentifierKeyToNamespaceTable(symbolName)
	default:
		panic("Bad symbol type being added to environment!")
	}
	return nil
}
