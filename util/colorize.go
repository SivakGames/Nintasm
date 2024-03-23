package util

import (
	"fmt"
	"strconv"
	"strings"
)

func Colorize(text string, color string, isBG bool) string {
	var colorText string
	colorValue := 33
	isAnsi := strings.HasPrefix(color, "ansi")

	switch color {
	case "red":
		colorValue = 31
	case "green":
		colorValue = 32
	case "yellow":
		colorValue = 33
	case "blue":
		colorValue = 34
	case "magenta":
		colorValue = 35
	case "cyan":
		colorValue = 36
	case "lightred":
		colorValue = 91
	case "lightgreen":
		colorValue = 92
	case "lightyellow":
		colorValue = 93
	case "lightblue":
		colorValue = 94
	case "lightmagenta":
		colorValue = 95
	case "lightcyan":
		colorValue = 96
	case "ansiPurple":
		colorValue = genAnsi("104")
	case "ansiTeal":
		colorValue = genAnsi("011")
	case "ansiGreen":
		colorValue = genAnsi("351")
	case "ansiSeaGreen":
		colorValue = genAnsi("142")
	case "ansiOrange":
		colorValue = genAnsi("520")
	case "ansiRed":
		colorValue = genAnsi("100")
	case "ansiGray1":
		colorValue = 236
	case "ansiGray2":
		colorValue = 238
	case "ansiGray3":
		colorValue = 240
	case "ansiGray4":
		colorValue = 242
	case "ansiGray5":
		colorValue = 244
	}

	if !isAnsi {
		if isBG {
			colorValue += 10
		}
		colorText = fmt.Sprintf("\x1b[%dm", colorValue)
	} else {
		ansiG := 38
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
