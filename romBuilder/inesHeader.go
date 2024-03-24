package romBuilder

import (
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
	enumSizeAliases "misc/nintasm/constants/enums/sizeAliases"
	"misc/nintasm/interpreter/operandFactory"
	"misc/nintasm/util"
	"misc/nintasm/util/validateSizeAlias"
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
		return errorHandler.AddNew(enumErrorCodes.INESValueAlreadySet, "INES mapper")
	}
	if !operandFactory.ValidateNodeIsNumeric(inesNode) {
		return errorHandler.AddNew(enumErrorCodes.NodeTypeNotNumeric)
	} else if !operandFactory.ValidateNumericNodeIsPositive(inesNode) {
		return errorHandler.AddNew(enumErrorCodes.NodeValueNotPositive)
	} else if !(operandFactory.ValidateNumericNodeIs8BitValue(inesNode)) {
		return errorHandler.AddNew(enumErrorCodes.NodeValueNot8Bit)
	}

	INESHeader.hasSetMapper = true
	INESHeader.mapper = int(inesNode.AsNumber)
	return nil
}

// -----------------------------------------

// INES Mirroring declaration
func ValidateInesMirroring(inesNode *Node) error {
	if INESHeader.hasSetMirroring {
		return errorHandler.AddNew(enumErrorCodes.INESValueAlreadySet, "INES mirroring")
	} else if !operandFactory.ValidateNodeIsNumeric(inesNode) {
		return errorHandler.AddNew(enumErrorCodes.NodeTypeNotNumeric)
	} else if !operandFactory.ValidateNumericNodeIsGTEandLTEValues(inesNode, 0, 1) {
		return errorHandler.AddNew(enumErrorCodes.NodeValueNotGTEandLTE, 0, 1)
	}

	INESHeader.hasSetMirroring = true
	INESHeader.mirroring = int(inesNode.AsNumber)
	return nil
}

// -----------------------------------------

// INES PRG declaration
func ValidateInesPrg(inesNode *Node) error {
	inesOperationDescription := "INES PRG"

	if INESHeader.hasSetPrg {
		return errorHandler.AddNew(enumErrorCodes.INESValueAlreadySet, inesOperationDescription)
	}

	if operandFactory.ValidateNodeIsString(inesNode) {
		err := validateSizeAlias.ValidateSizeStringAliasUsable(inesNode, &inesPrgSizeEnumAliases, inesOperationDescription)
		if err != nil {
			return err
		}
	}

	if !operandFactory.ValidateNodeIsNumeric(inesNode) {
		return errorHandler.AddNew(enumErrorCodes.NodeTypeNotNumeric)
	} else if !operandFactory.ValidateNumericNodeIsPositive(inesNode) {
		return errorHandler.AddNew(enumErrorCodes.NodeValueNotPositive)
	}

	err := validateSizeAlias.ValidateSizeNumberAliasUsable(inesNode, &inesPrgSizeEnumAliases, inesOperationDescription)
	if err != nil {
		return err
	}

	if !util.ValidateIsPowerOfTwo(int(inesNode.AsNumber)) {
		return errorHandler.AddNew(enumErrorCodes.NodeValueNotPowerOf2)
	} else if operandFactory.ValidateNumericNodeIsGTValue(inesNode, 128) {
		return errorHandler.AddNew(enumErrorCodes.NodeValueNotLTE, 128)
	}

	INESHeader.hasSetPrg = true
	INESHeader.prgHeaderValue = int(inesNode.AsNumber)
	INESHeader.prgSizeInKb = INESHeader.prgHeaderValue * INES_PRG_SIZE_MULTIPLIER
	return nil
}

// -----------------------------------------

// INES CHR declaration
func ValidateInesChr(inesNode *Node) error {
	inesOperationDescription := "INES CHR"

	if INESHeader.hasSetChr {
		return errorHandler.AddNew(enumErrorCodes.INESValueAlreadySet, inesOperationDescription)
	}

	if operandFactory.ValidateNodeIsString(inesNode) {
		err := validateSizeAlias.ValidateSizeStringAliasUsable(inesNode, &inesChrSizeEnumAliases, inesOperationDescription)
		if err != nil {
			return err
		}
	}

	if !operandFactory.ValidateNodeIsNumeric(inesNode) {
		return errorHandler.AddNew(enumErrorCodes.NodeTypeNotNumeric)
	} else if !operandFactory.ValidateNumericNodeIsPositive(inesNode) {
		return errorHandler.AddNew(enumErrorCodes.NodeValueNotPositive)
	}

	err := validateSizeAlias.ValidateSizeNumberAliasUsable(inesNode, &inesChrSizeEnumAliases, inesOperationDescription)
	if err != nil {
		return err
	}

	if !util.ValidateIsPowerOfTwo(int(inesNode.AsNumber)) {
		return errorHandler.AddNew(enumErrorCodes.NodeValueNotPowerOf2)
	} else if operandFactory.ValidateNumericNodeIsGTValue(inesNode, 256) {
		return errorHandler.AddNew(enumErrorCodes.NodeValueNotLTE, 256)
	}

	INESHeader.hasSetChr = true
	INESHeader.chrHeaderValue = int(inesNode.AsNumber)
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

//xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

func ClearINES() {
	INESHeader.hasSetChr = false
	INESHeader.hasSetPrg = false
	INESHeader.hasSetMirroring = false
	INESHeader.hasSetMapper = false
}

//=================================================

func GenerateINESHeader() []uint8 {
	iNESHeader := make([]uint8, 16)
	iNESHeader[0] = 'N'
	iNESHeader[1] = 'E'
	iNESHeader[2] = 'S'
	iNESHeader[3] = 0x1a

	mirrorFlag := GetInesMirroring()
	batteryFlag := 0 //TODO
	mapperHighNibble := GetInesMap() & 0x00f0
	mapperLowNibble := GetInesMap() & 0x000f

	flags6 := (mapperLowNibble << 4) | mirrorFlag | batteryFlag
	flags7 := (mapperHighNibble)

	iNESHeader[4] = uint8(GetInesPrgHeaderValue())
	iNESHeader[5] = uint8(GetInesChrHeaderValue())
	iNESHeader[6] = uint8(flags6)
	iNESHeader[7] = uint8(flags7)
	return iNESHeader
}
