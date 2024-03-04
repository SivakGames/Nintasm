package exprmapTable

import (
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
)

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++

type ExprmapTableType = map[string]int

var exprmapSymbolTable = map[string]ExprmapTableType{}

var defaultExprmapName = ""
var currentExprmapName = ""
var lastAddedExprmapName = ""

// ----------------------------------

// When first defining a charmap
func AddIdentifierKeyToExprmapTable(symbolName string) error {
	exprmapSymbolTable[symbolName] = ExprmapTableType{}
	lastAddedExprmapName = symbolName

	if defaultExprmapName == "" {
		defaultExprmapName = lastAddedExprmapName
		currentExprmapName = defaultExprmapName
	}
	return nil
}

func AddExprToExprmap(newExpr string, exprValue int) error {
	exprmapSymbolTable[lastAddedExprmapName][newExpr] = exprValue
	return nil
}

// ----------------------------------

func GetCurrentExprmap() (ExprmapTableType, error) {
	if currentExprmapName == "" {
		return nil, errorHandler.AddNew(enumErrorCodes.ExprMapNoneDefined) // ❌ Fails
	}

	return exprmapSymbolTable[currentExprmapName], nil
}

func GetSpecifiedExprmap(specifiedExprmapName string) (ExprmapTableType, error) {
	if currentExprmapName == "" {
		return nil, errorHandler.AddNew(enumErrorCodes.ExprMapNoneDefined) // ❌ Fails
	}
	specifiedExprmap, exists := exprmapSymbolTable[specifiedExprmapName]
	if !exists {
		return nil, errorHandler.AddNew(enumErrorCodes.ExprMapNotExist) // ❌ Fails
	}
	return specifiedExprmap, nil
}

// ----------------------------------

func CheckIfDefinedInExprmap(lookupExpr string) (int, bool) {
	target, exists := exprmapSymbolTable[currentExprmapName][lookupExpr]
	return target, exists
}

func CheckIfAlreadyExistsInExprmap(lookupExpr string) (int, error) {
	target, exists := exprmapSymbolTable[lastAddedExprmapName][lookupExpr]
	if exists {
		return target, errorHandler.AddNew(enumErrorCodes.ExprMapDuplicateKey, lookupExpr) // ❌ Fails
	}
	return target, nil
}

// ----------------------------------

func SetExprmapToDefault() error {
	if defaultExprmapName == "" {
		return errorHandler.AddNew(enumErrorCodes.ExprMapNoneDefined)
	}
	currentExprmapName = defaultExprmapName
	return nil
}

func SetExprmapTo__(newExprmapName string) error {
	if defaultExprmapName == "" {
		return errorHandler.AddNew(enumErrorCodes.ExprMapNoneDefined) // ❌ Fails
	}
	_, err := GetSpecifiedExprmap(newExprmapName)
	if err != nil {
		return err
	}

	currentExprmapName = newExprmapName
	return nil
}
