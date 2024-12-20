package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*commandConfig) error
}

type commandConfig struct {
	next     string
	previous string
	pageNum  int
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

		command, ok := supportedCommands[commandName]
		if ok {
			command.callback(&config)
		} else {
			fmt.Println("Unknown command")
		}
	}
}

func commandExit(config *commandConfig) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *commandConfig) error {
	helpMsg := "\nWelcome to the Pokedex!\nUsage:\n\n"
	for _, c := range supportedCommands {
		helpMsg += fmt.Sprintf("%v: %v\n", c.name, c.description)
	}
	fmt.Println(helpMsg)
	return nil
}

func commandMap(config *commandConfig) error {
	if config.next == "" {
		config.next = "https://pokeapi.co/api/v2/location-area/"
	}
	list, err := getFromAPI[locationAreaList](config.next)
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

func commandMapb(config *commandConfig) error {
	if config.previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	if config.pageNum == 1 {
		fmt.Println("you're on the first page")
		return nil
	}
	list, err := getFromAPI[locationAreaList](config.previous)
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
