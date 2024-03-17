package blockStack

// +++++++++++++++++++++++++++++++++++++++++++++
var invokeOperations CaptureBlockList = newCaptureBlockList()

func getCurrentCaptureBlockList() *CaptureBlockList {
	var list *CaptureBlockList = &invokeOperations
	for list.nextList != nil {
		list = list.nextList
	}
	return list
}

func getCurrentCaptureBlockListCaptureBlockStack() *[]captureBlock {
	list := getCurrentCaptureBlockList()
	return &list.captureBlockStack
}

func getCurrentCaptureBlockListCaptureBlockStackTop() *captureBlock {
	blockStack := getCurrentCaptureBlockListCaptureBlockStack()
	return &(*blockStack)[len(*blockStack)-1]
}

func getCurrentCaptureBlockListCaptureBlockStackTopFurthestAlternate() *captureBlock {
	block := getCurrentCaptureBlockListCaptureBlockStackTop()
	for block.AlternateCaptureBlock != nil {
		block = block.AlternateCaptureBlock
	}
	return block
}

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
// <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
func createAndAppendNewCaptureBlockList() {
	highestOp := getCurrentCaptureBlockList()
	newOp := newCaptureBlockList()
	(*highestOp).nextList = &newOp
}

func destroyCaptureBlockListAtPointer(ptr *CaptureBlockList) {
	var invokeOp *CaptureBlockList = &invokeOperations
	for invokeOp.nextList != nil {
		if invokeOp.nextList == ptr {
			invokeOp.nextList = nil
			break
		}
		invokeOp = invokeOperations.nextList
	}
}

func pushOntoCurrentCaptureBlockListCaptureBlockStack(blockOperationName string, operandList []Node) {
	blockEntries := getCurrentCaptureBlockListCaptureBlockStack()
	*blockEntries = append(*blockEntries, newCaptureBlock(blockOperationName, operandList))
}
func popFromCurrentCaptureBlockListCaptureBlockStack() {
	blockEntries := getCurrentCaptureBlockListCaptureBlockStack()
	*blockEntries = (*blockEntries)[:len(*blockEntries)-1]
}
