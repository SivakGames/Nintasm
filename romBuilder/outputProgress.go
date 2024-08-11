package romBuilder

import (
	"fmt"
	enumTerminalColors "misc/nintasm/constants/enums/terminalColors"
	"misc/nintasm/util"
)

const PASS_CAPTION_LEN = len("Pass _") + 2

// Will draw the basic captions for passes
func DrawPassCaptions() {
	drawEmptyPass(1)
	drawEmptyPass(2)
}

func drawEmptyPass(passNum int) {
	passBox := fmt.Sprintf(" Pass %d ", passNum)
	passCaption := util.Colorize(passBox, enumTerminalColors.AnsiGray3, true)
	fmt.Println(passCaption, util.GeneratePercentBar(0))
}

// 1111111111111111111111111111111111111111111111111111111111111111

func DrawPass1Progress() {
	totalINESByteSize := float64(GetInesTotalRomSizeInKb())
	romSegmentIndex := GetRomSegmentIndex()
	prevSegmentBytes := 0
	segmentIndex := 0

	for romSegmentIndex > 0 {
		prevSegmentBytes += GetRomSegmentByIndexTotalSize(segmentIndex)
		segmentIndex++
		romSegmentIndex--
	}

	romBankBytes := len(*GetCurrentBankSegmentBytes())
	bankIndex := GetBankIndex()

	currentlyUsed := float64(romBankBytes*bankIndex + prevSegmentBytes)
	percentage := float64((currentlyUsed)/totalINESByteSize) * 100
	DrawPassBar(percentage, 2)
}

func DrawPass1Complete() {
	DrawPassBar(100, 2)
}

// 2222222222222222222222222222222222222222222222222222222222222222

func DrawPass2Progress(value float64) {
	DrawPassBar(value, 1)
}

func DrawPass2Complete() {
	DrawPassBar(100, 1)
}

func DrawPassBar(percentage float64, moveUpAmt int) {
	fmt.Print(fmt.Sprintf("\x1b[0G\x1b[%dA\x1b[%dC", moveUpAmt, PASS_CAPTION_LEN+1))
	passBar := util.GeneratePercentBar(percentage)
	fmt.Print(passBar)
	if (moveUpAmt) > 0 {
		fmt.Print(fmt.Sprintf("\x1b[%dE", moveUpAmt))
	}
}
