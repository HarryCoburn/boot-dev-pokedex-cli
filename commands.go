package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/HarryCoburn/boot-dev-pokedex-cli/internal/apiclient"
)

type Config struct { // Holds our place in the API pages
	Next     *string
	Previous *string
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

func commandExit(config *Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *Config, commands map[string]cliCommand) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Printf("Usage:\n\n")
	for _, command := range commands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}

func commandMap(config *Config) error {
	body, err := apiclient.CallAPI(*config.Next)
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
	body, err := apiclient.CallAPI(*config.Previous)
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
