package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
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

func commandExit(c *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *config) error {
	commandsDB := getCommands()
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
	for _, v := range commandsDB {
		fmt.Printf("%s: %s\n", v.name, v.decsription)
	}
	return nil
}

func commandMap(c *config) error {
	type JSONresp struct {
		Count    int     `json:"count"`
		Next     *string `json:"next"`
		Previous *string `json:"previous"`
		Results  []struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"results"`
	}
	resp, err := http.Get(*c.next)
	if err != nil {
		return errors.New("error with connection")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.New("problem with reading")
	}
	locations := JSONresp{}
	defer resp.Body.Close()
	err = json.Unmarshal(body, &locations)
	if err != nil {
		return errors.New("Unmarshalling error")
	}
	for _, loc := range locations.Results {
		fmt.Println(loc.Name)
	}
	c.next = locations.Next
	c.previous = locations.Previous
	return nil
}

func commandMapb(c *config) error {
	if c.previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}

	type JSONresp struct {
		Count    int     `json:"count"`
		Next     *string `json:"next"`
		Previous *string `json:"previous"`
		Results  []struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"results"`
	}
	resp, err := http.Get(*c.previous)
	if err != nil {
		return errors.New("error with connection")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.New("problem with reading")
	}
	locations := JSONresp{}
	defer resp.Body.Close()
	err = json.Unmarshal(body, &locations)
	if err != nil {
		return errors.New("Unmarshalling error")
	}
	for _, loc := range locations.Results {
		fmt.Println(loc.Name)
	}
	c.next = locations.Next
	c.previous = locations.Previous
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
			decsription: "Commands manual",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			decsription: "Show locations",
			callback:    commandMap,
		},
		"mapb": {
			name: "mapb",
			decsription: "Show previos locations",
			callback: commandMapb,
		},
	}
}

type cliCommand struct {
	name        string
	decsription string
	callback    func(c *config) error
}

type config struct {
	next     *string
	previous *string
}
