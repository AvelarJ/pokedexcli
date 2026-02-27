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
	config := &config{
		cache: internal.NewCache(5 * 60 * time.Second), // Cache with a 5-minute expiration interval
	} //Start with empty config, will be updated by commands that need it

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		cleanInput := cleanInput(input)
		command := cleanInput[0]

		// Check if the command exists in the availableCommands map
		if commandFunc, exists := availableCommands[command]; exists {

			err := commandFunc.callback(config)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Printf("Unknown command: %s\n", command)
		}
	}
}
