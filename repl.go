package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/AvelarJ/pokedexcli/internal"
)

// Command registry
type cliCommand struct {
	name        string
	description string
	callback    func(*config, []string) error
}

// Config struct
type config struct {
	Next     *string
	Previous *string
	cache    *internal.Cache
	pokedex  *Pokedex
}

// Pokedex struct to store caught pokemon (Can be later expanded)
type Pokedex struct {
	Pokemon map[string]internal.Pokemon
}

var availableCommands map[string]cliCommand

// Initialize the command registry with available commands
func init() {
	availableCommands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Save the current state of the Pokedex and exit",
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
		"explore": {
			name:        "explore",
			description: "Explore a location area by name, showing the pokemon that can be found there (Ex. explore eterna-forest-area)",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Attempt to catch a pokemon by name (Ex. catch pikachu)",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect a pokemon by name (Ex. inspect pikachu)",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "List all pokemon saved in the pokedex",
			callback:    commandPokedex,
		},
	}
}

func loadPokedex() (*Pokedex, error) {
	// Load pokedex data from file if it exists, otherwise return an empty pokedex

	data, err := os.ReadFile("pokedex.json")
	if err != nil {
		if os.IsNotExist(err) {
			return &Pokedex{Pokemon: make(map[string]internal.Pokemon)}, nil
		} else { // Any other error should be returned (Also an empty pokedex)
			return &Pokedex{Pokemon: make(map[string]internal.Pokemon)}, fmt.Errorf("Error opening pokedex.json: %v", err)
		}
	} else {
		var pokedex *Pokedex

		err := json.Unmarshal(data, &pokedex)
		if err != nil { // Returns error and still an empty pokedex
			return &Pokedex{Pokemon: make(map[string]internal.Pokemon)}, fmt.Errorf("Error Unmarshalling pokedex.json: %v", err)
		}
		return pokedex, nil
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

func commandExit(config *config, _ []string) error {
	// Marshal the pokedex data and save it to a file before exiting

	data, err := json.Marshal(config.pokedex)
	if err != nil {
		return fmt.Errorf("error saving pokedex: %v", err)
	}
	err = os.WriteFile("pokedex.json", data, 0644)
	if err != nil {
		return fmt.Errorf("error writing pokedex file: %v", err)
	}

	fmt.Println("Pokedex saved. Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *config, _ []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, cmd := range availableCommands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}

	return nil
}

func commandMap(config *config, _ []string) error {
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

func commandMapb(config *config, _ []string) error {
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

func commandExplore(config *config, parameters []string) error {
	// Will take a location area name as a parameter, fetch data for pokemon found in that location area and print.

	if len(parameters) == 0 {
		return fmt.Errorf("Please provide a location area name to explore (hint: use map command to see location area names).")
	} else if len(parameters) > 1 {
		return fmt.Errorf("Please provide only one location area to explore.")
	}

	locationAreaName := parameters[0]

	fmt.Println("Exploring " + locationAreaName + "...")

	exploreResult, err := config.cache.FetchExploreItem("https://pokeapi.co/api/v2/location-area/" + locationAreaName)
	if err != nil {
		return fmt.Errorf("API request has given an error: %v", err)
	}

	fmt.Println("Found Pokemon:")
	// Loop through and print the names of the pokemon found in the location area
	for _, pokemon := range exploreResult.PokemonEncounters {
		fmt.Println(" - " + pokemon.Pokemon.Name)
	}

	return nil

}

func commandCatch(config *config, parameters []string) error {
	// Will take a pokemon name as a parameter, attempt to catch based on the base_experience and print the result
	// The caught pokemon will be stored in a "pokedex" (Likely a map[string]Pokemon struct) that will be accessed by its own command

	if len(parameters) == 0 {
		return fmt.Errorf("Please provide a pokemon name to catch.")
	} else if len(parameters) > 1 {
		return fmt.Errorf("Please provide only one pokemon name to catch.")
	}

	// Catch probability based on base_experience (Higher base_experience means harder to catch, with a minimum 10% catch rate and a maximum 90% catch rate)
	pokemonData, err := config.cache.FetchPokemon("https://pokeapi.co/api/v2/pokemon/" + parameters[0])
	if err != nil {
		return fmt.Errorf("Error reading pokemon endpoint: %v", err)
	}
	// Print to indicate clean API read and start of catch attempt
	fmt.Println("Throwing a Pokeball at " + parameters[0] + "...")

	// Catch probability: higher BaseExperience = harder to catch, scaled to 0.1–0.9
	const minExp, maxExp = 36.0, 608.0
	normalized := (float64(pokemonData.BaseExperience) - minExp) / (maxExp - minExp)
	catchProbability := 0.9 - normalized*0.8

	if rand.Float64() < catchProbability {
		fmt.Println(parameters[0] + " was caught!")
		config.pokedex.Pokemon[parameters[0]] = pokemonData
		fmt.Println("You may now inspect it with the inspect command.")
	} else {
		fmt.Println(parameters[0] + " escaped!")
	}

	return nil
}

func commandInspect(config *config, parameters []string) error {
	// Will take a pokemon name as a parameter and print it's info if it is in the pokedex (has been caught)

	if len(parameters) == 0 {
		return fmt.Errorf("Please provide a pokemon name to inspect.")
	} else if len(parameters) > 1 {
		return fmt.Errorf("Please provide only one pokemon name to inspect.")
	}

	pokemonName := parameters[0]
	pokemon, exists := config.pokedex.Pokemon[pokemonName]
	if !exists {
		return fmt.Errorf("you have not caught that pokemon")
	} else {
		// Print all of the pokemon's info
		fmt.Println("Name: " + pokemon.Name)
		fmt.Printf("Height: %d\n", pokemon.Height)
		fmt.Printf("Weight: %d\n", pokemon.Weight)
		fmt.Println("Stats:")
		// Loop through and print the stats and types of the pokemon (Not always multiple types)
		for _, stat := range pokemon.Stats {
			fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
		}
		fmt.Println("Types:")
		for _, t := range pokemon.Types {
			fmt.Printf("  - %s\n", t.Type.Name)
		}
	}

	return nil
}

func commandPokedex(config *config, _ []string) error {
	// Will print the names of all the pokemon in the pokedex (pokemon that have been caught)

	fmt.Println("Your pokedex:")
	for name := range config.pokedex.Pokemon {
		fmt.Println(" - " + name)
	}

	return nil
}
