package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	first_adress := "https://pokeapi.co/api/v2/location-area?offset=0&limit=20"
	conf := config{
		next: &first_adress,
		previous: nil,
	}
	commandsDB := getCommands()
	fmt.Print("Pokedex > ")
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		_, ok := commandsDB[sc.Text()]
		if !ok {
			fmt.Println("Unknown command")
			fmt.Print("Pokedex > ")
		} else {
			commandsDB[sc.Text()].callback(&conf)
			fmt.Print("Pokedex > ")
		}
	}
}
