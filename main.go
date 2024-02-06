package main

import (
	"fmt"
	"misc/nintasm/assemble"
	"time"
)

func main() {
	lines := []string{" .ds 2,4"}
	//lines := make([]string, 0x2000)
	//for i := range lines {
	//	lines[i] = " lda [555], y "
	//}

	start := time.Now()
	err := assemble.Start(lines)
	if err != nil {
		fmt.Println(err)
	}
	assemblyTime := fmt.Sprintf("%.2f", time.Since(start).Seconds())
	finalMessage := fmt.Sprintf("Assembly took: \x1b[33m%v\x1b[0m seconds", assemblyTime)
	fmt.Println(finalMessage)
}
