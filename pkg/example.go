package main

import (
	"fmt"
	"strconv"
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

func main() {
	number := uint64(123456789)
	shortenedUrl := convertToBase62(number)
	fmt.Printf("Original number: %d\nShortened URL: %s\n", number, shortenedUrl)

	recoveredNumber, _ := convertFromBase62(shortenedUrl)
	fmt.Printf("Recovered number: %d\n", recoveredNumber)

	// Для тестирования произвольного числа
	testNumberStr := "1234567890"
	testNumber, _ := strconv.ParseUint(testNumberStr, 10, 64)
	convertedTestString := convertToBase62(testNumber)
	fmt.Printf("\nTesting conversion of %s:\nShortened URL: %s\n", testNumberStr, convertedTestString)

	recoveredTestNumber, _ := convertFromBase62(convertedTestString)
	fmt.Printf("Recovered number: %d\n", recoveredTestNumber)
}
