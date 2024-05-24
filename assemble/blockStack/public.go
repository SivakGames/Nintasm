package blockStack

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
// Dealing with capture list itself
// <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<

// When processing block ops, add another node to the list which will take precedence
func AddNewCaptureBlockListNode() {
	createAndAppendNewCaptureBlockList()
}

// When finished processing, destroy this node
func DestroyCaptureBlockListNodeWithPointer(ptr *CaptureBlockListNode) {
	destroyCaptureBlockListAtPointer(ptr)
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++
func GetCaptureBlockListForcedCapturingFlag() bool {
	blockList := getCurrentCaptureBlockListNode()
	return blockList.forcedCapturing
}

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
// Dealing with capture blocks in the current list
// <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<

func PushCaptureBlock(blockOperationName string, operandList []Node) {
	pushOntoCurrentCaptureBlockListCaptureBlockStack(blockOperationName, operandList)
	defaultFlags := getStartOperationFlags(blockOperationName)
	blockList := getCurrentCaptureBlockListNode()
	blockList.forcedCapturing = defaultFlags.ForcedCapture
	blockList.overwriteForcedCapture = &defaultFlags.OverwriteForcedCapture
}

// **************************
func NEW_PopCaptureBlockPrepProcessBlock(processedLine []ProcessLine) {
	captureBlock := getCurrentCaptureBlockListNodeCaptureBlockStackTopEntry()
	captureBlock.ProcessLines = processedLine
}

func GenerateProcessedLine(scope ProcessLineScope, finalCapturedLines []CapturedLine) ProcessLine {
	return newProcessLine(scope, finalCapturedLines)
}

// **************************

func CreateNewAlternateForCaptureBlock(blockOperationName string, operandList []Node) {
	captureBlock := getCurrentCaptureBlockListNodeCaptureBlockStackTopEntryFurthestAlternate()
	altCaptureBlock := newCaptureBlock(blockOperationName, operandList)
	captureBlock.AlternateCaptureBlock = &altCaptureBlock
}

// +++++++++++++++++++++++++++++++++++++++++++++++++++

func GetCurrentCaptureBlockStack() *[]CaptureBlock {
	return getCurrentCaptureBlockListNodeCaptureBlockStack()
}

func GetCurrentCaptureBlockStackLen() int {
	blockStack := getCurrentCaptureBlockListNodeCaptureBlockStack()
	return len(*blockStack)
}

func GetCurrentCaptureBlock() *CaptureBlock {
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

func SetCurrentCaptureBlockPostFn(postFn func()) {
	captureBlock := getCurrentCaptureBlockListNodeCaptureBlockStackTopEntry()
	captureBlock.ProcessPostFn = postFn
}

// ***************************************************

func SetInlineEval() {
	captureBlock := GetCurrentCaptureBlock()
	captureBlock.IsInlineEval = true
}

func CheckInlineEval() bool {
	captureBlock := GetCurrentCaptureBlock()
	return captureBlock.IsInlineEval
}

// ***************************************************

func CopyCapturedLinesToProcessedWithEmptyScope() {
	capturedLines := GetCurrentCaptureBlockCapturedLines()
	processedLines := []ProcessLine{}
	pl := GenerateProcessedLine(ProcessLineScope{}, *capturedLines)
	processedLines = append(processedLines, pl)
	NEW_PopCaptureBlockPrepProcessBlock(processedLines)
}

func CopyPresetCapturedLinesToProcessedWithEmptyScope(capturedLines *[]CapturedLine) {
	processedLines := []ProcessLine{}
	pl := GenerateProcessedLine(ProcessLineScope{}, *capturedLines)
	processedLines = append(processedLines, pl)
	NEW_PopCaptureBlockPrepProcessBlock(processedLines)
}

// ***************************************************
func GetCurrentOpPtr() *CaptureBlockListNode {
	return getCurrentCaptureBlockListNode()
}
func GetCurrentActiveOpPtr() *CaptureBlockListNode {
	return getActiveCaptureBlockListNode()
}
func GetBlockEntriesWithPtr(pointer *CaptureBlockListNode) *[]CaptureBlock {
	return &pointer.captureBlockStack
}
func GetCapturedLinesWithPtr(pointer *CaptureBlockListNode) *[]CapturedLine {
	return &pointer.captureBlockStack[0].CapturedLines
}
func GetCapturedLinesOpNameWithPtr(pointer *CaptureBlockListNode) string {
	return (&pointer.captureBlockStack[0]).BlockOperationName
}
func GetProcessedLinesWithPtr(pointer *CaptureBlockListNode) *[]ProcessLine {
	return &pointer.captureBlockStack[0].ProcessLines
}
func ClearBlockEntriesWithPtr(pointer *CaptureBlockListNode) {
	pointer.captureBlockStack = (*pointer).captureBlockStack[:0]
}
func GetPostFnWithPtr(pointer *CaptureBlockListNode) func() {
	return pointer.captureBlockStack[0].ProcessPostFn
}

// -------------------------------------------------
// For exit macro
func SetExitOpName(opName string) {
	exitOpName = opName
}

func GetExitOpName() string {
	return exitOpName
}
