package game

import (
	"fmt"
	"time"

	"github.com/jensderond/sudoku-cli/internal/sudoku"
)

// Game state
type Game struct {
	Sudoku     sudoku.Sudoku
	Difficulty sudoku.Difficulty
	Lives      int
	StartTime  time.Time
	Elapsed    time.Duration
	Solved     bool
	GameOver   bool
}

// Create a new game
func New(difficulty sudoku.Difficulty) *Game {
	return &Game{
		Sudoku:     sudoku.New(difficulty),
		Difficulty: difficulty,
		Lives:      3,
		StartTime:  time.Now(),
		Solved:     false,
		GameOver:   false,
	}
}

// Reset game with new puzzle
func (g *Game) Reset() {
	g.Sudoku = sudoku.New(g.Difficulty)
	g.Lives = 3
	g.StartTime = time.Now()
	g.Elapsed = 0
	g.Solved = false
	g.GameOver = false
}

// Update elapsed time
func (g *Game) UpdateTime() {
	if !g.Solved && !g.GameOver {
		g.Elapsed = time.Since(g.StartTime)
	}
}

// Switch to next difficulty
func (g *Game) SwitchDifficulty() {
	g.Difficulty = g.Difficulty.Next()
}

// Handle number input
func (g *Game) HandleNumberInput(num int) bool {
	if g.Solved || g.GameOver {
		return false
	}

	oldValue := g.Sudoku.Grid[g.Sudoku.CursorY][g.Sudoku.CursorX]
	
	if !g.Sudoku.SetValue(num) {
		return false // Cannot modify initial cells
	}

	// Check if the move is incorrect
	if !g.Sudoku.IsCurrentMoveCorrect() && oldValue != num {
		g.Lives--
		if g.Lives <= 0 {
			g.GameOver = true
		}
	}

	// Check if solved
	if g.Sudoku.IsSolved() {
		g.Solved = true
	}

	return true
}

// Handle delete/clear input
func (g *Game) HandleClear() bool {
	if g.Solved || g.GameOver {
		return false
	}
	return g.Sudoku.ClearCurrentCell()
}

// Handle cursor movement
func (g *Game) HandleMovement(dx, dy int) {
	if g.GameOver {
		return
	}
	g.Sudoku.MoveCursor(dx, dy)
}

// Get formatted time string
func (g *Game) GetTimeString() string {
	minutes := int(g.Elapsed.Minutes())
	seconds := int(g.Elapsed.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

// Get lives display
func (g *Game) GetLivesDisplay() string {
	display := ""
	for i := 0; i < 3; i++ {
		if i < g.Lives {
			display += "❤️ "
		} else {
			display += "🩶 "
		}
	}
	return display
}
