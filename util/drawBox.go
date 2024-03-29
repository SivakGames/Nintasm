package util

import (
	"fmt"
	enumTerminalColors "misc/nintasm/constants/enums/terminalColors"
	"strings"
)

const BOX_TL = "╭"
const BOX_TB = "─"
const BOX_TR = "╮"
const BOX_LR = "│"
const BOX_BL = "╰"
const BOX_BR = "╯"

func DrawBox(
	content string,
	lineColor enumTerminalColors.Def,
	textColor enumTerminalColors.Def,
	bgColor enumTerminalColors.Def,
	margin bool,
) {
	adjustedContent := fmt.Sprintf(" %v ", content)
	bWidth := len(adjustedContent)

	topLine := BOX_TL + strings.Repeat(BOX_TB, bWidth) + BOX_TR
	botLine := BOX_BL + strings.Repeat(BOX_TB, bWidth) + BOX_BR

	cTopLine := Colorize(topLine, lineColor, false)
	cSideLine := Colorize(BOX_LR, lineColor, false)
	middle := Colorize(adjustedContent, textColor, false)
	cMiddle := Colorize(middle, bgColor, true)
	cBotLine := Colorize(botLine, lineColor, false)

	fmt.Println(cTopLine)
	fmt.Println(fmt.Sprintf("%v%v%v", cSideLine, cMiddle, cSideLine))
	fmt.Println(cBotLine)
	if margin {
		fmt.Println()
	}
	return
}
