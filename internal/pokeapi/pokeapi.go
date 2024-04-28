package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"internal/pokecache"
)

type Location struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"Results"`
}

type Explore struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type Pokemon struct {
	BaseExperience int    `json:"base_experience"`
	Name           string `json:"name"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Types          []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
}

func (pokemon *Pokemon) GetStat(statName string) (statValue int) {

	for _, stat := range pokemon.Stats {
		if stat.Stat.Name == statName {
			statValue = stat.BaseStat
		}
	}
	return statValue
}

func (pokemon *Pokemon) GetTypes() (types string) {
	for _, pokemonType := range pokemon.Types {
		types += "- " + pokemonType.Type.Name + "\n"
	}
	return types
}

type ApiResponse interface {
	Location | Explore | Pokemon
}

var cache pokecache.Cache

func init() {
	// Initialize cache
	cache = pokecache.NewCache(time.Duration(time.Minute * 2))
}

func getHTTPResponse(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("could not fetch data from URL due to some error: %v", err)
	}
	if res.StatusCode > 299 {
		return nil, fmt.Errorf("bad status code: %d", res.StatusCode)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("could not fetch data from URL due to some error: %v", err)
	}
	return body, nil
}

func GetApiResponse[T ApiResponse](url string) (T, error) {
	var zero T

	// Check cache has data already present
	apiResponse, err := cache.Get(url)
	if err != nil {
		// Cache miss, fetch data from url
		apiResponse, err = getHTTPResponse(url)
		if err != nil {
			return zero, fmt.Errorf("could not fetch data due to some error: %v", err)
		}
		cache.Add(url, apiResponse)
	}
	// Unmarshal data
	var apiResponseData T
	err = json.Unmarshal(apiResponse, &apiResponseData)
	if err != nil {
		return zero, fmt.Errorf("could not parse data due to some error: %v", err)
	}
	return apiResponseData, nil
}
