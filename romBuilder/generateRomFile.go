package romBuilder

import (
	"misc/nintasm/assemble/fileHandler"
	"misc/nintasm/util/commandLine"
	"os"
)

func GenerateRomBinFile() (string, error) {
	outFileName := fileHandler.GenerateOutFileName()
	file, err := os.OpenFile(outFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return "", err
	}
	defer file.Close()

	isRawFile := commandLine.SettingRawROMFile

	if !isRawFile {
		iNESHeader := GenerateINESHeader()

		_, err = file.Write(iNESHeader)
		if err != nil {
			return "", err
		}
	}

	for _, romSegment := range rom {
		for _, bank := range romSegment {
			_, err = file.Write(bank.bytes)
			if err != nil {
				return "", err
			}
		}
	}

	return outFileName, nil
}
