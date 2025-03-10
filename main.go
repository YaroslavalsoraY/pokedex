package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	commandsDB := getCommands()
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		_, ok := commandsDB[sc.Text()]
		if !ok {
			fmt.Println("Unknown command")
		} else {
			commandsDB[sc.Text()].callback()
		}
	}
}
