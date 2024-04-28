module github.com/harshalshethna12/pokedexcli

go 1.22.2

require (
	internal/pokeapi v1.0.0
	internal/pokecache v1.0.0
)

replace internal/pokeapi => ./internal/pokeapi

replace internal/pokecache => ./internal/pokecache
