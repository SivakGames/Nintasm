package charmapTable

import (
	"errors"
	"fmt"
	"misc/nintasm/interpreter/operandFactory"
)

// ++++++++++++++++++++++++++++++++++++++++++++++++++++++

type Node = operandFactory.Node

type CharmapTableType = map[rune][]Node

var charmapSymbolTable = map[string]CharmapTableType{}

var defaultCharmapName = ""
var currentCharmapName = ""
var lastAddedCharmapName = ""

// ----------------------------------

func AddCharmapToEnvironment(symbolName string) error {
	charmapSymbolTable[symbolName] = CharmapTableType{}
	lastAddedCharmapName = symbolName

	if defaultCharmapName == "" {
		defaultCharmapName = lastAddedCharmapName
		currentCharmapName = defaultCharmapName
	}
	return nil
}

func GetCurrentCharmap() (CharmapTableType, error) {
	if currentCharmapName == "" {
		return nil, errors.New("No charmaps have been defined!!!!")
	}

	return charmapSymbolTable[currentCharmapName], nil
}

func CheckIfDefinedInCharmap(lookupChar rune) ([]Node, bool) {
	target, exists := charmapSymbolTable[lastAddedCharmapName][lookupChar]
	return target, exists
}

func CheckIfAlreadyExistsInCharmap(lookupChar rune) ([]Node, error) {
	target, exists := CheckIfDefinedInCharmap(lookupChar)
	if exists {
		errMsg := fmt.Sprintf("Target char %c already defined!", lookupChar)
		return target, errors.New(errMsg)
	}
	return target, nil
}

func AddCharToCharmap(newChar rune, charNodes []Node) error {
	charmapSymbolTable[lastAddedCharmapName][newChar] = charNodes
	return nil
}
