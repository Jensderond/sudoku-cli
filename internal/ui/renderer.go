package ui

import (
	"fmt"
	"strings"

	"github.com/jensderond/sudoku-cli/internal/game"
)

// Render the complete UI
func Render(g *game.Game) string {
	var s strings.Builder

	// Title
	title := TitleStyle.Render("🎮 SUDOKU")
	s.WriteString(title + "\n\n")

	// Render grid
	s.WriteString(RenderGrid(g))

	// Status line
	s.WriteString(RenderStatus(g))

	return s.String()
}

// Render the Sudoku grid
func RenderGrid(g *game.Game) string {
	var s strings.Builder
	currentValue := g.Sudoku.GetCurrentValue()

	// Build the grid with borders
	s.WriteString("┏━━━┯━━━┯━━━┳━━━┯━━━┯━━━┳━━━┯━━━┯━━━┓\n")

	for i := 0; i < 9; i++ {
		s.WriteString("┃")

		for j := 0; j < 9; j++ {
			cell := " "
			if g.Sudoku.Grid[i][j] != 0 {
				cell = fmt.Sprintf("%d", g.Sudoku.Grid[i][j])
			}

			// Determine cell appearance based on state
			var cellDisplay string

			// Check if this cell should be highlighted (same number as cursor)
			isHighlighted := currentValue != 0 && g.Sudoku.Grid[i][j] == currentValue

			if i == g.Sudoku.CursorY && j == g.Sudoku.CursorX {
				// Current position - highlight with brackets
				cellDisplay = CursorStyle.Render(fmt.Sprintf("[%s]", cell))
			} else if g.Sudoku.Initial[i][j] {
				// Initial given numbers
				if isHighlighted {
					cellDisplay = HighlightedCellStyle.Render(fmt.Sprintf(" %s ", cell))
				} else {
					cellDisplay = InitialCellStyle.Render(fmt.Sprintf(" %s ", cell))
				}
			} else if g.Sudoku.Grid[i][j] != 0 {
				// User-entered numbers
				if g.Sudoku.Grid[i][j] == g.Sudoku.Solution[i][j] {
					// Correct
					if isHighlighted {
						cellDisplay = HighlightedCellStyle.Render(fmt.Sprintf(" %s ", cell))
					} else {
						cellDisplay = CorrectCellStyle.Render(fmt.Sprintf(" %s ", cell))
					}
				} else {
					// Incorrect
					if isHighlighted {
						cellDisplay = HighlightedCellStyle.Render(fmt.Sprintf(" %s ", cell))
					} else {
						cellDisplay = IncorrectCellStyle.Render(fmt.Sprintf(" %s ", cell))
					}
				}
			} else {
				// Empty cell
				cellDisplay = fmt.Sprintf(" %s ", cell)
			}

			s.WriteString(cellDisplay)

			// Add vertical separator
			if j < 8 {
				if (j+1)%3 == 0 {
					s.WriteString("┃")
				} else {
					s.WriteString("│")
				}
			}
		}
		s.WriteString("┃\n")

		// Add horizontal separator
		if i < 8 {
			if (i+1)%3 == 0 {
				s.WriteString("┣━━━┿━━━┿━━━╋━━━┿━━━┿━━━╋━━━┿━━━┿━━━┫\n")
			} else {
				s.WriteString("┠───┼───┼───╂───┼───┼───╂───┼───┼───┨\n")
			}
		}
	}

	s.WriteString("┗━━━┷━━━┷━━━┻━━━┷━━━┷━━━┻━━━┷━━━┷━━━┛\n")

	return s.String()
}

// Render the status line
func RenderStatus(g *game.Game) string {
	status := fmt.Sprintf("\nDifficulty: %s", g.Difficulty)

	// Lives
	livesDisplay := " | Lives: " + g.GetLivesDisplay()
	status += LivesStyle.Render(livesDisplay)

	// Timer
	status += TimerStyle.Render(fmt.Sprintf(" | Time: %s", g.GetTimeString()))

	if g.Solved {
		status += " | 🎉 SOLVED!"
	} else if g.GameOver {
		status += " | 💀 GAME OVER!"
	}

	return InfoStyle.Render(status)
}
