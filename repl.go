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

	commands := map[string]cliCommand{}

	commands["help"] = cliCommand{
		name:        "help",
		description: "Displays a help message",
		callback: func(config *Config) error {
			fmt.Println("Welcome to the Pokedex!")
			fmt.Printf("Usage:\n\n")
			for _, command := range commands {
				fmt.Printf("%s: %s\n", command.name, command.description)
			}
			return nil
		},
	}

	commands["exit"] = cliCommand{
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	}

	commands["map"] = cliCommand{
		name:        "map",
		description: "Displays the next 20 location areas in the Pokemon world",
		callback:    commandMap,
	}

	commands["mapb"] = cliCommand{
		name:        "mapb",
		description: "Displays the previous 20 location areas in the Pokemon world",
		callback:    commandMapb,
	}

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
