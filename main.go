package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
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
		caughtPokemons:      make(map[string]pokeapi.PokemonWithSpecie),
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
	caughtPokemons      map[string]pokeapi.PokemonWithSpecie
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
	area, err := cfg.pokeapiClient.GetLocationArea(areaName)
	if err != nil {
		return err
	}
	fmt.Println("Exploring " + area.Location.Name + "...")
	fmt.Println("Found pokemons:")
	for _, encounter := range area.PokemonEncounters {
		fmt.Println("- " + encounter.Pokemon.Name)
	}
	return nil
}

func commandCatch(cfg *commandConfig, args ...string) error {
	if len(args) != 1 {
		return errors.New("[catch] -> you must provide one pokemon name")
	}
	pokemonName := args[0]
	pokemon, err := cfg.pokeapiClient.GetPokemonWithSpecie(pokemonName)
	if err != nil {
		return err
	}
	fmt.Println("Throwing a pokeball at " + pokemon.Name + "...")
	captured := catchPokemon(pokemon.CaptureRate)
	if !captured {
		fmt.Println(pokemon.Name + " escaped!")
		return nil
	}
	fmt.Println(pokemon.Name + " was caught!")
	cfg.caughtPokemons[pokemon.Name] = pokemon
	return nil
}

func commandInspect(cfg *commandConfig, args ...string) error {
	if len(args) != 1 {
		return errors.New("[inspect] -> you must provide one pokemon name")
	}
	pokemonName := args[0]
	pokemon, ok := cfg.caughtPokemons[pokemonName]
	if !ok {
		fmt.Println("you have not caught " + pokemon.Name + "yet")
		return nil
	}
	fmt.Println("Name: " + pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Println("Stats: ")
	for _, stat := range pokemon.Stats {
		fmt.Printf(" - %s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types: " + pokemon.Name)
	for _, types := range pokemon.Types {
		fmt.Printf(" - %s\n", types.Type.Name)
	}
	return nil
}

func commandPokedex(cfg *commandConfig, args ...string) error {
	fmt.Println("Your pokemons: ")
	for _, pokemon := range cfg.caughtPokemons {
		fmt.Println("- " + pokemon.Name)
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
		"catch": {
			name:        "catch",
			description: "Try to catch a pokemon",
			callback:    commandCatch,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Show all pokemons that you have",
			callback:    commandPokedex,
		},
	}
}

func catchPokemon(captureRate int) bool {
	if captureRate < 0 || captureRate > 255 {
		fmt.Println("Invalid capture rate. It must be between 0 and 255.")
		return false
	}
	randomNumber := rand.Intn(256)

	return randomNumber <= captureRate
}
