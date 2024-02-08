package romBuilder

import (
	"errors"
	"fmt"
	"math"
	"misc/nintasm/parser/operandFactory"
)

type Node = operandFactory.Node

type romType []romSegmentType
type romSegmentType []bankType
type bankType struct {
	bytes    []uint8
	orgIsSet bool
	maxOrg   int
	minOrg   int
}

func newBank(bankSize int) bankType {
	return bankType{bytes: make([]uint8, bankSize), orgIsSet: false, minOrg: -1, maxOrg: -1}
}

// The final ROM that will be built
var rom = make(romType, 0)

var currentBankIndex = -1
var CurrentInsertionIndex = -1

//-------------------------------------------

func AddNewRomSegment(totalSize int, bankSize int) error {
	floatResult := float64(totalSize) / float64(bankSize)
	numBanks, deci := math.Modf(floatResult)
	if deci != 0 {
		return errors.New("Bank size is not evenly distributable")
	}

	newSegment := make([]bankType, int(numBanks))
	//newOrgDefs := make([]bankDefStruct, int(numBanks))

	for i := range newSegment {
		//newSegment[i] = make([]uint8, bankSize)
		//	newOrgDefs[i] = newBankDef(bankSize)
		newSegment[i] = newBank(bankSize)
	}

	rom = append(rom, newSegment)

	currentBankIndex = -1
	CurrentInsertionIndex = -1
	return nil
}

// The entire ROM layout (array of segments)
func getRom() *romType {
	return &rom
}

//+++++++++++++++++++++++++++

// The current ROM segment (array of bank segments)
func GetCurrentRomSegment() *romSegmentType {
	if len(rom) > 0 {
		return &rom[len(rom)-1]
	}
	panic("Attemped to access ROM with no segments!!!")
}

// How many rom segments are currently defined
func GetTotalRomSegmentsInRom() int {
	return len(*getRom())
}

//+++++++++++++++++++++++++++

// The current bank segment (array of uint8 )
func GetCurrentBankSegment() *bankType {
	currentRomSegment := GetCurrentRomSegment()
	return &(*currentRomSegment)[currentBankIndex]
}

// The current bank segment (array of uint8 )
func GetCurrentBankSegmentBytes() *[]uint8 {
	currentRomSegment := GetCurrentBankSegment()
	return &(*currentRomSegment).bytes
}

// How many banks are in the current ROM segment
func GetTotalBanksInCurrentRomSegment() int {
	return len(*GetCurrentRomSegment())
}

//+++++++++++++++++++++++++++

func GetBankIndex() int {
	return currentBankIndex
}
func SetBankIndex(newBankIndex int) error {
	currentBankIndex = newBankIndex
	return nil
}

//+++++++++++++++++++++++++++

func GetOrg() int {
	bank := GetCurrentBankSegment()
	return bank.minOrg + CurrentInsertionIndex
}

// TODO: Set ORG upper/lower bounds checks
func SetOrg(newOrg int) error {
	bank := GetCurrentBankSegment()
	bankSize := len(bank.bytes)

	if !bank.orgIsSet {
		newMinOrg := int((newOrg / bankSize) * bankSize)
		newMaxOrg := newMinOrg + bankSize - 1
		bank.minOrg = newMinOrg
		bank.maxOrg = newMaxOrg
		bank.orgIsSet = true
	} else {
		currentOrg := bank.minOrg + CurrentInsertionIndex
		if newOrg < bank.minOrg {
			errMsg := fmt.Sprintf("ORG is too small! Attempted: %d / Minimum Allowed: %d ", newOrg, bank.minOrg)
			return errors.New(errMsg)
		} else if newOrg > bank.maxOrg {
			errMsg := fmt.Sprintf("ORG is too big! Attempted: %d / Max Allowed: %d ", newOrg, bank.maxOrg)
			return errors.New(errMsg)
		} else if newOrg <= currentOrg {
			errMsg := fmt.Sprintf("Cannot set ORG to a value less than where the program counter currently is!\nThis would overwrite data!\n Attempted: %d / Currently at: %d ", newOrg, currentOrg)
			return errors.New(errMsg)
		}
	}
	CurrentInsertionIndex = newOrg % bankSize
	return nil
}
