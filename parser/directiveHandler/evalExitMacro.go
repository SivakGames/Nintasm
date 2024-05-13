package directiveHandler

import (
	"fmt"
	"misc/nintasm/assemble/blockStack"
)

func evalExitMacro() error {
	currentOpPtr := blockStack.GetCurrentActiveOpPtr()
	blockEntries := blockStack.GetBlockEntriesWithPtr(currentOpPtr)
	if len(*blockEntries) == 0 {
		fmt.Println("NO!")
	}
	blockStack.SetExitOpName("IM")

	return nil
}
