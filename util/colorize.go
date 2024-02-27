package util

import "fmt"

func Colorize(text string, color string, isBG bool) string {
	colorValue := 33
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
	}
	if isBG {
		colorValue += 10
	}

	colorText := fmt.Sprintf("\x1b[%dm", colorValue)
	return fmt.Sprintf("%v%v\x1b[0m", colorText, text)
}
