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
		return errorHandler.AddNew(enumErrorCodes.MacroMisplacedExitMacro)
	}
	blockStack.SetExitOpName("IM")

	return nil
}
