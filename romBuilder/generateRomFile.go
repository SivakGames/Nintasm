package romBuilder

import (
	"misc/nintasm/assemble/fileHandler"
	"os"
)

func GenerateRomBinFile() error {
	outFileName := fileHandler.AddRelativePathIncludeFileName("output.nes")
	file, err := os.OpenFile(outFileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, r := range rom {
		for _, b := range r {
			_, err = file.Write(b.bytes)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
