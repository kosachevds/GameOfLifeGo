package main

func main() {
	// game, err := ReadFromFileWithInt("./grids/glider")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	game := CreateRandom(40, 20, 200)
	game.RunWithConsole(200)
}
