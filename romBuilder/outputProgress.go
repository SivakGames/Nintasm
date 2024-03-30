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
	DrawPass1Bar(percentage)
}

func DrawPass1Complete() {
	DrawPass1Bar(100)
}

func DrawPass1Bar(percentage float64) {
	fmt.Print(fmt.Sprintf("\x1b[0G\x1b[2A\x1b[%dC", PASS_CAPTION_LEN+1))
	passBar := util.GeneratePercentBar(percentage)
	fmt.Println(passBar)
	fmt.Print("\x1b[1E")
}
