package blockStack

import "fmt"

// +++++++++++++++++++++++++++++++++++++++++++++

//Used by labeled directives. Backs up the label when initially opening a block
var currentBlockOperationLabel string = ""

// +++++++++++++++++++++++++++++++++++++++++++++++++++
func ClearCurrentOperationLabel() {
	currentBlockOperationLabel = ""
}
func GetCurrentOperationLabel() string {
	return currentBlockOperationLabel
}

// Will set the label of the labeled operation that will be captured.
// If one was previously set then error because it hasn't finished.
func SetCurrentOperationLabel(label string) error {
	if currentBlockOperationLabel != "" {
		panic(fmt.Sprintf("🛑 Somehow entering another label block operation while first (%v) is not done...", currentBlockOperationLabel))
	}
	currentBlockOperationLabel = label
	return nil
}
