package main

import "fmt"

/*
Description:

	Encryptor for the site domain

attr:

	line - line for encryption
	compression - lavel of compression. Range [min: 1, max: 6]
*/
type Scrambler struct {
	ConvertBytes      map[int]uint64
	MemoryOfLongBytes []uint8
	baseChars         string
}

func NewScrambler() *Scrambler {
	return &Scrambler{
		baseChars:         "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz",
		ConvertBytes:      make(map[int]uint64, 10),
		MemoryOfLongBytes: make([]uint8, 10, 10),
	}
}

func (s *Scrambler) changeBitOfByte(currentByte uint8, shift int) uint8 {
	DeltaByte := byte(1 << shift)
	return currentByte | DeltaByte
}

// Slices bytes [104 116 116 112 115 58 47 47 112 114 97 99 116 105
// 				99 117 109 46 121 97 110 100 101 120 46 114 117 47]

//  After PreparForConversion:

// Байтовое представление     / Битовое представление разрядов
//	    ConvertBytes         /        MemoryOfLongBytes
// 0: 104 116 116 112 115 58 | 00111110
// 1: 47 47 112 114 97 99    | 00001100
// 2: 116 105 99 117 109 46  | 00110110
// 3: 121 97 110 100 101 120 | 00101111
// 4: 46 114 117 47          | 00011000

func (s *Scrambler) PreparForConversion(line []byte, compression int) {
	var increment uint64 = 100
	var tmpResult uint64 = 0
	cntB := 0
	cntU := 0

	for i, b := range line {

		// Проверка, что значение b (1 byte) 3-x значное. Например 116
		if 0 < (uint8(b) / uint8(100)) {
			increment = 1000

			// В массиве MemoryOfLongBytes запоминаем разрядность чисел, где
			//     0 -это пустота/двуразрядное число
			//     1 -это трехразрядное число
			newByte := s.changeBitOfByte(s.MemoryOfLongBytes[cntU], (compression - cntB - 1))
			s.MemoryOfLongBytes[cntU] = newByte

		} else {
			// Тут логика с 2-х значными значениями. Например 99
			increment = 100
		}

		// Собираем число bytes с количеством == compression
		tmpResult = (tmpResult * increment) + uint64(b)

		// Собрали нужное число из byte. Добавляем значение в ConvertBytes
		if cntB == compression-1 || (i == len(line)-1) {
			s.ConvertBytes[cntU] = tmpResult
			tmpResult = 0
			cntB = 0
			cntU++
		} else {
			cntB++
		}
	}
}

func (s *Scrambler) Execute(line []byte, compression int) string {
	// Подготовка данных для последующей конвертации
	s.PreparForConversion(line, compression)
	// Конвертация

	var result string

	for _, v := range s.ConvertBytes {

		for v > 0 {
			tmp := v % 62
			v /= 62
			result = string(s.baseChars[int(tmp)]) + result
		}

		fmt.Println(result)

		break
	}

	return "QWER"
}

func main() {

	line := "https://practicum.yandex.ru/"
	lineByte := []byte(line)

	scrambler := NewScrambler()
	scrambler.Execute(lineByte, 6)

	fmt.Println(scrambler.ConvertBytes)

	for _, v := range scrambler.MemoryOfLongBytes {
		fmt.Printf("%08b\n", v)
	}

}
