package main

import (
	"bufio"
	"fmt"
	"internal/pokeapi"
	"os"
	"strings"
)

var cliCommands map[string]cliCommand
var cfg config
var commandsWithParameters map[string]bool

func init() {
	cliCommands = map[string]cliCommand{
		"help": {
			name:        "help",
			description: "displays a help message",
			callback:    helpCallback,
		},
		"exit": {
			name:        "exit",
			description: "exit from the program",
			callback:    exitCallback,
		},
		"map": {
			name:        "map",
			description: "displays next 20 map locations",
			callback:    mapCallback,
		},
		"mapb": {
			name:        "mapb",
			description: "displays previous 20 map locations",
			callback:    mapbCallback,
		},
		"explore": {
			name:        "explore",
			description: "explore a map location",
			callback:    exploreCallback,
		},
		"catch": {
			name:        "catch",
			description: "catch a pokemon!",
			callback:    catchCallback,
		},
		"inspect": {
			name:        "inspect",
			description: "inspect a caught pokemon!",
			callback:    inspectCallback,
		},
		"pokedex": {
			name:        "pokedex",
			description: "shows your pokedex!",
			callback:    pokedexCallback,
		},
	}
	initialUrl := "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"
	pokedex := make(map[string]pokeapi.Pokemon, 0)
	cfg = config{
		nextPageUrl:     &initialUrl, // initial value
		previousPageUrl: nil,
		areaName:        nil,
		pokemon:         nil,
		pokedex:         pokedex,
	}

	commandsWithParameters = map[string]bool{
		"explore": true,
		"catch":   true,
		"inspect": true,
	}
}

func main() {

	for {
		fmt.Printf("Pokedex> ")
		reader := bufio.NewReader(os.Stdin)
		scanner := bufio.NewScanner(reader)
		scanner.Scan()
		commandString := strings.TrimSpace(scanner.Text())
		commandStringList := strings.Split(commandString, " ")
		command := commandStringList[0]
		if commandsWithParameters[command] {
			if len(commandStringList) != 2 {
				fmt.Print("value not provided in command\n")
				continue
			}
			parameter := commandStringList[1]

			switch command {
			case "explore":
				cfg.areaName = &parameter
			case "catch", "inspect":
				cfg.pokemon = &parameter
			}
		}
		commandFunc, ok := cliCommands[command]
		if !ok {
			fmt.Printf("Invalid command\n")
			continue
		}
		commandFunc.callback(&cfg)
	}
}
