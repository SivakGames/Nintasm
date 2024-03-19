package romBuildingSettings

import (
	"misc/nintasm/assemble/errorHandler"
	enumErrorCodes "misc/nintasm/constants/enums/errorCodes"
)

// ------------------------------------

var autoZerpPage = true

func GetAutoZeroPage() bool {
	return autoZerpPage
}
func SetAutoZeroPage(newAutoZPValue bool) {
	autoZerpPage = newAutoZPValue
}

// ------------------------------------

var emptyRomFillValue uint8 = 0xff

func GetEmptyRomFillValue() uint8 {
	return emptyRomFillValue
}
func SetEmptyRomFillValue(newEmptyRomFillValue uint8) {
	emptyRomFillValue = newEmptyRomFillValue
}

// ------------------------------------

var rsValue uint
var rsHasBeenSetOnce bool

func GetRSValue() (uint, error) {
	if !rsHasBeenSetOnce {
		return rsValue, errorHandler.AddNew(enumErrorCodes.RsNotSet)
	}
	return rsValue, nil
}
func SetRSValue(newRSValue uint) {
	rsValue = newRSValue
	rsHasBeenSetOnce = true
}
func AddToRSValue(addRSValue uint) error {
	if !rsHasBeenSetOnce {
		return errorHandler.AddNew(enumErrorCodes.RsNotSet)
	}
	rsValue += addRSValue
	return nil
}
