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

func AddExprToExprmap(newExpr string, exprValue int) error {
	exprmapSymbolTable[lastAddedExprmapName][newExpr] = exprValue
	return nil
}

// ----------------------------------

func GetCurrentExprmap() (ExprmapTableType, error) {
	if currentExprmapName == "" {
		return nil, errors.New("No exprmaps have been defined!!!!")
	}

	return exprmapSymbolTable[currentExprmapName], nil
}

func GetSpecifiedExprmap(specifiedExprmapName string) (ExprmapTableType, error) {
	if currentExprmapName == "" {
		return nil, errors.New("No exprmaps have been defined!!!!")
	}
	specifiedExprmap, exists := exprmapSymbolTable[specifiedExprmapName]
	if !exists {
		return nil, errors.New("Specified exprmap doesn't exist!")
	}
	return specifiedExprmap, nil
}

// ----------------------------------

func CheckIfDefinedInExprmap(lookupExpr string) (int, bool) {
	target, exists := exprmapSymbolTable[currentExprmapName][lookupExpr]
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

// ----------------------------------

func SetExprmapToDefault() error {
	if defaultExprmapName == "" {
		return errors.New("Cannot reset - No exprmaps have been defined!!!!")
	}
	currentExprmapName = defaultExprmapName
	return nil
}

func SetExprmapTo__(newExprmapName string) error {
	if defaultExprmapName == "" {
		return errors.New("Cannot set - No exprmaps have been defined!!!!")
	}
	_, err := GetSpecifiedExprmap(newExprmapName)
	if err != nil {
		return err
	}

	currentExprmapName = newExprmapName
	return nil
}
