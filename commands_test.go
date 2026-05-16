package main

import "testing"

func TestCommandMapCompleteness(t *testing.T) {
	for key, command := range commandMap {
		if command.name == "" {
			t.Errorf("command %q has no name", key)
		}
		if command.description == "" {
			t.Errorf("command %q has no description", key)
		}
		if command.callback == nil {
			t.Errorf("command %q has no callback", key)
		}
	}
}
