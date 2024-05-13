package blockStack

//Main block capture linked list
var captureBlockList CaptureBlockListNode
var exitOpName string = ""

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func init() {
	captureBlockList = newCaptureBlockListNode()
}

// +++++++++++++++++++++++++++++++++++++++++++++

//Traverse the list until finding the last (and current) node
func getCurrentCaptureBlockListNode() *CaptureBlockListNode {
	var listNode *CaptureBlockListNode = &captureBlockList
	for listNode.nextNode != nil {
		listNode = listNode.nextNode
	}
	return listNode
}

//Traverse the list until finding the current active Node
func getActiveCaptureBlockListNode() *CaptureBlockListNode {
	var listNode *CaptureBlockListNode = &captureBlockList
	prevNode := listNode
	for listNode.nextNode != nil {
		prevNode = listNode
		listNode = listNode.nextNode
	}
	return prevNode
}

//Get the entire capture block stack for the current node
func getCurrentCaptureBlockListNodeCaptureBlockStack() *[]captureBlock {
	listNode := getCurrentCaptureBlockListNode()
	return &listNode.captureBlockStack
}

func getCurrentCaptureBlockListNodeCaptureBlockStackTopEntry() *captureBlock {
	blockStack := getCurrentCaptureBlockListNodeCaptureBlockStack()
	if len(*blockStack) > 0 {
		return &(*blockStack)[len(*blockStack)-1]
	}
	panic("🛑 Block stack has no length!")
}

func getCurrentCaptureBlockListNodeCaptureBlockStackTopEntryFurthestAlternate() *captureBlock {
	block := getCurrentCaptureBlockListNodeCaptureBlockStackTopEntry()
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
		listNode = listNode.nextNode
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
