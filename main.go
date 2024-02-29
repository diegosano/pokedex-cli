package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/diegosano/pokedex-cli/internal/pokeapi"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	pokeapiClient := pokeapi.NewClient(time.Minute * 5)
	cc := &commandConfig{
		previousLocationURL: nil,
		nextLocationURL:     nil,
		pokeapiClient:       pokeapiClient,
	}

	commandHelp(cc)

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()

		words := normalizeInput(scanner.Text())
		if len(words) == 0 {
			continue
		}

		command := words[0]
		args := words[1:]

		cmd, exists := getCommands()[command]
		if exists {
			err := cmd.callback(cc, args...)
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
	pokeapiClient       pokeapi.Client
}

func commandExit(cfg *commandConfig, args ...string) error {
	os.Exit(0)
	return nil
}

func commandHelp(cfg *commandConfig, args ...string) error {
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

func commandMap(cfg *commandConfig, args ...string) error {
	result, err := cfg.pokeapiClient.GetLocationAreas(cfg.nextLocationURL)
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

func commandMapB(cfg *commandConfig, args ...string) error {
	if cfg.previousLocationURL == nil {
		return errors.New("[mapb] -> you cannot move back from first page of locations")
	}
	result, err := cfg.pokeapiClient.GetLocationAreas(cfg.previousLocationURL)
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

func commandExplore(cfg *commandConfig, args ...string) error {
	if len(args) != 1 {
		return errors.New("[explore] -> you must provide one location name to explore")
	}
	areaName := args[0]
	result, err := cfg.pokeapiClient.GetLocationArea(areaName)
	if err != nil {
		return err
	}
	fmt.Println("Exploring " + result.Location.Name + "...")
	fmt.Println("Found pokemons:")
	for _, encounter := range result.PokemonEncounters {
		fmt.Println("- " + encounter.Pokemon.Name)
	}
	return nil
}

type cliCommand struct {
	name        string
	description string
	callback    func(*commandConfig, ...string) error
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
		"explore": {
			name:        "explore",
			description: "Explore a location area",
			callback:    commandExplore,
		},
	}
}
