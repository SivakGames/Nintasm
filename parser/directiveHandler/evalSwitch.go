package directiveHandler

import (
	"misc/nintasm/assemble/blockStack"
)

func evalSwitch(directiveName string, operandList *[]Node) error {
	blockStack.PushOntoTopEntry(directiveName, *operandList)
	return nil
}

func evalEndSwitch() error {
	return nil
}
