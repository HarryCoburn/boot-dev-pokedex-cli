package main

import (
	"bufio"
	"fmt"
	"strings"
	"time"

	"github.com/HarryCoburn/boot-dev-pokedex-cli/internal/pokecache"
)

func setupREPL(scanner *bufio.Scanner) {
	startURL := "https://pokeapi.co/api/v2/location-area/"
	pokemonURL := "https://pokeapi.co/api/v2/pokemon/"
	config := &Config{
		Start:      &startURL,
		Next:       &startURL,
		Cache:      pokecache.NewCache(5 * time.Second),
		PokemonURL: &pokemonURL,
	}

	buildCommands(config)

	runREPL(scanner, config)

}

func runREPL(scanner *bufio.Scanner, config *Config) {
	for {
		fmt.Print("Pokedex > ")
		if !scanner.Scan() {
			break
		}
		cleanedInput := cleanInput(scanner.Text())
		if len(cleanedInput) == 0 {
			continue
		}

		command, exists := config.Commands[cleanedInput[0]]
		var param string
		if len(cleanedInput) > 1 {
			param = cleanedInput[1]
		}

		if exists {
			command.callback(config, param)

		} else {
			fmt.Println("Unknown command")
		}
	}
}

func cleanInput(text string) []string {
	stringSlice := strings.Fields(strings.ToLower(text))
	return stringSlice
}
