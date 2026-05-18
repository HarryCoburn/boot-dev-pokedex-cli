package main

import (
	"encoding/json"
	"testing"
)

func TestCommandMapCompleteness(t *testing.T) {
	config := &Config{}
	buildCommands(config)
	expected := []string{"help", "exit", "map", "mapb"}
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
		Next: &startURL,
		apiCaller: func(url string) ([]byte, error) {
			return body, nil
		},
	}

	err := commandMap(config)

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

func TestCommandMapbNilPrevious(t *testing.T) {
	config := &Config{Previous: nil}
	err := commandMapb(config)
	if err != nil {
		t.Errorf("Expected nil, but got %v", err)
	}
}

// Additional tests?
//
// TestCommandMapAPIError

//     Config with apiCaller returning (nil, error)
//     Assert commandMap returns a non-nil error

// TestCommandMapInvalidJSON

//     Config with apiCaller returning ([]byte("not-json"), nil)
//     Assert commandMap returns a non-nil error

// TestCommandMapbUpdatesConfig

//     Config with non-nil Previous and a mock apiCaller returning a valid LocationAreaResponse
//     Assert error is nil, config.Next and config.Previous are updated to values from response

// TestCommandMapbAPIError

//     Config with non-nil Previous and apiCaller returning an error
//     Assert commandMapb returns a non-nil error

// TestCommandMapbInvalidJSON

//     Config with non-nil Previous and apiCaller returning invalid JSON bytes
//     Assert commandMapb returns a non-nil error
