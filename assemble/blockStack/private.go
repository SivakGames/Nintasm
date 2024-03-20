package blockStack

// +++++++++++++++++++++++++++++++++++++++++++++

//Traverse the list until finding the last (and current) node
func getCurrentCaptureBlockListNode() *CaptureBlockListNode {
	var list *CaptureBlockListNode = &captureBlockList
	for list.nextNode != nil {
		list = list.nextNode
	}
	return list
}

func getCurrentCaptureBlockListNodeCaptureBlockStack() *[]captureBlock {
	list := getCurrentCaptureBlockListNode()
	return &list.captureBlockStack
}

func getCurrentCaptureBlockListCaptureBlockStackTop() *captureBlock {
	blockStack := getCurrentCaptureBlockListNodeCaptureBlockStack()
	return &(*blockStack)[len(*blockStack)-1]
}

func getCurrentCaptureBlockListNodeCaptureBlockStackTopFurthestAlternate() *captureBlock {
	block := getCurrentCaptureBlockListCaptureBlockStackTop()
	for block.AlternateCaptureBlock != nil {
		block = block.AlternateCaptureBlock
	}
	return block
}

// >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
// <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
func createAndAppendNewCaptureBlockList() {
	highestOp := getCurrentCaptureBlockListNode()
	newOp := newCaptureBlockListNode()
	(*highestOp).nextNode = &newOp
}

func destroyCaptureBlockListAtPointer(ptr *CaptureBlockListNode) {
	var invokeOp *CaptureBlockListNode = &captureBlockList
	for invokeOp.nextNode != nil {
		if invokeOp.nextNode == ptr {
			invokeOp.nextNode = nil
			break
		}
		invokeOp = captureBlockList.nextNode
	}
}

func pushOntoCurrentCaptureBlockListCaptureBlockStack(blockOperationName string, operandList []Node) {
	blockEntries := getCurrentCaptureBlockListNodeCaptureBlockStack()
	*blockEntries = append(*blockEntries, newCaptureBlock(blockOperationName, operandList))
}
func popFromCurrentCaptureBlockListCaptureBlockStack() {
	blockEntries := getCurrentCaptureBlockListNodeCaptureBlockStack()
	*blockEntries = (*blockEntries)[:len(*blockEntries)-1]
}
