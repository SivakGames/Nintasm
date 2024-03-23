package util

import (
	"fmt"
	enumTerminalColors "misc/nintasm/constants/enums/terminalColors"
	"strconv"
)

const ANSI_COLOR_BASIS = 38

var colorMap = map[enumTerminalColors.Def]int{
	enumTerminalColors.Red:              31,
	enumTerminalColors.Green:            32,
	enumTerminalColors.Yellow:           33,
	enumTerminalColors.Blue:             34,
	enumTerminalColors.Magenta:          35,
	enumTerminalColors.Cyan:             36,
	enumTerminalColors.LightRed:         91,
	enumTerminalColors.LightGreen:       92,
	enumTerminalColors.LightYellow:      93,
	enumTerminalColors.LightBlue:        94,
	enumTerminalColors.LightMagenta:     95,
	enumTerminalColors.LightCyan:        96,
	enumTerminalColors.AnsiPurple:       genAnsi("104"),
	enumTerminalColors.AnsiTeal:         genAnsi("011"),
	enumTerminalColors.AnsiGreen:        genAnsi("351"),
	enumTerminalColors.AnsiDarkSeaGreen: genAnsi("021"),
	enumTerminalColors.AnsiSeaGreen:     genAnsi("142"),
	enumTerminalColors.AnsiOrange:       genAnsi("520"),
	enumTerminalColors.AnsiRed:          genAnsi("100"),
	enumTerminalColors.AnsiGray1:        236,
	enumTerminalColors.AnsiGray2:        238,
	enumTerminalColors.AnsiGray3:        240,
	enumTerminalColors.AnsiGray4:        242,
	enumTerminalColors.AnsiGray5:        244,
}

func Colorize(text string, color enumTerminalColors.Def, isBG bool) string {
	var colorText string
	isAnsi := color > enumTerminalColors.ANSI_START

	colorValue := colorMap[color]

	if !isAnsi {
		if isBG {
			colorValue += 10
		}
		colorText = fmt.Sprintf("\x1b[%dm", colorValue)
	} else {
		ansiG := ANSI_COLOR_BASIS
		if isBG {
			ansiG += 10
		}
		colorText = fmt.Sprintf("\x1b[%d;5;%dm", ansiG, colorValue)
	}

	return fmt.Sprintf("%v%v\x1b[0m", colorText, text)
}

func genAnsi(value string) int {
	genValue, _ := strconv.ParseInt(value, 6, 64)
	return int(genValue) + 16
}

func ColorizeWithAnsiRGBCode(text string, code string, isBG bool) string {
	ansiG := ANSI_COLOR_BASIS
	if isBG {
		ansiG += 10
	}
	colorValue := genAnsi(code)

	colorText := fmt.Sprintf("\x1b[%d;5;%dm", ansiG, colorValue)
	return fmt.Sprintf("%v%v\x1b[0m", colorText, text)
}
