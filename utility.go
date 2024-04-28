package main

import (
	"fmt"
	"internal/pokeapi"
)

func extractLocations(locationsApiResponse pokeapi.Location) string {
	// Extract locations
	locations := ""
	for _, location := range locationsApiResponse.Results {
		locations += location.Name + "\n"
	}
	return locations
}

func findPokemons(exploreApiResponseData pokeapi.Explore) string {
	pokemons := ""
	for _, pokemon := range exploreApiResponseData.PokemonEncounters {
		pokemons += "- " + pokemon.Pokemon.Name + "\n"
	}
	return pokemons
}

func extractPokemonDetails(pokemon pokeapi.Pokemon) string {
	pokemonDetails := ""
	pokemonDetails += fmt.Sprintf("Name: %s\n", pokemon.Name)
	pokemonDetails += fmt.Sprintf("Height: %d\n", pokemon.Height)
	pokemonDetails += fmt.Sprintf("Weight: %d\n", pokemon.Weight)
	pokemonDetails += "Stats:\n"
	pokemonDetails += fmt.Sprintf("-hp: %d\n", pokemon.GetStat("hp"))
	pokemonDetails += fmt.Sprintf("-attack: %d\n", pokemon.GetStat("attack"))
	pokemonDetails += fmt.Sprintf("-defense: %d\n", pokemon.GetStat("defense"))
	pokemonDetails += fmt.Sprintf("-special-attack: %d\n", pokemon.GetStat("special-attack"))
	pokemonDetails += fmt.Sprintf("-special-defense: %d\n", pokemon.GetStat("special-defense"))
	pokemonDetails += fmt.Sprintf("-speed: %d\n", pokemon.GetStat("speed"))
	pokemonDetails += "Types:\n"
	pokemonDetails += pokemon.GetTypes()
	return pokemonDetails
}
