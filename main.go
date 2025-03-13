package main

import (
	"bufio"
	"fmt"
	"os"

	pokecache "github.com/YaroslavalsoraY/pokedex/internal"
)

func main() {
	first_adress := "https://pokeapi.co/api/v2/location-area?offset=0&limit=20"
	conf := config{
		next:     &first_adress,
		previous: nil,
		cache:    pokecache.NewCache(5),
	}
	commandsDB := getCommands()
	fmt.Print("Pokedex > ")
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		command := cleanInput(sc.Text())
		if len(command) == 0 {
			fmt.Println("Unknown command")
			fmt.Print("Pokedex > ")
			continue
		}
		_, ok := commandsDB[command[0]]
		if !ok {
			fmt.Println("Unknown command")
			fmt.Print("Pokedex > ")
		} else if len(command) == 1 {
			commandsDB[command[0]].callback(&conf, "")
			fmt.Print("Pokedex > ")
		} else {
			commandsDB[command[0]].callback(&conf, command[1])
			fmt.Print("Pokedex > ")
		}
	}
}
