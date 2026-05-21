package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strings"

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
		description: "Explore what Pokemon are in a location. Requires a location from map/mapb",
		callback:    commandExplore,
	}

	config.Commands["catch"] = cliCommand{
		name:        "catch",
		description: "Attempt to catch a Pokemon. Requires the name of a pokemon.",
		callback:    commandCatch,
	}

	config.Commands["inspect"] = cliCommand{
		name:        "inspect",
		description: "Inspect the stats of a pokemon in your pokedex",
		callback:    commandInspect,
	}

	config.Commands["pokedex"] = cliCommand{
		name:        "pokedex",
		description: "Returns the pokemon you have caught for inspection",
		callback:    commandPokedex,
	}
}

func commandExit(config *Config, p string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *Config, p string) error {
	fmt.Println()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Printf("Usage:\n\n")

	keys := make([]string, 0, len(config.Commands))
	for k := range config.Commands {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		cmd := config.Commands[k]
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	fmt.Println()
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
	fmt.Println()
	fmt.Println("Locations:")
	for _, location := range locations.Locations {
		fmt.Println(location.Name)
	}
	fmt.Println()
	return nil
}

func commandMapb(config *Config, p string) error {
	if config.Previous == nil {
		fmt.Println("You're on the first page")
		fmt.Println()
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
	fmt.Println()
	fmt.Println("Locations:")
	for _, location := range locations.Locations {
		fmt.Println(location.Name)
	}
	fmt.Println()
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
	fmt.Println()
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
	fmt.Println()
	return nil
}

func commandCatch(config *Config, p string) error {
	fmt.Println()
	fmt.Printf("Throwing a Pokeball at %s...\n", p)
	var pokemon PokemonStats
	res, err := fetch(config, *config.PokemonURL+p+"/")
	if err != nil {
		return fmt.Errorf("Catch request failed")
	}

	if err := json.Unmarshal(res, &pokemon); err != nil {
		return fmt.Errorf("response returned no data to match the CatchResponse struct")
	}

	catchAttempt := rand.Intn(pokemon.Chance)
	if catchAttempt <= 20 {
		fmt.Printf("%s was caught!\n", strings.ToUpper(p[:1])+p[1:])
		config.Pokedex[p] = pokemon
		fmt.Printf("You may now inspect your %s with the inspect command.", p)
		fmt.Println()
	} else {
		fmt.Printf("%s escaped!\n", strings.ToUpper(p[:1])+p[1:])
	}
	fmt.Println()
	return nil
}

func commandInspect(config *Config, p string) error {
	pokemon, exists := config.Pokedex[p]
	if !exists {
		fmt.Println()
		fmt.Println("You have not caught that pokemon.")
		return nil
	}
	fmt.Println()
	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Printf("Stats:\n")
	for _, stat := range pokemon.Stats {
		fmt.Printf("- %s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Printf("Types:\n")
	for _, poketype := range pokemon.Types {
		fmt.Printf("- %s\n", poketype.Type.Name)
	}
	fmt.Println()
	return nil
}

func commandPokedex(config *Config, p string) error {
	fmt.Println()
	if len(config.Pokedex) == 0 {
		fmt.Println("Your Pokedex is empty! Go catch some Pokemon.")
		fmt.Println()
		return nil
	}
	fmt.Println("Your Pokedex:")
	for _, pokemon := range config.Pokedex {
		fmt.Printf("- %s", pokemon.Name)
	}
	fmt.Println()
	fmt.Println()
	return nil
}
