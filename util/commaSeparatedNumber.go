package util

import "fmt"

func CommaSeparatedNumber(value int) string {
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
