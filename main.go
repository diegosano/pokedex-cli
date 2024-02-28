package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/diegosano/pokedex-cli/internal/pokeapi"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	cc := &commandConfig{
		previousLocationURL: nil,
		nextLocationURL:     nil,
	}

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()

		words := normalizeInput(scanner.Text())
		if len(words) == 0 {
			continue
		}

		command := words[0]

		cmd, exists := getCommands()[string(command)]
		if exists {
			err := cmd.callback(cc)
			if err != nil {
				fmt.Println(err)
			}
			continue
		} else {
			fmt.Println("Unknown command")
			continue
		}
	}
}

func normalizeInput(text string) []string {
	lowered := strings.ToLower(text)
	words := strings.Fields(lowered)
	return words
}

type commandConfig struct {
	previousLocationURL *string
	nextLocationURL     *string
}

func commandExit(cfg *commandConfig) error {
	os.Exit(0)
	return nil
}

func commandHelp(cfg *commandConfig) error {
	fmt.Println()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, cmd := range getCommands() {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	fmt.Println()
	return nil
}

func commandMap(cfg *commandConfig) error {
	result, err := pokeapi.GetLocationArea(cfg.nextLocationURL)
	if err != nil {
		return err
	}
	for _, location := range result.Results {
		fmt.Println(location.Name)
	}
	cfg.previousLocationURL = result.Previous
	cfg.nextLocationURL = result.Next
	return nil
}

func commandMapB(cfg *commandConfig) error {
	if cfg.previousLocationURL == nil {
		return errors.New("you cannot move back from first page of locations")
	}
	result, err := pokeapi.GetLocationArea(cfg.previousLocationURL)
	if err != nil {
		return err
	}
	for _, location := range result.Results {
		fmt.Println(location.Name)
	}
	cfg.previousLocationURL = result.Previous
	cfg.nextLocationURL = result.Next
	return nil
}

type cliCommand struct {
	name        string
	description string
	callback    func(*commandConfig) error
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Display next 20 location names",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Display previous 20 location names",
			callback:    commandMapB,
		},
	}
}
