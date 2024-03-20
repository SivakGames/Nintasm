package blockStack

//Main block capture linked list
var captureBlockList CaptureBlockListNode

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func init() {
	captureBlockList = newCaptureBlockListNode()
}

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
}

func ForcePopCaptureBlock() {
	popFromCurrentCaptureBlockListCaptureBlockStack()
}

func PopCaptureBlockThenExtendCapturedLines(extendedLines []CapturedLine) {
	blockStack := getCurrentCaptureBlockListNodeCaptureBlockStack()

	// More than 1 will
	if len(*blockStack) > 1 {
		popFromCurrentCaptureBlockListCaptureBlockStack()
		captureBlock := getCurrentCaptureBlockListNodeCaptureBlockStackTopFurthestAlternate()
		for _, line := range extendedLines {
			captureBlock.CapturedLines = append(captureBlock.CapturedLines, line)
		}

	} else if len(*blockStack) == 1 {
		//Set eval operands to true
		captureBlock := getCurrentCaptureBlockListNodeCaptureBlockStackTopFurthestAlternate()
		captureBlock.CapturedLines = extendedLines
		GoToProcessingFlag = true

	} else {
		panic("🛑 Popping nothing/extending nothing!!!")
	}
}

func CreateNewAlternateForCaptureBlock(blockOperationName string, operandList []Node) {
	captureBlock := getCurrentCaptureBlockListNodeCaptureBlockStackTopFurthestAlternate()
	altCaptureBlock := newCaptureBlock(blockOperationName, operandList)
	captureBlock.AlternateCaptureBlock = &altCaptureBlock
}

//================================================

func GetCurrentCaptureBlockStack() *[]captureBlock {
	return getCurrentCaptureBlockListNodeCaptureBlockStack()
}

func GetCurrentCaptureBlock() *captureBlock {
	captureBlock := getCurrentCaptureBlockListCaptureBlockStackTop()
	return captureBlock
}

func GetCurrentCaptureBlockCapturedLines() *[]CapturedLine {
	captureBlock := getCurrentCaptureBlockListNodeCaptureBlockStackTopFurthestAlternate()
	return &captureBlock.CapturedLines
}

func GetCurrentCaptureBlockCapturedLinesAndOperandList() (*[]CapturedLine, *[]Node) {
	captureBlock := getCurrentCaptureBlockListCaptureBlockStackTop()
	return &captureBlock.CapturedLines, &captureBlock.OperandList
}

func GetCurrentCaptureBlockOperationName() string {
	captureBlock := getCurrentCaptureBlockListNodeCaptureBlockStackTopFurthestAlternate()
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
func GetLinesWithPtr(pointer *CaptureBlockListNode) *[]CapturedLine {
	return &pointer.captureBlockStack[0].CapturedLines
}
func GetBlockEntriesWithPtr(pointer *CaptureBlockListNode) *[]captureBlock {
	return &pointer.captureBlockStack
}
func ClearBlockEntriesWithPtr(pointer *CaptureBlockListNode) {
	pointer.captureBlockStack = (*pointer).captureBlockStack[:0]
}
