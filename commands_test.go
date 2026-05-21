package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/HarryCoburn/boot-dev-pokedex-cli/internal/pokecache"
)

func captureOutput(f func()) string {
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestCommandMapCompleteness(t *testing.T) {
	config := &Config{}
	buildCommands(config)
	expected := []string{"help", "exit", "map", "mapb", "explore", "catch", "inspect"}
	if len(config.Commands) != len(expected) {
		t.Errorf("expected %d commands, got %d", len(expected), len(config.Commands))
	}
	for _, key := range expected {
		cmd, ok := config.Commands[key]
		if !ok {
			t.Errorf("missing command %q", key)
		}

		if cmd.name == "" {
			t.Errorf("command %q has no name", key)
		}
		if cmd.description == "" {
			t.Errorf("command %q has no description", key)
		}
		if cmd.callback == nil {
			t.Errorf("command %q has no callback", key)
		}
	}
}

func TestCommandHelp(t *testing.T) {
	config := &Config{}
	buildCommands(config)
	output := captureOutput(func() {
		commandHelp(config, "")
	})

	if !strings.Contains(output, "Welcome to the Pokedex!") {
		t.Errorf("expected welcome message, got: %s", output)
	}
}

func TestCommandMapUpdatesConfig(t *testing.T) {
	next := "https://example.com/page2"
	prev := "https://example.com/page1"
	fakeResponse := LocationAreaResponse{
		Next:      &next,
		Previous:  &prev,
		Locations: []Location{{Name: "pallet-town", Url: "..."}},
	}
	body, _ := json.Marshal(fakeResponse)

	startURL := "https://example.com/start"
	config := &Config{
		Next:  &startURL,
		Cache: pokecache.NewCache(5 * time.Second),
		apiCaller: func(url string) ([]byte, error) {
			return body, nil
		},
	}

	err := commandMap(config, "")

	if err != nil {
		t.Errorf("Expected nil, but got %v", err)
	}
	if *config.Next != next {
		t.Errorf("Config.Next is incorrect. Set to: %v", config.Next)
	}
	if *config.Previous != prev {
		t.Errorf("Config.Previous is incorrect. Set to: %v", config.Previous)
	}
}

func TestCommandMapFetchError(t *testing.T) {
	startURL := "https://example.com/start"
	config := &Config{
		Next:  &startURL,
		Cache: pokecache.NewCache(5 * time.Second),
		apiCaller: func(url string) ([]byte, error) {
			return nil, fmt.Errorf("network failure")
		},
	}

	err := commandMap(config, "")
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestCommandMapBadJSON(t *testing.T) {
	startURL := "https://example.com/start"
	config := &Config{
		Next:  &startURL,
		Cache: pokecache.NewCache(5 * time.Second),
		apiCaller: func(url string) ([]byte, error) {
			return []byte("not valid json"), nil
		},
	}

	err := commandMap(config, "")
	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestCommandMapNilNext(t *testing.T) {
	config := &Config{Next: nil}
	err := commandMap(config, "")
	if err != nil {
		t.Errorf("Expected nil, but got %v", err)
	}
}

func TestCommandMapbNilPrevious(t *testing.T) {
	config := &Config{Previous: nil}
	err := commandMapb(config, "")
	if err != nil {
		t.Errorf("Expected nil, but got %v", err)
	}
}
