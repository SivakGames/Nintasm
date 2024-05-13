package directiveHandler

import (
	"misc/nintasm/assemble/blockStack"
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
)

func evalExitMacro() error {
	currentOpPtr := blockStack.GetCurrentActiveOpPtr()
	blockEntries := blockStack.GetBlockEntriesWithPtr(currentOpPtr)
	if len(*blockEntries) == 0 {
		return errorHandler.AddNew(enumErrorCodes.Other, "Exit Macro without being in a macro")
	}
	blockStack.SetExitOpName("IM")

	return nil
}
