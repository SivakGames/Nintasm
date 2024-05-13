package blockStack

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
// Dealing with capture list itself
// <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<

// When processing block ops (macros in particular), add another capture list
func AddNewCaptureBlockListNode() {
	createAndAppendNewCaptureBlockList()
}

func DestroyCaptureBlockListNodeWithPointer(ptr *CaptureBlockListNode) {
	destroyCaptureBlockListAtPointer(ptr)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++
func ClearCaptureBlockListEvalFlag() {
	blockList := getCurrentCaptureBlockListNode()
	blockList.evalutesInsteadOfCapturing = false
}
func GetCaptureBlockListEvalFlag() bool {
	blockList := getCurrentCaptureBlockListNode()
	return blockList.evalutesInsteadOfCapturing
}
func SetCaptureBlockListEvalFlag() {
	blockList := getCurrentCaptureBlockListNode()
	blockList.evalutesInsteadOfCapturing = true
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++
func ClearCaptureBlockListForcedCapturingFlag() {
	blockList := getCurrentCaptureBlockListNode()
	blockList.forcedCapturing = false
}
func GetCaptureBlockListForcedCapturingFlag() bool {
	blockList := getCurrentCaptureBlockListNode()
	return blockList.forcedCapturing
}
func SetCaptureBlockListForcedCapturingFlag() {
	blockList := getCurrentCaptureBlockListNode()
	blockList.forcedCapturing = true
}

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
// Dealing with capture blocks in the current list
// <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<

func PushCaptureBlock(blockOperationName string, operandList []Node) {
	pushOntoCurrentCaptureBlockListCaptureBlockStack(blockOperationName, operandList)
	defaultFlags := getStartOperationFlags(blockOperationName)
	blockList := getCurrentCaptureBlockListNode()
	blockList.evalutesInsteadOfCapturing = defaultFlags.ForcedEval
	blockList.forcedCapturing = defaultFlags.ForcedCapture
	blockList.overwriteForcedCapture = &defaultFlags.OverwriteForcedCapture
}

func PopCaptureBlockThenExtendCapturedLines(extendedLines []CapturedLine) {
	blockStack := getCurrentCaptureBlockListNodeCaptureBlockStack()

	// More than 1 will
	if len(*blockStack) > 1 {
		popFromCurrentCaptureBlockListCaptureBlockStack()
		captureBlock := getCurrentCaptureBlockListNodeCaptureBlockStackTopEntryFurthestAlternate()
		for _, line := range extendedLines {
			captureBlock.CapturedLines = append(captureBlock.CapturedLines, line)
		}

	} else if len(*blockStack) == 1 {
		captureBlock := getCurrentCaptureBlockListNodeCaptureBlockStackTopEntryFurthestAlternate()
		captureBlock.CapturedLines = extendedLines
		GoToProcessingFlag = true

	} else {
		panic("🛑 Popping nothing/extending nothing!!!")
	}
}

func CreateNewAlternateForCaptureBlock(blockOperationName string, operandList []Node) {
	captureBlock := getCurrentCaptureBlockListNodeCaptureBlockStackTopEntryFurthestAlternate()
	altCaptureBlock := newCaptureBlock(blockOperationName, operandList)
	captureBlock.AlternateCaptureBlock = &altCaptureBlock
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++

func GetCurrentCaptureBlockStack() *[]captureBlock {
	return getCurrentCaptureBlockListNodeCaptureBlockStack()
}

func GetCurrentCaptureBlockStackLen() int {
	blockStack := getCurrentCaptureBlockListNodeCaptureBlockStack()
	return len(*blockStack)
}

func GetCurrentCaptureBlock() *captureBlock {
	return getCurrentCaptureBlockListNodeCaptureBlockStackTopEntry()
}

func GetCurrentCaptureBlockCapturedLines() *[]CapturedLine {
	captureBlock := getCurrentCaptureBlockListNodeCaptureBlockStackTopEntryFurthestAlternate()
	return &captureBlock.CapturedLines
}

func GetCurrentCaptureBlockCapturedLinesAndOperandList() (*[]CapturedLine, *[]Node) {
	captureBlock := getCurrentCaptureBlockListNodeCaptureBlockStackTopEntry()
	return &captureBlock.CapturedLines, &captureBlock.OperandList
}

func GetCurrentCaptureBlockOperationName() string {
	captureBlock := getCurrentCaptureBlockListNodeCaptureBlockStackTopEntryFurthestAlternate()
	return captureBlock.BlockOperationName
}

func SetBottomOfStackToEmptyBlock() {
	blockStack := getCurrentCaptureBlockListNodeCaptureBlockStack()
	(*blockStack)[0] = newCaptureBlock("nil", nil)
}

// ***************************************************
func GetCurrentOpPtr() *CaptureBlockListNode {
	return getCurrentCaptureBlockListNode()
}
func GetCurrentActiveOpPtr() *CaptureBlockListNode {
	return getActiveCaptureBlockListNode()
}
func GetLinesWithPtr(pointer *CaptureBlockListNode) *[]CapturedLine {
	return &pointer.captureBlockStack[0].CapturedLines
}
func GetBlockEntriesWithPtr(pointer *CaptureBlockListNode) *[]captureBlock {
	return &pointer.captureBlockStack
}
func ClearBlockEntriesWithPtr(pointer *CaptureBlockListNode) {
	pointer.captureBlockStack = (*pointer).captureBlockStack[:0]
}

// -------------------------------------------------
func SetExitOpName(opName string) {
	exitOpName = opName
}
