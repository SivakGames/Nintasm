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

func AddCharToCharmap(newChar rune, charNodes []Node) error {
	charmapSymbolTable[lastAddedCharmapName][newChar] = charNodes
	return nil
}

// ----------------------------------

func GetCurrentCharmap() (CharmapTableType, error) {
	if currentCharmapName == "" {
		return nil, errors.New("No charmaps have been defined!!!!")
	}
	return charmapSymbolTable[currentCharmapName], nil
}

func GetSpecifiedCharmap(specifiedCharmapName string) (CharmapTableType, error) {
	if currentCharmapName == "" {
		return nil, errors.New("No charmaps have been defined!!!!")
	}
	specifiedCharmap, exists := charmapSymbolTable[specifiedCharmapName]
	if !exists {
		return nil, errors.New("Specified charmap doesn't exist!")
	}
	return specifiedCharmap, nil
}

// ----------------------------------

func MapStringToCharmap(stringToConvert string) (string, error) {
	replacedString := ""
	stringAsRuneArray := []rune(stringToConvert)

	currCharmap, err := GetCurrentCharmap()
	if err != nil {
		return replacedString, err
	}

	for _, r := range stringAsRuneArray {
		nodes, exists := currCharmap[r]
		if !exists {
			errMsg := fmt.Sprintf("Character `%c` is not defined in currently used charmap %v", r, currentCharmapName)
			return replacedString, errors.New(errMsg)
		}
		for _, v := range nodes {
			replacedString += string(rune(v.AsNumber))
		}
	}

	return replacedString, nil
}

// ----------------------------------

func checkIfDefinedInCharmap(lookupChar rune) ([]Node, bool) {
	target, exists := charmapSymbolTable[lastAddedCharmapName][lookupChar]
	return target, exists
}

func CheckIfCharAlreadyExistsInCharmap(lookupChar rune) ([]Node, error) {
	target, exists := checkIfDefinedInCharmap(lookupChar)
	if exists {
		errMsg := fmt.Sprintf("Target char %c already defined!", lookupChar)
		return target, errors.New(errMsg)
	}
	return target, nil
}

// ----------------------------------

func SetCharmapToDefault() error {
	if defaultCharmapName == "" {
		return errors.New("Cannot reset - No charmaps have been defined!!!!")
	}
	currentCharmapName = defaultCharmapName
	return nil
}

func SetCharmapTo__(newCharmapName string) error {
	if defaultCharmapName == "" {
		return errors.New("Cannot set - No charmaps have been defined!!!!")
	}
	_, err := GetSpecifiedCharmap(newCharmapName)
	if err != nil {
		return err
	}

	currentCharmapName = newCharmapName
	return nil
}
