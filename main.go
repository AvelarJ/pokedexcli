package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	config := &config{} //Start with empty config, will be updated by commands that need it
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
