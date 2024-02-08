package romSegmentation

import (
	"errors"
	enumSizeAliases "misc/nintasm/enums/sizeAliases"
	"misc/nintasm/parser/operandFactory"
	"misc/nintasm/romBuilder"
	"misc/nintasm/util"
)

type Node = operandFactory.Node

const ROM_SEGMENT_MIN_SIZE = 0x0000400
const ROM_SEGMENT_MAX_SIZE = 0x0200000

var romSegmentEnumAliases = map[enumSizeAliases.Def]int{
	enumSizeAliases.Size1kb:   ROM_SEGMENT_MIN_SIZE,
	enumSizeAliases.Size2kb:   0x0000800,
	enumSizeAliases.Size4kb:   0x0001000,
	enumSizeAliases.Size8kb:   0x0002000,
	enumSizeAliases.Size16kb:  0x0004000,
	enumSizeAliases.Size32kb:  0x0008000,
	enumSizeAliases.Size64kb:  0x0010000,
	enumSizeAliases.Size128kb: 0x0020000,
	enumSizeAliases.Size256kb: 0x0040000,
	enumSizeAliases.Size512kb: 0x0080000,
	enumSizeAliases.Size1mb:   0x0100000,
	enumSizeAliases.Size2mb:   ROM_SEGMENT_MAX_SIZE,
}

//-------------------------------------------

func ValidateAndAddRomSegment(segmentSizeNode *Node, segmentBankSizeNode *Node, segmentDescriptionNode *Node) error {
	var segmentSize int
	var segmentBankSize int
	var segmentDescription string
	var err error

	_ = segmentDescription

	// Check segment size

	if operandFactory.ValidateNodeIsString(segmentSizeNode) {
		err = util.ValidateSizeStringAliasUsable(segmentSizeNode, &romSegmentEnumAliases, "ROM SEGMENT - segment size")
		if err != nil {
			return err
		}
	}
	err = validateNodeNumerics(segmentSizeNode)
	if err != nil {
		return err
	}
	segmentSize = segmentSizeNode.AsNumber

	// Check segment's bank size

	if segmentBankSizeNode == nil {
		segmentBankSize = segmentSize
	} else {
		if operandFactory.ValidateNodeIsString(segmentBankSizeNode) {
			err := util.ValidateSizeStringAliasUsable(segmentBankSizeNode, &romSegmentEnumAliases, "ROM SEGMENT - bank size")
			if err != nil {
				return err
			}
		}
		err = validateNodeNumerics(segmentBankSizeNode)
		if err != nil {
			return err
		}
		segmentBankSize = segmentBankSizeNode.AsNumber
	}

	// Extend ROM with new segment

	err = romBuilder.AddNewRomSegment(segmentSize, segmentBankSize)
	return err
}

func validateNodeNumerics(node *Node) error {
	if !operandFactory.ValidateNodeIsNumeric(node) {
		return errors.New("ROM Segment size must be a number")
	}
	if !operandFactory.ValidateNumericNodeIsGTEandLTEValues(node, ROM_SEGMENT_MIN_SIZE, ROM_SEGMENT_MAX_SIZE) {
		return errors.New("Unacceptable ROM Segment size!")
	}
	if !util.ValidateIsPowerOfTwo(node.AsNumber) {
		return errors.New("ROM Segment size must be a power of 2")
	}
	return nil
}
