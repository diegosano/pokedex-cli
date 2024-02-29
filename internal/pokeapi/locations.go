package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type LocationArea struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

type LocationAreas struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func (c *Client) GetLocationAreas(url *string) (LocationAreas, error) {
	defaultURL := baseURL + "/location-area"
	if url != nil {
		defaultURL = *url
	}

	if val, ok := c.cache.Get(defaultURL); ok {
		response := LocationAreas{}
		err := json.Unmarshal(val, &response)
		if err != nil {
			return LocationAreas{}, err
		}
		return response, nil
	}

	res, err := http.Get(defaultURL)
	if err != nil {
		return LocationAreas{}, err
	}

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()

	if err != nil {
		return LocationAreas{}, err
	}

	if res.StatusCode > 299 {
		return LocationAreas{}, fmt.Errorf("response failed with status code: %d and\nbody: %s", res.StatusCode, body)
	}

	response := LocationAreas{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return LocationAreas{}, err
	}

	c.cache.Add(defaultURL, body)

	return response, nil
}

func (c *Client) GetLocationArea(areaName string) (LocationArea, error) {
	url := baseURL + "/location-area/" + areaName

	if val, ok := c.cache.Get(url); ok {
		response := LocationArea{}
		err := json.Unmarshal(val, &response)
		if err != nil {
			return LocationArea{}, err
		}
		return response, nil
	}

	res, err := http.Get(url)
	if err != nil {
		return LocationArea{}, err
	}

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()

	if err != nil {
		return LocationArea{}, err
	}

	if res.StatusCode > 299 {
		return LocationArea{}, fmt.Errorf("location area \"%s\" not found", areaName)
	}

	response := LocationArea{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return LocationArea{}, err
	}

	c.cache.Add(url, body)

	return response, nil
}
