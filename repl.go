package main

import (
	"bufio"
	"fmt"
	"strings"
)

func setupREPL(scanner *bufio.Scanner) {
	startURL := "https://pokeapi.co/api/v2/location-area/"
	config := &Config{
		Previous:  nil,
		Next:      &startURL,
		apiCaller: nil,
	}

	commands := buildCommands(config)

	runREPL(scanner, commands, config)

}

func runREPL(scanner *bufio.Scanner, commands map[string]cliCommand, config *Config) {
	for {
		fmt.Print("Pokedex > ")
		if !scanner.Scan() {
			break
		}
		cleanedInput := cleanInput(scanner.Text())
		if len(cleanedInput) == 0 {
			continue
		}

		if cleanedInput[0] == "help" {
			commandHelp(config, commands)
			continue
		}

		command, exists := commands[cleanedInput[0]]
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
