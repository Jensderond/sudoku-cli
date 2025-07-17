package sudoku

// Difficulty levels
type Difficulty int

const (
	Easy Difficulty = iota
	Medium
	Hard
	Expert
)

func (d Difficulty) Next() Difficulty {
	return (d + 1) % 4
}

func (d Difficulty) String() string {
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

// Sudoku grid and game state
type Sudoku struct {
	Grid     [9][9]int  // Current grid state
	Solution [9][9]int  // Complete solution
	Initial  [9][9]bool // Which cells were given initially
	CursorX  int
	CursorY  int
}

// Generate a new Sudoku puzzle
func New(difficulty Difficulty) Sudoku {
	s := Sudoku{}

	// Generate a complete valid grid
	generateCompleteGrid(&s.Solution)

	// Copy solution to current grid
	for i := range s.Grid {
		for j := range s.Grid[i] {
			s.Grid[i][j] = s.Solution[i][j]
		}
	}

	// Remove numbers based on difficulty using optimized strategy
	cellsToRemove := getCellsToRemove(difficulty)
	removeCellsSymmetrically(&s.Grid, cellsToRemove)

	// Mark initial cells
	for i := range s.Initial {
		for j := range s.Initial[i] {
			s.Initial[i][j] = s.Grid[i][j] != 0
		}
	}

	return s
}

// Check if the puzzle is solved
func (s *Sudoku) IsSolved() bool {
	for i := range s.Grid {
		for j := range s.Grid[i] {
			if s.Grid[i][j] != s.Solution[i][j] {
				return false
			}
		}
	}
	return true
}

// Get the current cell value (for highlighting)
func (s *Sudoku) GetCurrentValue() int {
	return s.Grid[s.CursorY][s.CursorX]
}

// Move cursor
func (s *Sudoku) MoveCursor(dx, dy int) {
	newX := s.CursorX + dx
	newY := s.CursorY + dy

	if newX >= 0 && newX < 9 {
		s.CursorX = newX
	}
	if newY >= 0 && newY < 9 {
		s.CursorY = newY
	}
}

// Set value at current cursor position
func (s *Sudoku) SetValue(value int) bool {
	if s.Initial[s.CursorY][s.CursorX] {
		return false // Cannot modify initial cells
	}
	s.Grid[s.CursorY][s.CursorX] = value
	return true
}

// Check if current move is correct
func (s *Sudoku) IsCurrentMoveCorrect() bool {
	return s.Grid[s.CursorY][s.CursorX] == s.Solution[s.CursorY][s.CursorX]
}

// Clear current cell
func (s *Sudoku) ClearCurrentCell() bool {
	if s.Initial[s.CursorY][s.CursorX] {
		return false // Cannot modify initial cells
	}
	s.Grid[s.CursorY][s.CursorX] = 0
	return true
}
