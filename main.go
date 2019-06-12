package main

import (
	"gameOfLife/gameoflife"
)

func main() {
	// game, err := gameoflife.ReadFromFileWithInt("./grids/glider")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	game := gameoflife.CreateRandom(40, 20, 200)
	game.RunWithConsole(200)
}
