package exprmapTable

import (
	"errors"
	"fmt"
)

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++

type ExprmapTableType = map[string]int

var exprmapSymbolTable = map[string]ExprmapTableType{}

var defaultExprmapName = ""
var currentExprmapName = ""
var lastAddedExprmapName = ""

// ----------------------------------

func AddExprmapToEnvironment(symbolName string) error {
	exprmapSymbolTable[symbolName] = ExprmapTableType{}
	lastAddedExprmapName = symbolName

	if defaultExprmapName == "" {
		defaultExprmapName = lastAddedExprmapName
		currentExprmapName = defaultExprmapName
	}
	return nil
}

func GetCurrentExprmap() (ExprmapTableType, error) {
	if currentExprmapName == "" {
		return nil, errors.New("No exprmaps have been defined!!!!")
	}

	return exprmapSymbolTable[currentExprmapName], nil
}

func CheckIfDefinedInExprmap(lookupExpr string) (int, bool) {
	target, exists := exprmapSymbolTable[lastAddedExprmapName][lookupExpr]
	return target, exists
}

func CheckIfAlreadyExistsInExprmap(lookupExpr string) (int, error) {
	target, exists := CheckIfDefinedInExprmap(lookupExpr)
	if exists {
		errMsg := fmt.Sprintf("Target expr %v already defined in exprmap!", lookupExpr)
		return target, errors.New(errMsg)
	}
	return target, nil
}

func AddExprToExprmap(newExpr string, exprValue int) error {
	exprmapSymbolTable[lastAddedExprmapName][newExpr] = exprValue
	return nil
}
