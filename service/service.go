package service

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kenshaw/escpos"
)

// Max returns the maximum of two integers.
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Repeat returns a string consisting of repeated copies of s.
func Repeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}

// ReplaceThbSymbol replaces the THB currency symbol with an empty string.
func ReplaceThbSymbol(s string) string {
	return strings.ReplaceAll(s, "฿", "")
}

// PrintTextWithSpace prints text with spaces to reach a specified length.
func PrintTextWithSpace(printer *escpos.Escpos, text string, data float64, length int, maxLengthNumber int) {
	dataCurrency := FormatCurrency(data)

	totalLine := len(text) + maxLengthNumber
	spacesNeeded := Max(0, length-totalLine)
	spaces := Repeat(" ", spacesNeeded)

	//Center Spaces calcuate from amount
	leftValueSpaceCenter := maxLengthNumber - len(dataCurrency)
	spacesCenter := Repeat(" ", leftValueSpaceCenter)
	//Center Spaces calcuate from amount
	textData := spaces + text + "      " + spacesCenter + fmt.Sprint(dataCurrency)
	textData = ReplaceThbSymbol(textData)
	printer.SetAlign("right")
	fmt.Println(textData)
	if text == "Total: " {
		printer.Write(textData + "\n")
	} else {
		printer.Write(textData + "\n\n")
	}

}

// PrintTextWithSpace prints text with spaces to reach a specified length.
func PrintTextWithSpaceTest(printer *escpos.Escpos, text string, data float64, length int) {
	dataCurrency := FormatCurrency(data)
	totalLine := len(text) + 10
	spacesNeeded := Max(0, length-totalLine)
	spaces := Repeat(" ", spacesNeeded)
	textData := spaces + text + "      " + fmt.Sprint(dataCurrency)
	textData = ReplaceThbSymbol(textData)
	fmt.Println(textData)
	// fmt.Println("spacesNeeded", spacesNeeded, "spaces", len(spaces), "textLine", totalLine)
	// printer.SetAlign("right")
	// if text == "Total: " {
	// 	printer.Write(textData + "\n")
	// } else {
	// 	printer.Write(textData + "\n\n")
	// }

}

// Format formats a currency value.
func FormatCurrency(value float64) string {
	// Format the number with a thousands separator and two decimal places
	formatted := strconv.FormatFloat(value, 'f', 2, 64)

	// Add thousands separator (',') to the integer part of the number
	parts := strings.Split(formatted, ".")
	integerPart := parts[0]
	decimalPart := parts[1]

	var formattedWithComma string

	// Add a comma as a thousands separator
	for i, digit := range integerPart {
		if i > 0 && (len(integerPart)-i)%3 == 0 {
			formattedWithComma += ","
		}
		formattedWithComma += string(digit)
	}

	// Combine the integer and decimal parts
	if decimalPart != "" {
		formattedWithComma += "." + decimalPart
	}

	return formattedWithComma
}
