package charmapTable

import (
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
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

// When first defining a charmap
func AddIdentifierKeyToCharmapTable(symbolName string) {
	charmapSymbolTable[symbolName] = CharmapTableType{}
	lastAddedCharmapName = symbolName
	if defaultCharmapName == "" {
		defaultCharmapName = lastAddedCharmapName
		currentCharmapName = defaultCharmapName
	}
}

func AddCharToCharmap(newChar rune, charNodes []Node) error {
	charmapSymbolTable[lastAddedCharmapName][newChar] = charNodes
	return nil
}

// ----------------------------------

func GetCurrentCharmap() (CharmapTableType, error) {
	if currentCharmapName == "" {
		return nil, errorHandler.AddNew(enumErrorCodes.CharMapNoneDefined) // ❌ Fails
	}
	return charmapSymbolTable[currentCharmapName], nil
}

func GetSpecifiedCharmap(specifiedCharmapName string) (CharmapTableType, error) {
	if currentCharmapName == "" {
		return nil, errorHandler.AddNew(enumErrorCodes.CharMapNoneDefined) // ❌ Fails
	}
	specifiedCharmap, exists := charmapSymbolTable[specifiedCharmapName]
	if !exists {
		return nil, errorHandler.AddNew(enumErrorCodes.CharMapNotExist) // ❌ Fails
	}
	return specifiedCharmap, nil
}

// ----------------------------------

func MapStringToCharmap(stringToConvert string) ([]int, error) {
	stringAsRuneArray := []rune(stringToConvert)
	replacedStringAsBytes := []int{}

	currCharmap, err := GetCurrentCharmap()
	if err != nil {
		return replacedStringAsBytes, err
	}

	for _, r := range stringAsRuneArray {
		nodes, exists := currCharmap[r]
		if !exists {
			return replacedStringAsBytes, errorHandler.AddNew(enumErrorCodes.ToCharMapUndefChar, r, currentCharmapName) // ❌ Fails
		}
		for _, v := range nodes {
			replacedStringAsBytes = append(replacedStringAsBytes, int(v.AsNumber))
		}
	}

	return replacedStringAsBytes, nil
}

// ----------------------------------

func checkIfDefinedInCharmap(lookupChar rune) ([]Node, bool) {
	target, exists := charmapSymbolTable[lastAddedCharmapName][lookupChar]
	return target, exists
}

func CheckIfCharAlreadyExistsInCharmap(lookupChar rune) ([]Node, error) {
	target, exists := checkIfDefinedInCharmap(lookupChar)
	if exists {
		return target, errorHandler.AddNew(enumErrorCodes.CharMapDuplicateKey, lookupChar, lastAddedCharmapName) // ❌ Fails
	}
	return target, nil
}

// ----------------------------------

func SetCharmapToDefault() error {
	if defaultCharmapName == "" {
		return errorHandler.AddNew(enumErrorCodes.CharMapNoneDefined) // ❌ Fails
	}
	currentCharmapName = defaultCharmapName
	return nil
}

func SetCharmapTo__(newCharmapName string) error {
	if defaultCharmapName == "" {
		return errorHandler.AddNew(enumErrorCodes.CharMapNoneDefined) // ❌ Fails
	}
	_, err := GetSpecifiedCharmap(newCharmapName)
	if err != nil {
		return err
	}

	currentCharmapName = newCharmapName
	return nil
}
