package main

import (
	"log"
	"strconv"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var printer *message.Printer

func init() {
	printer = message.NewPrinter(language.English)
}

// formatStr parses string value received from coinbase and parses it as a
// float before casting it as an integer and writing it to a string variable
// with human formatting e.g. "1234567.89" -> "1,234,567". Bit of a shitty
// hack, definitely could be simplified.
func formatStr(s string) string {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Println(err)
		return s
	}

	out := printer.Sprintf("%d", int(f)) // NOTE: Gimme those commas
	return out
}
