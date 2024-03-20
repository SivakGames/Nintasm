package blockStack

// +++++++++++++++++++++++++++++++++++++++++++++

//Traverse the list until finding the last (and current) node
func getCurrentCaptureBlockListNode() *CaptureBlockListNode {
	var listNode *CaptureBlockListNode = &captureBlockList
	for listNode.nextNode != nil {
		listNode = listNode.nextNode
	}
	return listNode
}

func getCurrentCaptureBlockListNodeCaptureBlockStack() *[]captureBlock {
	listNode := getCurrentCaptureBlockListNode()
	return &listNode.captureBlockStack
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
	highestNode := getCurrentCaptureBlockListNode()
	newNode := newCaptureBlockListNode()
	(*highestNode).nextNode = &newNode
}

func destroyCaptureBlockListAtPointer(ptr *CaptureBlockListNode) {
	var listNode *CaptureBlockListNode = &captureBlockList
	for listNode.nextNode != nil {
		if listNode.nextNode == ptr {
			listNode.nextNode = nil
			break
		}
		listNode = captureBlockList.nextNode
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
