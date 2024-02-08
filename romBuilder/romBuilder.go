package romBuilder

import (
	"errors"
	"math"
	"misc/nintasm/parser/operandFactory"
)

type Node = operandFactory.Node

type bankDefStruct struct {
	bytes    []uint8
	orgIsSet bool
	maxOrg   int
	minOrg   int
}

func newBankDef(bankSize int) bankDefStruct {
	return bankDefStruct{bytes: make([]uint8, bankSize), orgIsSet: false, minOrg: -1, maxOrg: -1}
}

// The final ROM that will be built
var romLayout = make([][]bankDefStruct, 0)

// Corresponding ORG info for
//var romOrgInfos = make([][]bankDefStruct, 0)

//var allocatedRomSize = 0

var currentRomSegmentIndex = -1
var currentBankIndex = -1
var CurrentInsertionIndex = -1
var currentOrg = 0x8000

//-------------------------------------------

func AddNewRomSegment(totalSize int, bankSize int) error {
	floatResult := float64(totalSize) / float64(bankSize)
	numBanks, deci := math.Modf(floatResult)
	if deci != 0 {
		return errors.New("Bank size is not evenly distributable")
	}

	newSegment := make([]bankDefStruct, int(numBanks))
	//newOrgDefs := make([]bankDefStruct, int(numBanks))

	for i := range newSegment {
		//newSegment[i] = make([]uint8, bankSize)
		//	newOrgDefs[i] = newBankDef(bankSize)
		newSegment[i] = newBankDef(bankSize)
	}

	romLayout = append(romLayout, newSegment)

	currentRomSegmentIndex = len(romLayout) - 1
	currentBankIndex = -1
	CurrentInsertionIndex = -1
	return nil
}

// The entire ROM layout (array of segments)
func GetRomLayout() *[][]bankDefStruct {
	return &romLayout
}

// The current ROM segment (array of bank segments)
func GetCurrentRomSegment() *[]bankDefStruct {

	return &romLayout[currentRomSegmentIndex]
}

func GetTotalRomSegmentsInRom() int {
	return len(*GetRomLayout())
}

// The current bank segment (array of uint8 )
func GetCurrentBankSegment() *bankDefStruct {
	currentRomSegment := GetCurrentRomSegment()
	return &(*currentRomSegment)[currentBankIndex]
}

// The current bank segment (array of uint8 )
func GetCurrentBankSegmentBytes() *[]uint8 {
	currentRomSegment := GetCurrentBankSegment()
	return &(*currentRomSegment).bytes
}

func GetTotalBanksInCurrentRomSegment() int {
	return len(*GetCurrentRomSegment())
}

func GetBankIndex() int {
	return currentBankIndex
}
func SetBankIndex(newBankIndex int) error {
	currentBankIndex = newBankIndex
	return nil
}

func GetOrg() int {
	return currentOrg + CurrentInsertionIndex
}

// TODO: Set ORG upper/lower bounds checks
func SetOrg(newOrg int) {
	currentOrg = newOrg
	return
}
