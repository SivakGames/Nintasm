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
		segmentHeader := fmt.Sprintf(" Rom Segment %d - Byte Usage / Bank Size │ %% free ", rsIdx)
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
		segmentTotalCaption := util.Colorize(" Total ", enumTerminalColors.AnsiDarkSeaGreen, true)
		segmentTotalOutput := fmt.Sprintf("%v%v %v", strings.Repeat(" ", bankCaptionNumberWidth), segmentTotalCaption, output)
		fmt.Println(strings.Repeat("╌", 64))
		fmt.Println(segmentTotalOutput)
		fmt.Println()
	}
}

func calcAndGenOutput(maxBytes int, occupiedBytes int) string {
	unoccupiedBytes := maxBytes - occupiedBytes
	availablePercent := (float64(unoccupiedBytes) / float64(maxBytes)) * 100

	out2 := util.CommaSeparatedNumber(maxBytes)
	out1 := util.PadStringLeft(util.CommaSeparatedNumber(occupiedBytes), len(out2), ' ')
	//out3 := commaSeparated(unoccupiedBytes)
	percentVisual := util.PadStringLeft(fmt.Sprintf("(%.2f%%)", availablePercent), 9, ' ')
	percentVisual = strings.ReplaceAll(percentVisual, "(", util.Colorize("(", enumTerminalColors.Cyan, false))
	percentVisual = strings.ReplaceAll(percentVisual, ")", util.Colorize(")", enumTerminalColors.Cyan, false))

	barVisual := util.GeneratePercentBar(availablePercent)

	output := fmt.Sprintf("%v %v %v │ %v %v",
		out1, util.Colorize("/", enumTerminalColors.Cyan, false), out2,
		barVisual, percentVisual)
	return output
}
