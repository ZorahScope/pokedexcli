package main

import (
	"fmt"
	"strings"
)

func CleanInput(text string) []string {
	lowerCaseInput := strings.ToLower(text)
	splitWords := strings.Fields(lowerCaseInput)
	cleanedInput := make([]string, 0, len(splitWords))

	for _, word := range splitWords {
		cleanedWord := strings.Trim(word, " ")
		cleanedInput = append(cleanedInput, cleanedWord)
	}
	return cleanedInput

}

func main() {
	fmt.Println("Hello, World!")
}
