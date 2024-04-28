package main

import (
	"fmt"
	"internal/pokeapi"
	"math/rand"
	"os"
	"strconv"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

type config struct {
	nextPageUrl     *string
	previousPageUrl *string
	areaName        *string
	pokemon         *string
	pokedex         map[string]pokeapi.Pokemon
}

func helpCallback(cfg *config) error {
	fmt.Printf("\nWelcome to the Pokedex!\n")
	fmt.Printf("Usage:\n\n")
	fmt.Printf("help: Displays a help message\n")
	fmt.Printf("exit: Exit the Pokedex\n\n")
	fmt.Printf("map: Shows the list of next 20 in-game locations\n\n")
	fmt.Printf("mapb: Shows the list of previous 20 in-game locations\n\n")
	fmt.Printf("explore <location_area>: Explore a location area\n\n")
	fmt.Printf("catch <pokemon>: Catch a pokemon\n\n")
	fmt.Printf("inspect <pokemon>: Inspect a caught pokemon\n\n")
	fmt.Printf("pokedex: Shows your pokedex\n\n")
	return nil
}

func exitCallback(cfg *config) error {
	os.Exit(0)
	return nil
}

func mapCallback(cfg *config) error {
	url := cfg.nextPageUrl
	if url == nil {
		fmt.Printf("this is the last page\n")
		return fmt.Errorf("this is the last page")
	}
	locationsApiResponseData, err := pokeapi.GetApiResponse[pokeapi.Location](*url)
	if err != nil {
		return fmt.Errorf("error occured while fetching location data %v", err)
	}
	locations := extractLocations(locationsApiResponseData)
	fmt.Print(locations)

	// update nextUrl/previousUrl
	cfg.nextPageUrl = locationsApiResponseData.Next
	cfg.previousPageUrl = locationsApiResponseData.Previous
	return nil
}

func mapbCallback(cfg *config) error {
	url := cfg.previousPageUrl
	if url == nil {
		fmt.Printf("this is the first page\n")
		return fmt.Errorf("this is the first page")
	}
	locationsApiResponseData, err := pokeapi.GetApiResponse[pokeapi.Location](*url)
	if err != nil {
		return fmt.Errorf("error occured while fetching location data %v", err)
	}
	locations := extractLocations(locationsApiResponseData)
	fmt.Print(locations)

	// update nextUrl/previousUrl
	cfg.nextPageUrl = locationsApiResponseData.Next
	cfg.previousPageUrl = locationsApiResponseData.Previous
	return nil
}

func exploreCallback(cfg *config) error {
	areaName := cfg.areaName
	fmt.Printf("Exploring %s...\n", *areaName)

	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", *areaName)
	exploreApiResponseData, err := pokeapi.GetApiResponse[pokeapi.Explore](url)
	if err != nil {
		return fmt.Errorf("error occured while fetching location data for further exploring %v", err)
	}
	pokemons := findPokemons(exploreApiResponseData)
	fmt.Println("Found Pokemon:")
	fmt.Print(pokemons)
	return nil
}

func catchCallback(cfg *config) error {
	pokemon := cfg.pokemon
	fmt.Printf("Throwing a Pokeball at %s...\n", *pokemon)
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", *pokemon)

	pokemonApiResponseData, err := pokeapi.GetApiResponse[pokeapi.Pokemon](url)

	if err != nil {
		return fmt.Errorf("error occured while fetching %s pokemon data: %v", *pokemon, err)
	}

	// simulate catching pokemon
	randomNumber := rand.Intn(pokemonApiResponseData.BaseExperience + 1)
	threshold, err := strconv.Atoi(fmt.Sprintf("%.0f", 0.7*float64(pokemonApiResponseData.BaseExperience)))
	if err != nil {
		return fmt.Errorf("failed converting threshold to number due to error: %v", err)
	}
	if randomNumber > threshold {
		// catch pokemon
		cfg.pokedex[*pokemon] = pokemonApiResponseData
		fmt.Printf("%s was caught!\n", *pokemon)
	} else {
		fmt.Printf("%s escaped!\n", *pokemon)
	}
	return nil
}

func inspectCallback(cfg *config) error {
	pokemon := cfg.pokemon
	// see if pokemon present in pokedex
	pokemonData, ok := cfg.pokedex[*pokemon]
	if !ok {
		fmt.Print("you have not caught that pokemon.\n")
		return nil
	}
	pokemonDetails := extractPokemonDetails(pokemonData)
	fmt.Println(pokemonDetails)
	return nil
}

func pokedexCallback(cfg *config) error {
	fmt.Println("Your Pokedex:")
	for _, pokemon := range cfg.pokedex {
		pokemonName := pokemon.Name
		fmt.Printf("- %s\n", pokemonName)
	}
	return nil
}
