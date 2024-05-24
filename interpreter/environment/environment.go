package environment

import (
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	enumSymbolTableTypes "misc/nintasm/constants/enums/symbolTableTypes"
	"misc/nintasm/interpreter/environment/charmapTable"
	"misc/nintasm/interpreter/environment/exprmapTable"
	"misc/nintasm/interpreter/environment/funcTable"
	"misc/nintasm/interpreter/environment/macroTable"
	"misc/nintasm/interpreter/environment/namespaceTable"
	"misc/nintasm/interpreter/environment/symbolAsNodeTable"
	"misc/nintasm/interpreter/operandFactory"
	"misc/nintasm/romBuilder"
	"strings"
)

type Node = operandFactory.Node

var masterLookupTable = map[string]enumSymbolTableTypes.Def{}
var unresolvedAddsSilentError = true

// ----------------------------------

func CheckIfIdentifierExistsInMasterTable(symbolName string) (enumSymbolTableTypes.Def, bool) {
	symbolTableEnum, exists := masterLookupTable[symbolName]
	return symbolTableEnum, exists
}

func CheckIfAlreadyDefinedInMasterTable(symbolName string) error {
	symbolEnum, exists := masterLookupTable[symbolName]
	if exists {
		description := GetSymbolDescriptionFromEnum(symbolEnum)
		return errorHandler.AddNew(enumErrorCodes.InterpreterAlreadyDefined, symbolName, description)
	}
	return nil
}

// ----------------------------------

func addToMasterTable(symbolName string, symbolEnum enumSymbolTableTypes.Def) error {
	masterLookupTable[symbolName] = symbolEnum
	return nil
}
func removeFromMasterTable(symbolName string) error {
	delete(masterLookupTable, symbolName)
	return nil
}

// ============================================================================
// ============================================================================

// Add identifier to the symbol table
func AddIdentifierToSymbolAsNodeTable(symbolName string, node Node, symbolType enumSymbolTableTypes.Def) error {
	err := CheckIfAlreadyDefinedInMasterTable(symbolName)
	if err != nil {
		return err
	}
	addToMasterTable(symbolName, symbolType)
	symbolAsNodeTable.AddIdentifierKeyToSymbolAsNodeTable(symbolName, node)
	return nil
}

// See if identifier has a value in the symbol table.
// Returns node, resolved, err
func LookupIdentifierInSymbolAsNodeTable(symbolName string) (Node, bool, error) {
	node, exists := symbolAsNodeTable.GetNodeFromSymbolAsNodeTable(symbolName)
	if !exists {
		//Sees if the label exists but doesn't correspond to an actual value
		_, nonValueLabelExists := masterLookupTable[symbolName]
		if nonValueLabelExists {
			return node, false, errorHandler.AddNew(enumErrorCodes.InterpreterIdentifierNotValueSymbol, symbolName)
		}
		if unresolvedAddsSilentError {
			return node, false, errorHandler.AddUnresolved()
		}
		return node, false, errorHandler.AddNew(enumErrorCodes.InterpreterSymbolNotFound, symbolName)
	}
	return node, true, nil
}

// Get all local labels with a parent label. Used for generating namespace indexes (GNSI)
func GetLocalLabelsOfParent(symbolName string) []string {
	localLabelsFromParent := make([]string, 0)
	for key := range masterLookupTable {
		if strings.HasPrefix(key, symbolName+".") {
			localLabelsFromParent = append(localLabelsFromParent, key)
		}
	}
	return localLabelsFromParent
}

// -----------------------------------------------------------------------------

// Will check the topmost stack for a substitution ID and get the node if there is one
func LookupSubstitutionID(symbolName string) (Node, error) {
	substitutionNode, exists := symbolAsNodeTable.LookupSymbolInTopOfSymbolTableStack(symbolName)
	if !exists {
		if unresolvedAddsSilentError {
			return operandFactory.CreateUndefinedNode(symbolName), errorHandler.AddUnresolved()
		}
		return substitutionNode, errorHandler.AddNew(enumErrorCodes.InterpreterSymbolNotFound, symbolName)
	}
	return substitutionNode, nil
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
	case enumSymbolTableTypes.Function:
		funcTable.AddIdentifierKeyToFunctionTable(symbolName)
	case enumSymbolTableTypes.Macro:
		macroTable.AddIdentifierKeyToMacroTable(symbolName)
	case enumSymbolTableTypes.Namespace:
		namespaceTable.AddIdentifierKeyToNamespaceTable(symbolName)
	default:
		panic("🛑 Bad symbol type being added to environment!")
	}
	return nil
}

// ============================================================================
// ============================================================================

// Add things that do not hold any numeric value(s) to the master env. table and their own respective tables
func RemoveOtherIdentifierFromMasterTable(symbolName string, symbolEnum enumSymbolTableTypes.Def) error {
	identifierEnum, identifierExists := CheckIfIdentifierExistsInMasterTable(symbolName)

	if !identifierExists {
		return errorHandler.AddNew(enumErrorCodes.RemovedIdentifierNotFound, symbolName)
	}
	if identifierEnum != symbolEnum {
		return errorHandler.AddNew(enumErrorCodes.RemovedIdentifierWrongType, symbolName, GetSymbolDescriptionFromEnum(identifierEnum), GetSymbolDescriptionFromEnum(symbolEnum))
	}

	removeFromMasterTable(symbolName)
	switch symbolEnum {
	case enumSymbolTableTypes.Macro:
		macroTable.RemoveIdentifierKeyFromMacroTable(symbolName)

	default:
		panic("🛑 Bad symbol type being removed from environment!")
	}
	return nil
}

// ++++++++++++++++++++++++++++++
func ClearUnresolvedSilentErrorFlag() {
	unresolvedAddsSilentError = false
}
func GetUnresolvedSilentErrorFlag() bool {
	return unresolvedAddsSilentError
}
func SetUnresolvedSilentErrorFlagTo(value bool) {
	unresolvedAddsSilentError = value
}
func SetUnresolvedSilentErrorFlag() {
	unresolvedAddsSilentError = true
}

func GetSymbolDescriptionFromEnum(symbolEnum enumSymbolTableTypes.Def) string {
	switch symbolEnum {
	case enumSymbolTableTypes.CharMap:
		return "CHARMAP"
	case enumSymbolTableTypes.ExprMap:
		return "EXPRMAP"
	case enumSymbolTableTypes.Function:
		return "FUNC"
	case enumSymbolTableTypes.Label:
		return "LABEL"
	case enumSymbolTableTypes.Macro:
		return "MACRO"
	case enumSymbolTableTypes.Namespace:
		return "NAMESPACE"
	case enumSymbolTableTypes.SymbolAsNode:
		return "SYMBOL"
	default:
		return "??????????"
	}
}
