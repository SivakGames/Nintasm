package romBuilder

import (
	"fmt"
	enumTerminalColors "misc/nintasm/constants/enums/terminalColors"
	"misc/nintasm/util"
	"strings"
)

func OutputSegmentUsage() {
	for rsIdx, romSegment := range rom {
		fmt.Println(strings.Repeat("─", 64))
		segmentHeader := fmt.Sprintf(" Rom Segment %d - Byte Usage / Bank Size ", rsIdx)
		fmt.Println(util.Colorize(segmentHeader, enumTerminalColors.AnsiTeal, true))
		fmt.Println(strings.Repeat("─", 64))
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
		segmentCaption := util.Colorize(" Total ", enumTerminalColors.AnsiDarkSeaGreen, true)
		bankOutput := fmt.Sprintf("%v %v", segmentCaption, output)
		fmt.Println(strings.Repeat("╌", 64))
		fmt.Println(bankOutput)
		fmt.Println()
	}
}

func calcAndGenOutput(maxBytes int, occupiedBytes int) string {
	unoccupiedBytes := maxBytes - occupiedBytes
	availablePercent := (float64(unoccupiedBytes) / float64(maxBytes)) * 100

	out2 := commaSeparated(maxBytes)
	out1 := util.PadStringLeft(commaSeparated(occupiedBytes), len(out2), ' ')
	out3 := commaSeparated(unoccupiedBytes)
	percentVisual := util.PadStringRight(fmt.Sprintf("(%.2f%%)", availablePercent), 9, ' ')
	percentVisual = strings.ReplaceAll(percentVisual, "(", util.Colorize("(", enumTerminalColors.Cyan, false))
	percentVisual = strings.ReplaceAll(percentVisual, ")", util.Colorize(")", enumTerminalColors.Cyan, false))

	numBars := 0

	for i := 0; i < 8; i++ {
		threshold := float64(i) * 12.5
		if availablePercent > threshold {
			numBars++
		}
	}

	g := int(5 * (float64(numBars) / 8.0))
	r := int(5 * (float64(8-numBars) / 8.0))

	solidBarVisual := strings.Repeat("█", numBars)
	emptyBarVisual := strings.Repeat("░", 8-numBars)
	barVisual := util.ColorizeWithAnsiRGBCode(solidBarVisual, fmt.Sprintf("%d%d0", r, g), false) + util.Colorize(emptyBarVisual, enumTerminalColors.AnsiGray4, false)

	output := fmt.Sprintf("%v %v %v %v │ %v %v remaining",
		out1, util.Colorize("/", enumTerminalColors.Cyan, false), out2,
		percentVisual, barVisual, out3)
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
