package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jensderond/sudoku-cli/internal/game"
	"github.com/jensderond/sudoku-cli/internal/sudoku"
	"github.com/jensderond/sudoku-cli/internal/ui"
)

func main() {
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Initialize game
	g := game.New(sudoku.Medium)

	// Create UI model
	model := ui.NewModel(g)

	// Create and run the program
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
