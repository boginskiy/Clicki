package pkg

import (
	"fmt"
)

const baseChars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// convertToBase62 конвертирует десятичное число в строку в системе счисления с основанием 62
func convertToBase62(num uint64) string {
	var result string
	for num > 0 {
		remainder := num % 62
		num /= 62
		result = string(baseChars[int(remainder)]) + result
	}
	return result
}

// convertFromBase62 восстанавливает исходное число из строки в системе счисления с основанием 62
func convertFromBase62(str string) (uint64, error) {
	var num uint64
	multiplier := uint64(1)
	for i := len(str) - 1; i >= 0; i-- {
		index := -1
		for j, char := range baseChars {
			if rune(str[i]) == char {
				index = j
				break
			}
		}
		if index == -1 {
			return 0, fmt.Errorf("invalid character in the input string: '%c'", str[i])
		}
		num += uint64(index) * multiplier
		multiplier *= 62
	}
	return num, nil
}
