package main

import (
	"bufio"
	"fmt"
	"strings"
)

func runREPL(scanner *bufio.Scanner, commands map[string]cliCommand) {
	for {
		fmt.Print("Pokedex > ")
		if !scanner.Scan() {
			break
		}
		cleanedInput := cleanInput(scanner.Text())
		if len(cleanedInput) == 0 {
			continue
		}

		command, exists := commands[cleanedInput[0]]
		if exists {
			command.callback()
		} else {
			fmt.Println("Unknown command")
		}
	}
}

func cleanInput(text string) []string {
	stringSlice := strings.Fields(strings.ToLower(text))
	return stringSlice
}
