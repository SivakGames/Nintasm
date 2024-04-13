package interpreter

import (
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
)

//++++++++++++++++++++++++++++++++++++++++++++++++++++++++

var PopParentLabelWhenBlockOpDone bool = false
var parentLabelStack []string
var currentLocalLabel string

//++++++++++++++++++++++++++++++++++++++++++++++++++++++++

func GetParentLabel() (string, error) {
	if len(parentLabelStack) == 0 {
		return "", errorHandler.AddNew(enumErrorCodes.InterpreterNoParentLabel)
	}
	parentLabel := parentLabelStack[len(parentLabelStack)-1]
	return parentLabel, nil
}

func GetParentLabelNoError() string {
	if len(parentLabelStack) == 0 {
		return ""
	}
	parentLabel := parentLabelStack[len(parentLabelStack)-1]
	return parentLabel
}

// Will overwrite at current position or add if none
func OverwriteParentLabel(newLabel string) {
	if len(parentLabelStack) == 0 {
		parentLabelStack = append(parentLabelStack, newLabel)
		return
	}
	parentLabelStack[len(parentLabelStack)-1] = newLabel
}

func ClearParentLabel() {
	parentLabelStack = parentLabelStack[:0]
}

// ------------------------------------------------------------------

func AppendParentLabel(newLabel string) {
	parentLabelStack = append(parentLabelStack, newLabel)
}

func PopParentLabel() {
	parentLabelStack = parentLabelStack[:len(parentLabelStack)-1]
}

// ------------------------------------------------------------------
func ClearLocalLabel() {
	currentLocalLabel = ""
}
func GetLocalLabel() string {
	return currentLocalLabel
}
func SetLocalLabel(newLabel string) {
	currentLocalLabel = newLabel
}
