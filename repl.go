package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	pokecache "github.com/YaroslavalsoraY/pokedex/internal"
)

type cliCommand struct {
	name        string
	decsription string
	callback    func(c *config, param string) error
}

type config struct {
	next     *string
	previous *string
	cache    *pokecache.Cache
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

type AreaLocation struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

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

func commandExit(c *config, param string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *config, param string) error {
	commandsDB := getCommands()
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
	for _, v := range commandsDB {
		fmt.Printf("%s: %s\n", v.name, v.decsription)
	}
	return nil
}

func commandMap(c *config, param string) error {
	locations := JSONresp{}
	_, ok := c.cache.Cached[*c.next]
	if !ok {
		resp, err := http.Get(*c.next)
		if err != nil {
			return errors.New("error with connection")
		}
		body, err := io.ReadAll(resp.Body)
		c.cache.Add(*c.next, body)
		if err != nil {
			return errors.New("problem with reading")
		}
		err = json.Unmarshal(body, &locations)
		if err != nil {
			return errors.New("Unmarshalling error")
		}
		defer resp.Body.Close()
	} else {
		body := c.cache.Cached[*c.next].Val
		err := json.Unmarshal(body, &locations)
		if err != nil {
			return errors.New("Unmarshalling error")
		}
	}
	for _, loc := range locations.Results {
		fmt.Println(loc.Name)
	}
	c.next = locations.Next
	c.previous = locations.Previous
	return nil
}

func commandMapb(c *config, param string) error {
	if c.previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}
	locations := JSONresp{}
	_, ok := c.cache.Cached[*c.previous]
	if !ok {
		resp, err := http.Get(*c.previous)
		if err != nil {
			return errors.New("error with connection")
		}
		body, err := io.ReadAll(resp.Body)
		c.cache.Add(*c.previous, body)
		if err != nil {
			return errors.New("problem with reading")
		}
		err = json.Unmarshal(body, &locations)
		if err != nil {
			return errors.New("Unmarshalling error")
		}
		defer resp.Body.Close()
	} else {
		body := c.cache.Cached[*c.previous].Val
		err := json.Unmarshal(body, &locations)
		if err != nil {
			return errors.New("Unmarshalling error")
		}
	}
	for _, loc := range locations.Results {
		fmt.Println(loc.Name)
	}
	c.next = locations.Next
	c.previous = locations.Previous
	return nil

}

func commandExplore(c *config, area string) error {
	if area == "" {
		fmt.Println("Need argument 'location-area'")
		return nil
	}
	url :=  "https://pokeapi.co/api/v2/location-area/" + area
	pokemons := AreaLocation{}
	_, ok := c.cache.Cached[url]
	if !ok {
		resp, err := http.Get(url)
		if err != nil {
			return errors.New("error with connection")
		}
		body, err := io.ReadAll(resp.Body)
		c.cache.Add(url, body)
		if err != nil {
			return errors.New("problem with reading")
		}
		err = json.Unmarshal(body, &pokemons)
		if err != nil {
			return errors.New("Unmarshalling error")
		}
		defer resp.Body.Close()
	} else {
		body := c.cache.Cached[url].Val
		err := json.Unmarshal(body, &pokemons)
		if err != nil {
			return errors.New("Unmarshalling error")
		}
	}
	for _, pok := range pokemons.PokemonEncounters {
		fmt.Println(pok.Pokemon.Name)
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
			decsription: "Commands manual",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			decsription: "Show locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			decsription: "Show previos locations",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			decsription: "Show pokemons on 'location'",
			callback:    commandExplore,
		},
	}
}
