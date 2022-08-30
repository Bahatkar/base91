package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

// base91 character tables, the first one is used for encode and another one for decode
var (
	numberToLetter = map[int64]string{0: "A", 1: "B", 2: "C", 3: "D", 4: "E", 5: "F", 6: "G", 7: "H", 8: "I", 9: "J", 10: "K", 11: "L", 12: "M",
		13: "N", 14: "O", 15: "P", 16: "Q", 17: "R", 18: "S", 19: "T", 20: "U", 21: "V", 22: "W", 23: "X", 24: "Y", 25: "Z", 26: "a", 27: "b",
		28: "c", 29: "d", 30: "e", 31: "f", 32: "g", 33: "h", 34: "i", 35: "j", 36: "k", 37: "l", 38: "m", 39: "n", 40: "o", 41: "p", 42: "q",
		43: "r", 44: "s", 45: "t", 46: "u", 47: "v", 48: "w", 49: "x", 50: "y", 51: "z", 52: "0", 53: "1", 54: "2", 55: "3", 56: "4", 57: "5",
		58: "6", 59: "7", 60: "8", 61: "9", 62: "!", 63: "#", 64: "$", 65: "%", 66: "&", 67: "(", 68: ")", 69: "*", 70: "+", 71: ",", 72: ".",
		73: "/", 74: ":", 75: ";", 76: "<", 77: "=", 78: ">", 79: "?", 80: "@", 81: "[", 82: "]", 83: "^", 84: "_", 85: "`", 86: "{", 87: "|",
		88: "}", 89: "~", 90: "\""}

	letterToNumber = map[string]int{"A": 0, "B": 1, "C": 2, "D": 3, "E": 4, "F": 5, "G": 6, "H": 7, "I": 8, "J": 9, "K": 10, "L": 11, "M": 12,
		"N": 13, "O": 14, "P": 15, "Q": 16, "R": 17, "S": 18, "T": 19, "U": 20, "V": 21, "W": 22, "X": 23, "Y": 24, "Z": 25, "a": 26, "b": 27,
		"c": 28, "d": 29, "e": 30, "f": 31, "g": 32, "h": 33, "i": 34, "j": 35, "k": 36, "l": 37, "m": 38, "n": 39, "o": 40, "p": 41, "q": 42,
		"r": 43, "s": 44, "t": 45, "u": 46, "v": 47, "w": 48, "x": 49, "y": 50, "z": 51, "0": 52, "1": 53, "2": 54, "3": 55, "4": 56, "5": 57,
		"6": 58, "7": 59, "8": 60, "9": 61, "!": 62, "#": 63, "$": 64, "%": 65, "&": 66, "(": 67, ")": 68, "*": 69, "+": 70, ",": 71, ".": 72,
		"/": 73, ":": 74, ";": 75, "<": 76, "=": 77, ">": 78, "?": 79, "@": 80, "[": 81, "]": 82, "^": 83, "_": 84, "`": 85, "{": 86, "|": 87,
		"}": 88, "~": 89, "\"": 90}
)

func main() {
	fileName := os.Args[2]
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	currentTime := time.Now().Format("15.04.05_02.01.2006")
	newFile, err := os.Create(fmt.Sprintf("result_%s.txt", currentTime))
	if err != nil {
		fmt.Println(err)
	}
	defer newFile.Close()

	r := bufio.NewReader(file)
	action := os.Args[1]

	for {
		a, err := r.ReadString('\n')
		//"normalization" means throwing out every whitespace and LF symbols because there aren't that symbols in base91 character table
		line := stringNormalization(a)

		if action == "enc" {
			binaryRow := getBinaryWord(line)
			cipher := encode(binaryRow)

			newFile.WriteString(cipher + "\n")
		} else if action == "dec" {
			binaryRow := cipherToBinary(line)
			decodedRow := decode(binaryRow)

			newFile.WriteString(decodedRow + "\n")
		} else {
			fmt.Println("Second argument must be \"enc\" for encode or \"dec\" for decode")
			break
		}

		if err == io.EOF {
			break
		}
	}
}

func getBinaryWord(word string) string {
	wordAsBytes := []byte(word)
	var binaryLetter, binaryWord string
	for i := 0; i < len(wordAsBytes); i++ {
		binaryLetter = strconv.FormatInt(int64(wordAsBytes[i]), 2)
		binaryLetterLen := len(binaryLetter)
		//FormatInt returns binary value with leading "1", so I need to fill it with 0's to get 8 len value
		if binaryLetterLen < 8 {
			for i := 0; i < 8-binaryLetterLen; i++ {
				binaryLetter = "0" + binaryLetter
			}
		}
		binaryWord = binaryLetter + binaryWord
	}
	return binaryWord
}

func encode(binaryRow string) string {
	var encodedRow string

	/*
		If length of binary row is equal or greater than 13, getting 13 bits of it, transform it to integer, then getting a remainder of
		division by 91 and a quotient. Final integers are numbers of symbols from base91 character table.
		If length of binary row is lower than 13, but integer from transforming it greater than 88, it's ok, so acting normally.
		Else transform binary row to integer, which value is number of character from base91 table.
	*/
	for i := len(binaryRow); len(binaryRow) > 0; i -= 13 {
		if len(binaryRow) >= 13 {
			a := binaryRow[i-13:]
			b, err := strconv.ParseInt(a, 2, 64)
			if err != nil {
				fmt.Println(err)
			}
			encodedRow += numberToLetter[b%91] + numberToLetter[b/91]
			binaryRow = binaryRow[:i-13]
		} else if b, _ := strconv.ParseInt(binaryRow, 2, 64); len(binaryRow) < 13 && b > 88 {
			encodedRow += numberToLetter[b%91] + numberToLetter[b/91]
			break
		} else {
			b, err := strconv.ParseInt(binaryRow, 2, 64)
			if err != nil {
				fmt.Println(err)
			}
			encodedRow += numberToLetter[b]
			break
		}
	}
	return encodedRow
}

func cipherToBinary(word string) string {
	var cipher string
	for i := 0; i < len(word)-1; i += 2 {
		fstLetter := letterToNumber[string(word[i])]
		sndLetter := letterToNumber[string(word[i+1])]
		pairSum := sndLetter*91 + fstLetter

		sumAsBinary := strconv.FormatInt(int64(pairSum), 2)
		lenghtA := len(sumAsBinary)
		if lenghtA < 13 {
			for i := 0; i < 13-lenghtA; i++ {
				sumAsBinary = "0" + sumAsBinary
			}
		}

		cipher = sumAsBinary + cipher
	}

	if len(word)%2 == 1 {
		letter := letterToNumber[string(word[len(word)-1])]
		a := strconv.FormatInt(int64(letter), 2)
		cipher = a + cipher
	}
	return cipher
}

func decode(binaryRow string) string {
	//preparing binary row transform to bytes without a remainder
	if len(binaryRow)%8 != 0 {
		for i := 0; i < len(binaryRow)%8; i++ {
			binaryRow = "0" + binaryRow
		}
	}

	var decodedRow string
	for len(binaryRow) > 0 {
		a, err := strconv.ParseInt(string(binaryRow[len(binaryRow)-8:]), 2, 64)
		if err != nil {
			fmt.Println(err)
		}
		b := byte(a)
		decodedRow += string(b)
		binaryRow = binaryRow[:len(binaryRow)-8]
	}

	return decodedRow
}

func stringNormalization(word string) string {
	a := strings.ReplaceAll(word, "\n", "")
	b := strings.ReplaceAll(a, " ", "")

	return b
}
