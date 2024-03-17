package interpreter

import (
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
)

//++++++++++++++++++++++++++++++++++++++++++++++++++++++++

var PopParentLabelWhenBlockOpDone bool = false
var parentLabelStack []string

//++++++++++++++++++++++++++++++++++++++++++++++++++++++++

func GetParentLabel() (string, error) {
	if len(parentLabelStack) == 0 {
		return "", errorHandler.AddNew(enumErrorCodes.InterpreterNoParentLabel)
	}
	parentLabel := parentLabelStack[len(parentLabelStack)-1]
	return parentLabel, nil
}

// Will overwrite at current position or add if none
func OverwriteParentLabel(newLabel string) {
	if len(parentLabelStack) == 0 {
		parentLabelStack = append(parentLabelStack, newLabel)
		return
	}
	parentLabelStack[len(parentLabelStack)-1] = newLabel
	return
}

// ------------------------------------------------------------------

func AppendParentLabel(newLabel string) {
	parentLabelStack = append(parentLabelStack, newLabel)
	return
}

func PopParentLabel() {
	parentLabelStack = parentLabelStack[:len(parentLabelStack)-1]
	return
}
