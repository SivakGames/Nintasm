package romBuilder

import (
	"errors"
	enumSizeAliases "misc/nintasm/enums/sizeAliases"
	"misc/nintasm/parser/operandFactory"
	"misc/nintasm/util"
)

type iNESHeaderFormat struct {
	chrHeaderValue  int
	chrSizeInKb     int
	mapper          int
	mirroring       int
	prgHeaderValue  int
	prgSizeInKb     int
	hasSetChr       bool
	hasSetMapper    bool
	hasSetMirroring bool
	hasSetPrg       bool
}

var INESHeader = iNESHeaderFormat{}

const INES_PRG_SIZE_MULTIPLIER = 0x04000
const INES_CHR_SIZE_MULTIPLIER = 0x02000

var inesPrgSizeEnumAliases = map[enumSizeAliases.Def]int{
	enumSizeAliases.Size16kb:  1,
	enumSizeAliases.Size32kb:  2,
	enumSizeAliases.Size64kb:  4,
	enumSizeAliases.Size128kb: 8,
	enumSizeAliases.Size256kb: 16,
	enumSizeAliases.Size512kb: 32,
	enumSizeAliases.Size1mb:   64,
	enumSizeAliases.Size2mb:   128,
}

var inesChrSizeEnumAliases = map[enumSizeAliases.Def]int{
	enumSizeAliases.Size8kb:   1,
	enumSizeAliases.Size16kb:  2,
	enumSizeAliases.Size32kb:  4,
	enumSizeAliases.Size64kb:  8,
	enumSizeAliases.Size128kb: 16,
	enumSizeAliases.Size256kb: 32,
	enumSizeAliases.Size512kb: 64,
	enumSizeAliases.Size1mb:   128,
	enumSizeAliases.Size2mb:   256,
}

// -----------------------------------------

// INES Mapper declaration
func ValidateInesMap(inesNode *Node) error {
	if INESHeader.hasSetMapper {
		return errors.New("INES mapper has already been set!!!")
	}

	if !(operandFactory.ValidateNodeIsNumeric(inesNode) &&
		operandFactory.ValidateNumericNodeIsPositive(inesNode) &&
		operandFactory.ValidateNumericNodeIs8BitValue(inesNode)) {
		return errors.New("INES mapper must be a non-negative 8-bit number")
	}

	INESHeader.hasSetMapper = true
	INESHeader.mapper = inesNode.AsNumber
	return nil
}

// -----------------------------------------

// INES Mirroring declaration
func ValidateInesMirroring(inesNode *Node) error {
	if INESHeader.hasSetMirroring {
		return errors.New("INES mirroring has already been set!!!")
	}
	if !(operandFactory.ValidateNodeIsNumeric(inesNode) &&
		operandFactory.ValidateNumericNodeIsGTEandLTEValues(inesNode, 0, 1)) {
		return errors.New("INES mirroring must be either 0 or 1")
	}

	INESHeader.hasSetMirroring = true
	INESHeader.mirroring = inesNode.AsNumber
	return nil
}

// -----------------------------------------

// INES PRG declaration
func ValidateInesPrg(inesNode *Node) error {
	inesOperationDescription := "INES PRG"

	if INESHeader.hasSetPrg {
		return errors.New("INES PRG has already been set!!!")
	}

	if operandFactory.ValidateNodeIsString(inesNode) {
		err := util.ValidateSizeStringAliasUsable(inesNode, &inesPrgSizeEnumAliases, inesOperationDescription)
		if err != nil {
			return err
		}
	}

	if !(operandFactory.ValidateNodeIsNumeric(inesNode) &&
		operandFactory.ValidateNumericNodeIsPositive(inesNode)) {
		return errors.New("INES PRG must be >= 1 or use a valid alias")
	}

	err := util.ValidateSizeNumberAliasUsable(inesNode, &inesPrgSizeEnumAliases, inesOperationDescription)
	if err != nil {
		return err
	}

	if !util.ValidateIsPowerOfTwo(inesNode.AsNumber) ||
		operandFactory.ValidateNumericNodeIsGTValue(inesNode, 128) {
		return errors.New("Unacceptable INES PRG size declared!")
	}

	INESHeader.hasSetPrg = true
	INESHeader.prgHeaderValue = inesNode.AsNumber
	INESHeader.prgSizeInKb = INESHeader.prgHeaderValue * INES_PRG_SIZE_MULTIPLIER
	return nil
}

// -----------------------------------------

// INES CHR declaration
func ValidateInesChr(inesNode *Node) error {
	inesOperationDescription := "INES CHR"

	if INESHeader.hasSetChr {
		return errors.New("INES CHR has already been set!!!")
	}

	if operandFactory.ValidateNodeIsString(inesNode) {
		err := util.ValidateSizeStringAliasUsable(inesNode, &inesChrSizeEnumAliases, inesOperationDescription)
		if err != nil {
			return err
		}
	}

	if !(operandFactory.ValidateNodeIsNumeric(inesNode) &&
		operandFactory.ValidateNumericNodeIsPositive(inesNode)) {
		return errors.New("INES CHR must be >= 1 or use a valid alias")
	}

	err := util.ValidateSizeNumberAliasUsable(inesNode, &inesChrSizeEnumAliases, inesOperationDescription)
	if err != nil {
		return err
	}

	if !util.ValidateIsPowerOfTwo(inesNode.AsNumber) ||
		operandFactory.ValidateNumericNodeIsGTValue(inesNode, 256) {
		return errors.New("Unacceptable INES CHR size declared!")
	}

	INESHeader.hasSetChr = true
	INESHeader.chrHeaderValue = inesNode.AsNumber
	INESHeader.chrSizeInKb = INESHeader.chrHeaderValue * INES_CHR_SIZE_MULTIPLIER
	return nil
}

// +++++++++++++++++++++++++++++++++++++++++++++

func GetInesMap() int {
	return INESHeader.mapper
}
func GetInesPrgHeaderValue() int {
	return INESHeader.prgHeaderValue
}
func GetInesPrgSizeInKb() int {
	return INESHeader.prgSizeInKb
}

func GetInesChrHeaderValue() int {
	return INESHeader.chrHeaderValue
}
func GetInesChrSizeInKb() int {
	return INESHeader.chrSizeInKb
}

func GetInesTotalRomSizeInKb() int {
	return GetInesChrSizeInKb() + GetInesPrgSizeInKb()
}

func GetInesMirroring() int {
	return INESHeader.mirroring
}
