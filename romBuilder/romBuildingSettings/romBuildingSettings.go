package romBuildingSettings

import "errors"

// ------------------------------------

var autoZerpPage = true

func GetAutoZeroPage() bool {
	return autoZerpPage
}
func SetAutoZeroPage(newAutoZPValue bool) {
	autoZerpPage = newAutoZPValue
}

// ------------------------------------

var emptyRomFillValue uint8 = 0x00

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
		return rsValue, errors.New("RS has not been set!")
	}
	return rsValue, nil
}
func SetRSValue(newRSValue uint) {
	rsValue = newRSValue
	rsHasBeenSetOnce = true
}
func AddToRSValue(addRSValue uint) error {
	if !rsHasBeenSetOnce {
		return errors.New("RS has not yet been set!")
	}
	rsValue += addRSValue
	return nil
}
