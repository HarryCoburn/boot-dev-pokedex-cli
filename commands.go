package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"

	"github.com/HarryCoburn/boot-dev-pokedex-cli/internal/apiclient"
	"github.com/HarryCoburn/boot-dev-pokedex-cli/internal/pokecache"
)

type Config struct { // Holds our place in the API pages
	Start      *string
	Next       *string
	Previous   *string
	PokemonURL *string
	apiCaller  func(string) ([]byte, error)
	Cache      *pokecache.Cache
	Commands   map[string]cliCommand
	Pokedex    map[string]PokemonStats
}

type cliCommand struct {
	name        string
	description string
	callback    func(*Config, string) error
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

type Pokemon struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type PokemonEncounters []struct {
	Pokemon Pokemon `json:"pokemon"`
}

type ExploreResponse struct {
	PokemonEncounters PokemonEncounters `json:"pokemon_encounters"`
}

type PokemonStats struct {
	Name   string `json:"name"`
	Height int    `json:"height"`
	Weight int    `json:"weight"`
	Chance int    `json:"base_experience"`
	Stats  []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
}

func buildCommands(config *Config) {
	config.Commands = make(map[string]cliCommand)
	config.Pokedex = make(map[string]PokemonStats)

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

	config.Commands["explore"] = cliCommand{
		name:        "explore",
		description: "Explore what Pokemon are in a location.",
		callback:    commandExplore,
	}

	config.Commands["catch"] = cliCommand{
		name:        "catch",
		description: "Attempt to catch a Pokemon.",
		callback:    commandCatch,
	}
}

func commandExit(config *Config, p string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *Config, p string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Printf("Usage:\n\n")
	for _, command := range config.Commands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}

func commandMap(config *Config, p string) error {
	if config.Next == nil {
		fmt.Println("You're on the last page")
		return nil
	}

	var locations LocationAreaResponse
	res, err := fetch(config, *config.Next)
	if err != nil {
		return fmt.Errorf("Map failed")
	}

	if err := json.Unmarshal(res, &locations); err != nil {
		return fmt.Errorf("response returned no data to match the Location struct")
	}

	config.Next = locations.Next
	config.Previous = locations.Previous
	for _, location := range locations.Locations {
		fmt.Println(location.Name)
	}
	return nil
}

func commandMapb(config *Config, p string) error {
	if config.Previous == nil {
		fmt.Println("You're on the first page")
		return nil
	}
	var locations LocationAreaResponse
	res, err := fetch(config, *config.Previous)
	if err != nil {
		return fmt.Errorf("Mapb failed")
	}

	if err := json.Unmarshal(res, &locations); err != nil {
		return fmt.Errorf("response returned no data to match the Location struct")
	}

	config.Next = locations.Next
	config.Previous = locations.Previous
	for _, location := range locations.Locations {
		fmt.Println(location.Name)
	}
	return nil
}

func fetch(config *Config, url string) ([]byte, error) {
	caller := config.apiCaller
	if caller == nil {
		caller = apiclient.CallAPI
	}

	body, exists := config.Cache.Get(url)
	if !exists {
		var err error
		body, err = caller(url)
		if err != nil {
			return nil, fmt.Errorf("API call failed")
		}
		config.Cache.Add(url, body)
	}
	return body, nil
}

func commandExplore(config *Config, loc string) error {
	fmt.Println("Exploring " + loc)
	locURL := *config.Start + loc + "/"
	var pokemons ExploreResponse
	res, err := fetch(config, locURL)
	if err != nil {
		return fmt.Errorf("Explore failed")
	}
	if err := json.Unmarshal(res, &pokemons); err != nil {
		return fmt.Errorf("response returned no data to match the PokemonEncounters struct")
	}
	fmt.Println("Pokemon found:")
	for _, pokemon := range pokemons.PokemonEncounters {
		fmt.Println("- " + pokemon.Pokemon.Name)
	}
	return nil
}

func commandCatch(config *Config, p string) error {
	fmt.Printf("Throwing a Pokeball at %s...\n", p)
	var pokemon PokemonStats
	res, err := fetch(config, *config.PokemonURL+p+"/")
	if err != nil {
		return fmt.Errorf("Catch request failed")
	}

	if err := json.Unmarshal(res, &pokemon); err != nil {
		return fmt.Errorf("response returned no data to match the CatchResponse struct")
	}
	fmt.Printf("Catch chance is %d\n", pokemon.Chance)
	catchAttempt := rand.Intn(pokemon.Chance)
	if catchAttempt <= 20 {
		fmt.Printf("%s was caught!\n", p)
		config.Pokedex[p] = pokemon
	} else {
		fmt.Printf("%s escaped!\n", p)
	}
	return nil
}
