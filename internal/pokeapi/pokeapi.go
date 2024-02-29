package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/diegosano/pokedex-cli/internal/pokecache"
)

const baseURL = "https://pokeapi.co/api/v2/"

type Client struct {
	cache      pokecache.Cache
	httpClient http.Client
}

func NewClient(cacheInterval time.Duration) Client {
	return Client{
		cache:      pokecache.NewCache(cacheInterval),
		httpClient: http.Client{},
	}
}

type LocationResponse struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func (c *Client) GetLocationArea(url *string) (LocationResponse, error) {
	defaultURL := baseURL + "/location-area"
	if url != nil {
		defaultURL = *url
	}

	if val, ok := c.cache.Get(defaultURL); ok {
		response := LocationResponse{}
		err := json.Unmarshal(val, &response)
		if err != nil {
			return LocationResponse{}, err
		}
		return response, nil
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

	response := LocationResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return LocationResponse{}, err
	}

	return response, nil
}
