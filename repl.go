package main

import "strings"

func cleanInput(text string) []string {
	// Used to split the user input into a slice of strings, removing any extra spaces and converting to lowercase for easier processing.
	//hello world -> ["hello", "world"]
	//Charmander Bulbasaur PIKACHU -> ["charmander", "bulbasaur", "pikachu"]

	lCase := strings.ToLower(text)
	trimmed := strings.TrimSpace(lCase)
	split := strings.Fields(trimmed)

	return split
}
