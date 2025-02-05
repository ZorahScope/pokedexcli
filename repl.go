package main

import (
	"bufio"
	"fmt"
	"github.com/zorahscope/pokedexcli/internal/pokeapi"
	"math/rand"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func(config *commandConfig, args string) error
}

type commandConfig struct {
	next       string
	previous   string
	pageNum    int
	exploreURL string
	pokedex    map[string]pokeapi.Pokemon
}

var supportedCommands map[string]cliCommand

// init initializes the command registry with supported CLI commands.
// This is done during package initialization to avoid circular dependencies
// between the command definitions and their implementations.
func init() {
	supportedCommands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays list of location areas, each subsequent call will return the next page of location areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays list of location areas, each subsequent call will return the previous page of location areas",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Displays list of pokemon at given location",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Attempts to catch designated pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Displays information of captured pokemon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Displays list of pokemon that have been captured",
			callback:    commandPokedex,
		},
	}
}

func cleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}

func startRepl() {
	reader := bufio.NewScanner(os.Stdin)
	config := commandConfig{}
	for {
		fmt.Print("Pokedex > ")
		reader.Scan()

		words := cleanInput(reader.Text())
		if len(words) == 0 {
			continue
		}

		commandName := words[0]
		var args string
		if len(words) > 1 {
			args = words[1]
		}

		command, ok := supportedCommands[commandName]
		if ok {
			command.callback(&config, args)
		} else {
			fmt.Println("Unknown command")
		}
	}
}

func commandExit(config *commandConfig, args string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *commandConfig, args string) error {
	helpMsg := "\nWelcome to the Pokedex!\nUsage:\n\n"
	for _, c := range supportedCommands {
		helpMsg += fmt.Sprintf("%v: %v\n", c.name, c.description)
	}
	fmt.Println(helpMsg)
	return nil
}

func commandMap(config *commandConfig, args string) error {
	if config.next == "" {
		config.next = "https://pokeapi.co/api/v2/location-area/"
	}
	list, err := pokeapi.GetFromAPI[pokeapi.LocationAreaList](config.next)
	if err != nil {
		return fmt.Errorf("error getting data from API: %w", err)
	}
	config.next = list.Next
	if list.Previous != nil {
		config.previous = *list.Previous
	}
	page := ""
	for _, area := range list.Results {
		page += "\n" + area.Name
	}
	fmt.Println(page)
	config.pageNum++
	return nil
}

func commandMapb(config *commandConfig, args string) error {
	if config.previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	if config.pageNum == 1 {
		fmt.Println("you're on the first page")
		return nil
	}
	list, err := pokeapi.GetFromAPI[pokeapi.LocationAreaList](config.previous)
	if err != nil {
		return fmt.Errorf("error getting data from API: %w", err)
	}
	config.next = list.Next
	if list.Previous != nil {
		config.previous = *list.Previous
	}
	page := ""
	for _, area := range list.Results {
		page += "\n" + area.Name
	}
	fmt.Println(page)
	config.pageNum--
	return nil
}

func commandExplore(config *commandConfig, args string) error {
	if config.exploreURL == "" {
		config.exploreURL = "https://pokeapi.co/api/v2/location-area/"
	}
	if config.exploreURL+args == config.exploreURL {
		fmt.Println("Empty argument! Please try again")
		return nil
	}
	list, err := pokeapi.GetFromAPI[pokeapi.LocationArea](config.exploreURL + args)
	if err != nil {
		fmt.Printf("error getting data from API: %v\n", err)
		return fmt.Errorf("error getting data from API: %w", err)
	}
	fmt.Println("Exploring " + args + "...\nFound Pokemon:")
	for _, pokemon := range list.PokemonEncounters {
		fmt.Printf(" - %s\n", pokemon.Pokemon.Name)
	}
	return nil
}

func commandCatch(config *commandConfig, args string) error {
	pokemonURL := "https://pokeapi.co/api/v2/pokemon/"

	if pokemonURL+args == pokemonURL {
		fmt.Println("No pokemon selected! Please try again")
		return nil
	}

	pkmn, err := pokeapi.GetFromAPI[pokeapi.Pokemon](pokemonURL + args)
	if err != nil {
		fmt.Printf("error getting data from API: %v\n", err)
		return fmt.Errorf("error getting data from API: %w", err)
	}
	fmt.Printf("Throwing a Pokeball at %v...\n", pkmn.Name)

	if config.pokedex == nil {
		config.pokedex = make(map[string]pokeapi.Pokemon)
	}

	captureChance := 20.0 / float64(pkmn.BaseExperience)
	randomValue := rand.Float64()

	if randomValue < captureChance {
		config.pokedex[pkmn.Name] = pkmn
		fmt.Printf("%v was caught!\n", pkmn.Name)
		return nil
	}
	fmt.Printf("%v escaped!\n", pkmn.Name)
	return nil
}

func commandInspect(config *commandConfig, args string) error {
	pkmn, ok := config.pokedex[args]
	if !ok {
		fmt.Println("you have not caught that pokemon")
		return nil
	}
	var output strings.Builder

	output.WriteString(fmt.Sprintf("Name: %v\n", pkmn.Name))
	output.WriteString(fmt.Sprintf("Height: %d\n", pkmn.Height))
	output.WriteString(fmt.Sprintf("Weight: %d\n", pkmn.Weight))
	output.WriteString(fmt.Sprintf("Stats: \n"))

	for _, stat := range pkmn.Stats {
		output.WriteString(fmt.Sprintf("  -%v: %v\n", stat.Stat.Name, stat.BaseStat))
	}
	output.WriteString("Types:\n")
	for _, typ := range pkmn.Types {
		output.WriteString(fmt.Sprintf("  - %v\n", typ.Type.Name))
	}

	fmt.Println(output.String())

	return nil
}

func commandPokedex(config *commandConfig, args string) error {
	fmt.Println("Your Pokedex:")
	if len(config.pokedex) == 0 {
		fmt.Println("  - <empty>")
	}
	for k, _ := range config.pokedex {
		fmt.Printf("  - %v\n", k)
	}
	return nil
}
