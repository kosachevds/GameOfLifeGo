package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func CreateWithGrid(grid [][]bool) (*Game, error) {
	rowsCount := len(grid)
	columnsCount := len(grid[0])
	for i := range grid[1:] {
		if len(grid[i]) != columnsCount {
			return nil, fmt.Errorf("rows sizes must be the same")
		}
	}
	return &Game{
		grid:        grid,
		rowCount:    rowsCount,
		columnCount: columnsCount,
		stepCount:   0,
	}, nil
}

func CreateWithInt(grid [][]int) (*Game, error) {
	boolGrid := make([][]bool, len(grid))
	for i := range grid {
		boolGrid[i] = make([]bool, len(grid[i]))
		for j := range grid[i] {
			boolGrid[i][j] = grid[i][j] != 0
		}
	}
	return CreateWithGrid(boolGrid)
}

func ReadFromFileWithInt(filename string) (*Game, error) {
	grid, err := readGrid(filename, '1')
	if err != nil {
		return nil, fmt.Errorf("grid reading error: %e", err)
	}
	return CreateWithGrid(grid)
}

func CreateRandom(width, height int, livingCellCount int) *Game {
	rand.Seed(time.Now().UTC().UnixNano())
	game := newGame(width, height)
	currentLivingCount := 0
	for currentLivingCount < livingCellCount {
		var row, column int
		for {
			row = rand.Intn(height)
			column = rand.Intn(width)
			if !game.grid[row][column] {
				break
			}
		}
		game.grid[row][column] = true
		currentLivingCount++
	}
	return game
}

///////////////////////////////////////////////////////////////////////////////

func readGrid(filename string, livingMark rune) ([][]bool, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("open file error: %e", err)
	}
	scanner := bufio.NewScanner(file)
	lines := readLinesAsync(scanner, 2)
	grid := parseLines(lines, livingMark)
	return grid, nil
}

func readLinesAsync(scanner *bufio.Scanner, chanBufferSize int) <-chan string {
	lines := make(chan string, chanBufferSize)
	go func() {
		for scanner.Scan() {
			lines <- scanner.Text()
		}
		close(lines)
	}()
	return lines
}

func parseLines(lines <-chan string, livingMark rune) [][]bool {
	grid := make([][]bool, 0)
	for line := range lines {
		lineLength := len(line)
		row := make([]bool, lineLength)
		for i, item := range line {
			row[i] = item == livingMark
		}
		grid = append(grid, row)
	}
	return grid
}
