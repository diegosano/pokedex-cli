package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const baseURL = "https://pokeapi.co/api/v2/"

type LocationResponse struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func GetLocationArea(url *string) (LocationResponse, error) {
	defaultURL := baseURL + "/location-area"
	if url != nil {
		defaultURL = *url
	}
	res, err := http.Get(defaultURL)
	if err != nil {
		return LocationResponse{}, err
	}

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()

	if err != nil {
		return LocationResponse{}, err
	}

	if res.StatusCode > 299 {
		return LocationResponse{}, fmt.Errorf("response failed with status code: %d and\nbody: %s", res.StatusCode, body)
	}

	result := LocationResponse{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return LocationResponse{}, err
	}

	return result, nil
}
