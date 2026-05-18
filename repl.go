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
	config := &Config{
		Previous:  nil,
		Next:      &startURL,
		apiCaller: nil,
		Cache:     pokecache.NewCache(5 * time.Second),
		Commands:  make(map[string]cliCommand),
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
		if exists {
			command.callback(config)

		} else {
			fmt.Println("Unknown command")
		}
	}
}

func cleanInput(text string) []string {
	stringSlice := strings.Fields(strings.ToLower(text))
	return stringSlice
}
