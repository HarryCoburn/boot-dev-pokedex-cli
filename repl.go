package main

import (
	"bufio"
	"fmt"
	"strings"
)



func setupREPL(scanner *bufio.Scanner) {
	config := &Config{
		Previous: "",
		Next:     "https://pokeapi.co/api/v2/location-area/",
	}

	commandMap := map[string]cliCommand{}

	commandMap["help"] = cliCommand{
		name:        "help",
		description: "Displays a help message",
		callback: func(config *Config) error {
			fmt.Println("Welcome to the Pokedex!")
			fmt.Printf("Usage:\n\n")
			for _, command := range commandMap {
				fmt.Printf("%s: %s\n", command.name, command.description)
			}
			return nil
		},
	}

	commandMap["exit"] = cliCommand{
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	}

	commandMap["map"] = cliCommand{
		name:        "map",
		description: "Displays the next 20 location areas in the Pokemon world",
		callback:    commandLocationMap,
	}

	runREPL(scanner, commandMap, config)

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
