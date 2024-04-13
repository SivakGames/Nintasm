package directiveHandler

import (
	"misc/nintasm/interpreter"
	"misc/nintasm/romBuilder/romSegmentation"
)

// +++++++++++++++++++++++++

func evalRomBuildingOperands(directiveName string, operandList *[]Node) error {
	switch directiveName {
	case "ROMSEGMENT":
		var segmentSizeNode *Node
		var segmentBankSizeNode *Node = nil
		var segmentDescriptionNode *Node = nil

		romBuildingNodes := &(*operandList)
		segmentSizeNode = &(*romBuildingNodes)[0]

		if len(*romBuildingNodes) >= 2 {
			segmentBankSizeNode = &(*romBuildingNodes)[1]
		}
		if len(*romBuildingNodes) == 3 {
			segmentDescriptionNode = &(*romBuildingNodes)[2]
		}

		interpreter.ClearParentLabel()
		interpreter.ClearLocalLabel()
		return romSegmentation.ValidateAndAddRomSegment(segmentSizeNode, segmentBankSizeNode, segmentDescriptionNode)

	case "BANK":
		bankNode := &(*operandList)[0]
		interpreter.ClearParentLabel()
		interpreter.ClearLocalLabel()
		return romSegmentation.ValidateAndSetBank(bankNode)

	case "ORG":
		orgNode := &(*operandList)[0]
		return romSegmentation.ValidateAndSetOrg(orgNode)

	default:
		panic("🛑 Something is VERY wrong with ROM building directive")
	}

}
