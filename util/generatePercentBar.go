package util

import (
	"fmt"
	enumTerminalColors "misc/nintasm/constants/enums/terminalColors"
	"strings"
)

const ASCII_BAR_WIDTH = 20

func GeneratePercentBar(percentage float64) string {
	numBars := 0

	for i := 0; i < ASCII_BAR_WIDTH; i++ {
		threshold := float64(i) * float64(100/ASCII_BAR_WIDTH)
		if percentage > threshold {
			numBars++
		}
	}

	r := int(5 * (float64(ASCII_BAR_WIDTH-numBars) / ASCII_BAR_WIDTH))
	g := int(5 * (float64(numBars) / ASCII_BAR_WIDTH))
	b := int(float64(g) / 2.5)

	solidBarVisual := strings.Repeat("█", numBars)
	emptyBarVisual := strings.Repeat("░", ASCII_BAR_WIDTH-numBars)

	coloredBars := ColorizeWithAnsiRGBCode(solidBarVisual, fmt.Sprintf("%d%d%d", r, g, b), false)
	grayBars := Colorize(emptyBarVisual, enumTerminalColors.AnsiGray4, false)

	return coloredBars + grayBars
}
