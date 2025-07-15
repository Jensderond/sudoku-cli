package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Difficulty levels
type Difficulty int

func (d Difficulty) Next() Difficulty {
    return (d + 1) % 4
}

const (
	Easy Difficulty = iota
	Medium
	Hard
	Expert
)

// Sudoku grid and game state
type Sudoku struct {
	grid     [9][9]int  // Current grid state
	solution [9][9]int  // Complete solution
	initial  [9][9]bool // Which cells were given initially
	cursorX  int
	cursorY  int
}

// Model for BubbleTea
type model struct {
	sudoku     Sudoku
	difficulty Difficulty
	showErrors bool
	solved     bool
	gameOver   bool
	lives      int
	startTime  time.Time
	elapsed    time.Duration
	keys       keyMap
	help       help.Model
}

// Timer tick message
type tickMsg time.Time

// Key bindings
type keyMap struct {
	Up     key.Binding
	Down   key.Binding
	Left   key.Binding
	Right  key.Binding
	Num    key.Binding
	Delete key.Binding
	New    key.Binding
	Quit   key.Binding
	Help   key.Binding
	Difficulty  key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.Num, k.Delete},
		{k.New, k.Help, k.Quit},
	}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "move right"),
	),
	Num: key.NewBinding(
		key.WithKeys("1", "2", "3", "4", "5", "6", "7", "8", "9"),
		key.WithHelp("1-9", "enter number"),
	),
	Delete: key.NewBinding(
		key.WithKeys("delete", "backspace", "0", "x"),
		key.WithHelp("del/x/0", "clear cell"),
	),
	New: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "new game"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Difficulty: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "switch difficulty"),
	),
}

// Styles
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			MarginBottom(1)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginTop(1)

	livesStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true)

	timerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("46"))
)

// Initialize a new game
func initialModel(difficulty Difficulty) model {
	m := model{
		difficulty: difficulty,
		keys:       keys,
		help:       help.New(),
		lives:      3,
		startTime:  time.Now(),
	}
	m.sudoku = generateSudoku(difficulty)
	return m
}

// Generate a new Sudoku puzzle
func generateSudoku(difficulty Difficulty) Sudoku {
	s := Sudoku{}

	// Generate a complete valid grid
	generateCompleteGrid(&s.solution)

	// Copy solution to current grid
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			s.grid[i][j] = s.solution[i][j]
		}
	}

	// Remove numbers based on difficulty
	cellsToRemove := getCellsToRemove(difficulty)
	removed := 0
	attempts := 0

	for removed < cellsToRemove && attempts < 1000 {
		x := rand.Intn(9)
		y := rand.Intn(9)

		if s.grid[x][y] != 0 {
			backup := s.grid[x][y]
			s.grid[x][y] = 0

			// Check if puzzle still has unique solution
			if hasUniqueSolution(s.grid) {
				removed++
			} else {
				s.grid[x][y] = backup
			}
		}
		attempts++
	}

	// Mark initial cells
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			s.initial[i][j] = s.grid[i][j] != 0
		}
	}

	return s
}

// Generate a complete valid Sudoku grid
func generateCompleteGrid(grid *[9][9]int) {
	// Fill diagonal 3x3 boxes first (they don't affect each other)
	for i := 0; i < 9; i += 3 {
		fillBox(grid, i, i)
	}

	// Fill remaining cells using backtracking
	solveSudoku(grid, 0, 0)
}

// Fill a 3x3 box with random valid numbers
func fillBox(grid *[9][9]int, startRow, startCol int) {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	shuffle(nums)

	index := 0
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			grid[startRow+i][startCol+j] = nums[index]
			index++
		}
	}
}

// Shuffle slice
func shuffle(nums []int) {
	for i := len(nums) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		nums[i], nums[j] = nums[j], nums[i]
	}
}

// Solve Sudoku using backtracking
func solveSudoku(grid *[9][9]int, row, col int) bool {
	if row == 9 {
		return true
	}

	nextRow, nextCol := getNextCell(row, col)

	if grid[row][col] != 0 {
		return solveSudoku(grid, nextRow, nextCol)
	}

	// Try numbers 1-9 in random order
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	shuffle(nums)

	for _, num := range nums {
		if isValid(grid, row, col, num) {
			grid[row][col] = num
			if solveSudoku(grid, nextRow, nextCol) {
				return true
			}
			grid[row][col] = 0
		}
	}

	return false
}

// Get next cell position
func getNextCell(row, col int) (int, int) {
	col++
	if col == 9 {
		col = 0
		row++
	}
	return row, col
}

// Check if a number is valid at a position
func isValid(grid *[9][9]int, row, col, num int) bool {
	// Check row
	for i := 0; i < 9; i++ {
		if grid[row][i] == num {
			return false
		}
	}

	// Check column
	for i := 0; i < 9; i++ {
		if grid[i][col] == num {
			return false
		}
	}

	// Check 3x3 box
	boxRow, boxCol := (row/3)*3, (col/3)*3
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if grid[boxRow+i][boxCol+j] == num {
				return false
			}
		}
	}

	return true
}

// Check if puzzle has unique solution
func hasUniqueSolution(grid [9][9]int) bool {
	solutions := 0
	var testGrid [9][9]int

	// Copy grid
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			testGrid[i][j] = grid[i][j]
		}
	}

	countSolutions(&testGrid, 0, 0, &solutions)
	return solutions == 1
}

// Count number of solutions
func countSolutions(grid *[9][9]int, row, col int, count *int) {
	if *count > 1 {
		return // Early exit if more than one solution found
	}

	if row == 9 {
		*count++
		return
	}

	nextRow, nextCol := getNextCell(row, col)

	if grid[row][col] != 0 {
		countSolutions(grid, nextRow, nextCol, count)
		return
	}

	for num := 1; num <= 9; num++ {
		if isValid(grid, row, col, num) {
			grid[row][col] = num
			countSolutions(grid, nextRow, nextCol, count)
			grid[row][col] = 0
		}
	}
}

// Get number of cells to remove based on difficulty
func getCellsToRemove(difficulty Difficulty) int {
	switch difficulty {
	case Easy:
		return 40 + rand.Intn(6)
	case Medium:
		return 46 + rand.Intn(7)
	case Hard:
		return 53 + rand.Intn(6)
	case Expert:
		return 59 + rand.Intn(6)
	default:
		return 40
	}
}

// Check if the puzzle is solved
func (s *Sudoku) isSolved() bool {
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if s.grid[i][j] != s.solution[i][j] {
				return false
			}
		}
	}
	return true
}

// Get the current cell value (for highlighting)
func (s *Sudoku) getCurrentValue() int {
	return s.grid[s.cursorY][s.cursorX]
}

// Timer command
func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// BubbleTea Init
func (m model) Init() tea.Cmd {
	return tickCmd()
}

// BubbleTea Update
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		if !m.solved && !m.gameOver {
			m.elapsed = time.Since(m.startTime)
		}
		return m, tickCmd()

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll

		case key.Matches(msg, m.keys.Difficulty):
			m.difficulty = m.difficulty.Next()

		case key.Matches(msg, m.keys.New):
			m.sudoku = generateSudoku(m.difficulty)
			m.solved = false
			m.gameOver = false
			m.lives = 3
			m.startTime = time.Now()
			m.elapsed = 0

		case key.Matches(msg, m.keys.Up):
			if m.sudoku.cursorY > 0 && !m.gameOver {
				m.sudoku.cursorY--
			}

		case key.Matches(msg, m.keys.Down):
			if m.sudoku.cursorY < 8 && !m.gameOver {
				m.sudoku.cursorY++
			}

		case key.Matches(msg, m.keys.Left):
			if m.sudoku.cursorX > 0 && !m.gameOver {
				m.sudoku.cursorX--
			}

		case key.Matches(msg, m.keys.Right):
			if m.sudoku.cursorX < 8 && !m.gameOver {
				m.sudoku.cursorX++
			}

		case key.Matches(msg, m.keys.Delete):
			if !m.sudoku.initial[m.sudoku.cursorY][m.sudoku.cursorX] && !m.solved && !m.gameOver {
				m.sudoku.grid[m.sudoku.cursorY][m.sudoku.cursorX] = 0
			}

		default:
			// Handle number input
			if len(msg.String()) == 1 && msg.String() >= "1" && msg.String() <= "9" {
				if !m.sudoku.initial[m.sudoku.cursorY][m.sudoku.cursorX] && !m.solved && !m.gameOver {
					num := int(msg.String()[0] - '0')
					oldValue := m.sudoku.grid[m.sudoku.cursorY][m.sudoku.cursorX]
					m.sudoku.grid[m.sudoku.cursorY][m.sudoku.cursorX] = num

					// Check if the move is incorrect
					if num != m.sudoku.solution[m.sudoku.cursorY][m.sudoku.cursorX] && oldValue != num {
						m.lives--
						if m.lives <= 0 {
							m.gameOver = true
						}
					}

					// Check if solved
					if m.sudoku.isSolved() {
						m.solved = true
					}
				}
			}
		}
	}

	return m, nil
}

// BubbleTea View
func (m model) View() string {
	var s strings.Builder

	// Title
	title := titleStyle.Render("🎮 SUDOKU")
	s.WriteString(title + "\n\n")

	// Get current cell value for highlighting
	currentValue := m.sudoku.getCurrentValue()

	// Build the grid with borders
	s.WriteString("┏━━━┯━━━┯━━━┳━━━┯━━━┯━━━┳━━━┯━━━┯━━━┓\n")

	for i := 0; i < 9; i++ {
		s.WriteString("┃")

		for j := 0; j < 9; j++ {
			cell := " "
			if m.sudoku.grid[i][j] != 0 {
				cell = fmt.Sprintf("%d", m.sudoku.grid[i][j])
			}

			// Determine cell appearance based on state
			var cellDisplay string

			// Check if this cell should be highlighted (same number as cursor)
			isHighlighted := currentValue != 0 && m.sudoku.grid[i][j] == currentValue

			if i == m.sudoku.cursorY && j == m.sudoku.cursorX {
				// Current position - highlight with brackets
				cellDisplay = lipgloss.NewStyle().
					Foreground(lipgloss.Color("51")).  // Bright cyan
					Bold(true).
					Render(fmt.Sprintf("[%s]", cell))
			} else if m.sudoku.initial[i][j] {
				// Initial given numbers
				if isHighlighted {
					cellDisplay = lipgloss.NewStyle().
						Foreground(lipgloss.Color("45")).  // Lighter cyan for highlighted
						Bold(true).
						Render(fmt.Sprintf(" %s ", cell))
				} else {
					cellDisplay = lipgloss.NewStyle().
						Foreground(lipgloss.Color("241")).
						Render(fmt.Sprintf(" %s ", cell))
				}
			} else if m.sudoku.grid[i][j] != 0 {
				// User-entered numbers
				if m.sudoku.grid[i][j] == m.sudoku.solution[i][j] {
					// Correct
					if isHighlighted {
						cellDisplay = lipgloss.NewStyle().
							Foreground(lipgloss.Color("45")).  // Lighter cyan for highlighted
							Bold(true).
							Render(fmt.Sprintf(" %s ", cell))
					} else {
						cellDisplay = lipgloss.NewStyle().
							Foreground(lipgloss.Color("46")).
							Render(fmt.Sprintf(" %s ", cell))
					}
				} else {
					// Incorrect
					if isHighlighted {
						cellDisplay = lipgloss.NewStyle().
							Foreground(lipgloss.Color("45")).  // Lighter cyan for highlighted
							Bold(true).
							Render(fmt.Sprintf(" %s ", cell))
					} else {
						cellDisplay = lipgloss.NewStyle().
							Foreground(lipgloss.Color("196")).
							Render(fmt.Sprintf(" %s ", cell))
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

	// Status line
	status := fmt.Sprintf("\nDifficulty: %s", difficultyString(m.difficulty))

	// Lives
	livesDisplay := " | Lives: "
	for i := 0; i < 3; i++ {
		if i < m.lives {
			livesDisplay += "❤️ "
		} else {
			livesDisplay += "🩶 "
		}
	}
	status += livesStyle.Render(livesDisplay)

	// Timer
	minutes := int(m.elapsed.Minutes())
	seconds := int(m.elapsed.Seconds()) % 60
	status += timerStyle.Render(fmt.Sprintf(" | Time: %02d:%02d", minutes, seconds))

	if m.solved {
		status += " | 🎉 SOLVED!"
	} else if m.gameOver {
		status += " | 💀 GAME OVER!"
	}

	s.WriteString(infoStyle.Render(status))

	// Help
	s.WriteString("\n\n" + m.help.View(m.keys))

	return s.String()
}

func difficultyString(d Difficulty) string {
	switch d {
	case Easy:
		return "Easy"
	case Medium:
		return "Medium"
	case Hard:
		return "Hard"
	case Expert:
		return "Expert"
	default:
		return "Unknown"
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	p := tea.NewProgram(initialModel(Medium))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
