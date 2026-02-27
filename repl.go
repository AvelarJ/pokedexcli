package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/AvelarJ/pokedexcli/internal"
)

// Command registry
type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

// Config struct
type config struct {
	Next     *string
	Previous *string
	cache    *internal.Cache
}

var availableCommands map[string]cliCommand

// Initialize the command registry with available commands
func init() {
	availableCommands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Display help information",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Display names of 20 location areas in the pokemon world",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Display the previous 20 location areas in the pokemon world",
			callback:    commandMapb,
		},
	}
}

func cleanInput(text string) []string {
	// Used to split the user input into a slice of strings, removing any extra spaces and converting to lowercase for easier processing.
	//hello world -> ["hello", "world"]
	//Charmander Bulbasaur PIKACHU -> ["charmander", "bulbasaur", "pikachu"]

	lCase := strings.ToLower(text)
	trimmed := strings.TrimSpace(lCase)
	split := strings.Fields(trimmed)

	return split
}

// Command functions

func commandExit(config *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, cmd := range availableCommands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}

	return nil
}

func commandMap(config *config) error {
	//Will fetch from the API then Println the name of 20 location areas
	var apiURL string

	if config.Previous == nil && config.Next == nil {
		apiURL = "https://pokeapi.co/api/v2/location-area?limit=20"
	} else if config.Next != nil {
		apiURL = *config.Next
	} else {
		apiURL = *config.Previous
	}

	locationArea, err := config.cache.FetchLocationAreas(apiURL)
	if err != nil {
		return fmt.Errorf("error fetching location areas: %v", err)
	}

	//Set config for next and previous page
	config.Next = locationArea.Next
	config.Previous = locationArea.Previous

	for _, area := range locationArea.Results {
		fmt.Println(area.Name)
	}

	return nil
}

func commandMapb(config *config) error {
	//(map back) a way to look at the previous 20 location areas, if there are any. If there are no previous location areas, it should print a message saying so.
	if config.Previous == nil {
		fmt.Println("No previous location areas to display.")
		return nil
	}

	apiURL := *config.Previous

	locationArea, err := config.cache.FetchLocationAreas(apiURL)
	if err != nil {
		return fmt.Errorf("error fetching location areas: %v", err)
	}

	//Set config for next and previous page
	config.Next = locationArea.Next
	config.Previous = locationArea.Previous

	for _, area := range locationArea.Results {
		fmt.Println(area.Name)
	}

	return nil
}
