package main

// TODO: какое максимальное число изменений может произойти одном шаге при некотором размере сетки

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"
)

const deadCellMark rune = ' '
const livingCellMark rune = '*'

type Game struct {
	grid        [][]bool
	rowCount    int
	columnCount int
	stepCount   int
}

type cellChange struct {
	i, j     int
	newValue bool
}

func (game *Game) RunWithConsole(msDelay int) {
	delayDuration := time.Duration(msDelay) * time.Millisecond
	changesChan := make(chan []*cellChange)
	game.stepCount = 0
	for {
		game.stepCount++
		err := clearConsole()
		if err != nil {
			fmt.Println(err)
			return
		}
		go func() {
			changesChan <- game.createChanges()
		}()
		// TODO: do changes while Sleep
		game.print()
		changes := <-changesChan
		if len(changes) == 0 {
			break
		}
		game.doChanges(changes)
		time.Sleep(delayDuration)
	}
	fmt.Println("Grid is empty or cannot change")
}

func (game *Game) IsEmpty() bool {
	for i := range game.grid {
		for _, item := range game.grid[i] {
			if item {
				return false
			}
		}
	}
	return true
}

///////////////////////////////////////////////////////////////////////////////

func clearConsole() error {
	command := exec.Command("cmd", "/c", "cls")
	command.Stdout = os.Stdout
	return command.Run()
}

func newGame(width, height int) *Game {
	grid := make([][]bool, height)
	for i := range grid {
		grid[i] = make([]bool, width)
	}
	return &Game{
		grid:        grid,
		columnCount: width,
		rowCount:    height,
		stepCount:   0,
	}
}

func (game *Game) print() {
	// TODO: return sync.WaitGroup
	fmt.Printf("Step: %d\n\n", game.stepCount)
	for i := range game.grid {
		for _, item := range game.grid[i] {
			if item {
				fmt.Printf("%c", livingCellMark)
			} else {
				fmt.Printf("%c", deadCellMark)
			}
		}
		fmt.Println()
	}
}

func (game *Game) createChanges() []*cellChange {
	// TODO: check for empty in here
	changesChan := make(chan *cellChange, game.rowCount)
	var wg sync.WaitGroup
	wg.Add(game.rowCount)
	for rowIndex := range game.grid {
		go func(rowIndex int) {
			for j := range game.grid[rowIndex] {
				change := game.createCellChange(rowIndex, j)
				if change != nil {
					changesChan <- change
				}
			}
			wg.Done()
		}(rowIndex)
	}
	go func() {
		wg.Wait()
		close(changesChan)
	}()
	changes := make([]*cellChange, 0, game.rowCount)
	for change := range changesChan {
		changes = append(changes, change)
	}
	return changes
}

func (game *Game) createCellChange(row, column int) *cellChange {
	count := game.countLivingNeighbors(row, column)
	if !game.grid[row][column] {
		if count == 3 {
			return &cellChange{row, column, true}
		}
	} else {
		if count < 2 || count > 3 {
			return &cellChange{row, column, false}
		}
	}
	return nil
}

func (game *Game) doChanges(changes []*cellChange) {
	for _, change := range changes {
		game.grid[change.i][change.j] = change.newValue
	}
}

func (game *Game) countLivingNeighbors(i, j int) int {
	shifts := []int{-1, 0, 1}
	count := 0
	for _, iShift := range shifts {
		for _, jShift := range shifts {
			if iShift == 0 && jShift == 0 {
				continue
			}
			rowIndex := (game.rowCount + i + iShift) % game.rowCount
			columnIndex := (game.columnCount + j + jShift) % game.columnCount
			if game.grid[rowIndex][columnIndex] {
				count++
			}
		}
	}
	return count
}
