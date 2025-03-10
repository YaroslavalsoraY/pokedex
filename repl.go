package main

import (
	"fmt"
	"os"
	"strings"
)

func cleanInput(text string) []string {
	var temp_result, result []string
	newText := strings.TrimSpace(text)
	newText = strings.ToLower(newText)
	temp_result = strings.Split(newText, " ")
	for _, e := range temp_result {
		if len([]byte(e)) == 0 {
			continue
		} else {
			result = append(result, e)
		}
	}
	return result
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	commandsDB := getCommands()
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
	for _, v := range commandsDB {
		fmt.Printf("%s: %s\n", v.name, v.decsription)
	}
	return nil
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name:        "exit",
			decsription: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			decsription: "commands manual",
			callback:    commandHelp,
		},
	}
}

type cliCommand struct {
	name        string
	decsription string
	callback    func() error
}
