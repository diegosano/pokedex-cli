package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type PokemonSpecie struct {
	CaptureRate int `json:"capture_rate"`
}

type Pokemon struct {
	Name    string `json:"name"`
	Species struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"species"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	Weight int `json:"weight"`
	Height int `json:"height"`
}

type PokemonWithSpecie struct {
	Pokemon
	PokemonSpecie
}

func (c *Client) GetPokemon(name string) (Pokemon, error) {
	url := baseURL + "/pokemon/" + name

	if val, ok := c.cache.Get(url); ok {
		response := Pokemon{}
		err := json.Unmarshal(val, &response)
		if err != nil {
			return Pokemon{}, err
		}
		return response, nil
	}

	res, err := http.Get(url)
	if err != nil {
		return Pokemon{}, err
	}

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()

	if err != nil {
		return Pokemon{}, err
	}

	if res.StatusCode > 299 {
		return Pokemon{}, fmt.Errorf("pokemon \"%s\" not found", name)
	}

	response := Pokemon{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return Pokemon{}, err
	}

	c.cache.Add(url, body)

	return response, nil
}

func (c *Client) GetPokemonSpecie(name string) (PokemonSpecie, error) {
	url := baseURL + "/pokemon-species/" + name

	if val, ok := c.cache.Get(url); ok {
		response := PokemonSpecie{}
		err := json.Unmarshal(val, &response)
		if err != nil {
			return PokemonSpecie{}, err
		}
		return response, nil
	}

	res, err := http.Get(url)
	if err != nil {
		return PokemonSpecie{}, err
	}

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()

	if err != nil {
		return PokemonSpecie{}, err
	}

	if res.StatusCode > 299 {
		return PokemonSpecie{}, fmt.Errorf("pokemon \"%s\" not found", name)
	}

	response := PokemonSpecie{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return PokemonSpecie{}, err
	}

	c.cache.Add(url, body)

	return response, nil
}

func (c *Client) GetPokemonWithSpecie(name string) (PokemonWithSpecie, error) {
	pokemon, err := c.GetPokemon(name)
	if err != nil {
		return PokemonWithSpecie{}, nil
	}
	specie, err := c.GetPokemonSpecie(name)
	if err != nil {
		return PokemonWithSpecie{}, nil
	}
	return PokemonWithSpecie{
		Pokemon:       pokemon,
		PokemonSpecie: specie,
	}, nil
}
