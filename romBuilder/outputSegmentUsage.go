package romBuilder

import (
	"fmt"
	enumTerminalColors "misc/nintasm/constants/enums/terminalColors"
	"misc/nintasm/util"
)

func OutputSegmentUsage() {
	for rsIdx, romSegment := range rom {
		fmt.Println("Rom Segment", rsIdx)
		totalSegmentOccupied := 0
		totalSegmentMax := 0

		numBanks := len(romSegment)
		bankCaptionNumberWidth := 0
		for numBanks > 0 {
			numBanks /= 10
			bankCaptionNumberWidth++
		}

		for bIdx, bank := range romSegment {
			maxBytes := len(bank.bytes)
			occupiedBytes := int(bank.occupiedBytes)
			output := calcAndGenOutput(maxBytes, occupiedBytes)
			totalSegmentOccupied += occupiedBytes
			totalSegmentMax += maxBytes
			bankCaption := fmt.Sprintf(" Bank %v ", util.PadStringLeft(
				fmt.Sprintf("%d", bIdx),
				bankCaptionNumberWidth, ' '),
			)
			color := enumTerminalColors.AnsiGray2
			if bIdx%2 == 0 {
				color = enumTerminalColors.AnsiGray3
			}

			bankOutput := fmt.Sprintf("%v %v", util.Colorize(bankCaption, color, true), output)
			fmt.Println(bankOutput)
		}
		output := calcAndGenOutput(totalSegmentMax, totalSegmentOccupied)
		bankOutput := fmt.Sprintf("Total: %v", output)
		fmt.Println(bankOutput)
	}
}

func calcAndGenOutput(maxBytes int, occupiedBytes int) string {
	unoccupiedBytes := maxBytes - occupiedBytes
	remainPercent := 100 - (float64(unoccupiedBytes)/float64(maxBytes))*100

	out2 := commaSeparated(maxBytes)
	out1 := util.PadStringLeft(commaSeparated(occupiedBytes), len(out2), ' ')
	out3 := commaSeparated(unoccupiedBytes)

	output := fmt.Sprintf("%v%v%v (%.2f%%) - %v remaining",
		out1, util.Colorize("/", enumTerminalColors.Cyan, false), out2,
		remainPercent, out3)
	return output
}

func commaSeparated(value int) string {
	var digits []string

	if value == 0 {
		digits = append(digits, "0")
	}

	for value > 0 {
		insertValue := value % 1000
		value /= 1000
		if value > 0 {
			digits = append(digits, fmt.Sprintf(",%03d", insertValue))
		} else {
			digits = append(digits, fmt.Sprintf("%d", insertValue))
		}
	}

	finalString := ""
	for len(digits) > 0 {
		num := digits[len(digits)-1]
		finalString += num
		digits = digits[:len(digits)-1]
	}

	return finalString
}
