package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/AvelarJ/pokedexcli/internal"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	// Load pokedex from file if it exists
	pokedex, err := loadPokedex()
	if err != nil {
		fmt.Println(err)
	} else {
		if len(pokedex.Pokemon) > 0 { // Correctly loaded pokedex
			fmt.Printf("Loaded pokedex with %d pokemon\n", len(pokedex.Pokemon))
		} else { // Otherwise create empty pokedex
			fmt.Println("Created new pokedex, let the adventure begin!")
		}
	}

	config := &config{
		cache:   internal.NewCache(5 * 60 * time.Second), // Cache with a 5-minute expiration interval
		pokedex: pokedex,                                 // Start with an empty pokedex
	} // Start with empty config, will be updated by commands that need it

	for { // Inf loop - run until exit command is called
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		cleanInput := cleanInput(input)
		command := cleanInput[0]

		// Additional parameters for commands that require them, otherwise an empty slice
		var parameters []string
		if len(cleanInput) > 1 {
			parameters = cleanInput[1:]
		} else {
			parameters = []string{}
		}

		// Check if the command exists in the availableCommands map
		if commandFunc, exists := availableCommands[command]; exists {

			err := commandFunc.callback(config, parameters)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Printf("Unknown command: %s\n", command)
		}
	}
}
