package blockStack

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
// Dealing with capture list itself
// <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<

// When processing block ops (macros in particular), add another capture list
func AddNewCaptureBlockList() {
	createAndAppendNewCaptureBlockList()
}

func DestroyCaptureBlockListWithPointer(ptr *CaptureBlockList) {
	destroyCaptureBlockListAtPointer(ptr)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++
func ClearCaptureBlockListEvalFlag() {
	blockList := getCurrentCaptureBlockList()
	blockList.evalutesInsteadOfCapturing = false
}
func GetCaptureBlockListEvalFlag() bool {
	blockList := getCurrentCaptureBlockList()
	return blockList.evalutesInsteadOfCapturing
}
func SetCaptureBlockListEvalFlag() {
	blockList := getCurrentCaptureBlockList()
	blockList.evalutesInsteadOfCapturing = true
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++
func ClearCaptureBlockListForcedCapturingFlag() {
	blockList := getCurrentCaptureBlockList()
	blockList.forcedCapturing = false
}
func GetCaptureBlockListForcedCapturingFlag() bool {
	blockList := getCurrentCaptureBlockList()
	return blockList.forcedCapturing
}
func SetCaptureBlockListForcedCapturingFlag() {
	blockList := getCurrentCaptureBlockList()
	blockList.forcedCapturing = true
}

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
// Dealing with capture blocks in the current list
// <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<

func PushCaptureBlock(blockOperationName string, operandList []Node) {
	pushOntoCurrentCaptureBlockListCaptureBlockStack(blockOperationName, operandList)
	defaultFlags := getStartOperationFlags(blockOperationName)
	blockList := getCurrentCaptureBlockList()
	blockList.evalutesInsteadOfCapturing = defaultFlags.ForcedEval
	blockList.forcedCapturing = defaultFlags.ForcedCapture
}

func ForcePopCaptureBlock() {
	popFromCurrentCaptureBlockListCaptureBlockStack()
}

func PopCaptureBlockThenExtendCapturedLines(extendedLines []CapturedLine) {
	blockStack := getCurrentCaptureBlockListCaptureBlockStack()

	// More than 1 will
	if len(*blockStack) > 1 {
		popFromCurrentCaptureBlockListCaptureBlockStack()
		captureBlock := getCurrentCaptureBlockListCaptureBlockStackTopFurthestAlternate()
		for _, line := range extendedLines {
			captureBlock.CapturedLines = append(captureBlock.CapturedLines, line)
		}

	} else if len(*blockStack) == 1 {
		//Set eval operands to true
		captureBlock := getCurrentCaptureBlockListCaptureBlockStackTopFurthestAlternate()
		captureBlock.CapturedLines = extendedLines
		GoToProcessingFlag = true

	} else {
		panic("🛑 Popping nothing/extending nothing!!!")
	}
}

func CreateNewAlternateForCaptureBlock(blockOperationName string, operandList []Node) {
	captureBlock := getCurrentCaptureBlockListCaptureBlockStackTopFurthestAlternate()
	altCaptureBlock := newCaptureBlock(blockOperationName, operandList)
	captureBlock.AlternateCaptureBlock = &altCaptureBlock
}

//================================================

func GetCurrentCaptureBlockStack() *[]captureBlock {
	return getCurrentCaptureBlockListCaptureBlockStack()
}

func GetCurrentCaptureBlock() *captureBlock {
	captureBlock := getCurrentCaptureBlockListCaptureBlockStackTop()
	return captureBlock
}

func GetCurrentCaptureBlockCapturedLines() *[]CapturedLine {
	captureBlock := getCurrentCaptureBlockListCaptureBlockStackTopFurthestAlternate()
	return &captureBlock.CapturedLines
}

func GetCurrentCaptureBlockCapturedLinesAndOperandList() (*[]CapturedLine, *[]Node) {
	captureBlock := getCurrentCaptureBlockListCaptureBlockStackTop()
	return &captureBlock.CapturedLines, &captureBlock.OperandList
}

func GetCurrentCaptureBlockOperationName() string {
	captureBlock := getCurrentCaptureBlockListCaptureBlockStackTopFurthestAlternate()
	return captureBlock.BlockOperationName
}

func SetBottomOfStackToEmptyBlock() {
	blockStack := getCurrentCaptureBlockListCaptureBlockStack()
	(*blockStack)[0] = newCaptureBlock("nil", nil)
}

// ***************************************************
func GetCurrentOpPtr() *CaptureBlockList {
	return getCurrentCaptureBlockList()
}
func GetLinesWithPtr(pointer *CaptureBlockList) *[]CapturedLine {
	return &pointer.captureBlockStack[0].CapturedLines
}
func GetBlockEntriesWithPtr(pointer *CaptureBlockList) *[]captureBlock {
	return &pointer.captureBlockStack
}
func ClearBlockEntriesWithPtr(pointer *CaptureBlockList) {
	pointer.captureBlockStack = (*pointer).captureBlockStack[:0]
}
