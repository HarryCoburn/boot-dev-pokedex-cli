package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/HarryCoburn/boot-dev-pokedex-cli/internal/apiclient"
	"github.com/HarryCoburn/boot-dev-pokedex-cli/internal/pokecache"
)

type Config struct { // Holds our place in the API pages
	Next      *string
	Previous  *string
	apiCaller func(string) ([]byte, error)
	Cache     *pokecache.Cache
	Commands  map[string]cliCommand
}

type cliCommand struct {
	name        string
	description string
	callback    func(*Config) error
}

type Location struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type LocationAreaResponse struct {
	Next      *string    `json:"next"`
	Previous  *string    `json:"previous"`
	Locations []Location `json:"results"`
}

func buildCommands(config *Config) {
	config.Commands = make(map[string]cliCommand)

	config.Commands["help"] = cliCommand{
		name:        "help",
		description: "Displays a help message",
		callback:    commandHelp,
	}

	config.Commands["exit"] = cliCommand{
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	}

	config.Commands["map"] = cliCommand{
		name:        "map",
		description: "Displays the next 20 location areas in the Pokemon world",
		callback:    commandMap,
	}

	config.Commands["mapb"] = cliCommand{
		name:        "mapb",
		description: "Displays the previous 20 location areas in the Pokemon world",
		callback:    commandMapb,
	}
}

func commandExit(config *Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *Config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Printf("Usage:\n\n")
	for _, command := range config.Commands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}

func commandMap(config *Config) error {
	caller := config.apiCaller
	if caller == nil {
		caller = apiclient.CallAPI
	}
	body, err := caller(*config.Next)
	if err != nil {
		return fmt.Errorf("API call failed")
	}
	var locations LocationAreaResponse
	locationErr := json.Unmarshal(body, &locations)
	if locationErr != nil {
		return fmt.Errorf("Response returned no data to match the Location struct")
	}
	config.Next = locations.Next
	config.Previous = locations.Previous
	for _, location := range locations.Locations {
		fmt.Println(location.Name)
	}

	return nil
}

func commandMapb(config *Config) error {
	if config.Previous == nil {
		fmt.Println("You're on the first page")
		return nil
	}
	caller := config.apiCaller
	if caller == nil {
		caller = apiclient.CallAPI
	}
	body, err := caller(*config.Previous)
	if err != nil {
		return fmt.Errorf("API call failed")
	}
	var locations LocationAreaResponse
	locationErr := json.Unmarshal(body, &locations)
	if locationErr != nil {
		return fmt.Errorf("Response returned no data to match the Location struct")
	}
	config.Next = locations.Next
	config.Previous = locations.Previous
	for _, location := range locations.Locations {
		fmt.Println(location.Name)
	}

	return nil
}
