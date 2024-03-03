package romBuilder

import (
	"misc/nintasm/assemble/fileHandler"
	"os"
)

func GenerateRomBinFile() error {
	outFileName := fileHandler.AddRelativePathIncludeFileName("output.nes")
	file, err := os.OpenFile(outFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	iNESHeader := GenerateINESHeader()

	_, err = file.Write(iNESHeader)
	if err != nil {
		return err
	}

	for _, romSegment := range rom {
		for _, bank := range romSegment {
			_, err = file.Write(bank.bytes)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
