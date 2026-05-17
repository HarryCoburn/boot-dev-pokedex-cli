package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

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
	Next      string     `json:"next"`
	Previous  string     `json:"previous"`
	Locations []Location `json:"results"`
}

type Config struct { // Holds our place in the API pages
	Next     string
	Previous string
}

var commandMap map[string]cliCommand
var config *Config

func init() {
	commandMap = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays the next 20 location areas in the Pokemon world",
			callback:    commandLocationMap,
		},
	}
	config = &Config{
		Previous: "",
		Next:     "https://pokeapi.co/api/v2/location-area/",
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
	for _, command := range commandMap {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}

func commandLocationMap(config *Config) error {
	res, err := http.Get(config.Next)
	if err != nil {
		return errors.New("Cannot reach the PokeAPI")
	}
	body, err := io.ReadAll(res.Body)

	res.Body.Close()
	fmt.Println(string(body))
	if res.StatusCode > 299 {
		return fmt.Errorf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
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
