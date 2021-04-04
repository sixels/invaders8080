package main

import (
	"os"

	"github.com/protoshark/invaders8080/invaders"
)

func main() {
	args := os.Args[1:]

	romPath := args[0]

	game := invaders.New()
	game.Run(romPath)
}
