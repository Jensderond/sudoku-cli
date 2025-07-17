package sudoku

import "math/rand"

// Generate a complete valid Sudoku grid
func generateCompleteGrid(grid *[9][9]int) {
	// Used number trackers for rows, columns, and boxes
	var rowUsed, colUsed, boxUsed [9][10]bool

	// Fill diagonal 3x3 boxes first (they don't affect each other)
	for _, i := range []int{0, 3, 6} {
		fillBox(grid, i, i, &rowUsed, &colUsed, &boxUsed)
	}

	// Fill remaining cells using backtracking
	solveSudokuFast(grid, 0, 0, &rowUsed, &colUsed, &boxUsed)
}

// Fill a 3x3 box with random valid numbers
func fillBox(grid *[9][9]int, startRow, startCol int, rowUsed, colUsed, boxUsed *[9][10]bool) {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	shuffle(nums)

	index := 0
	for i := range 3 {
		for j := range 3 {
			val := nums[index]
			grid[startRow+i][startCol+j] = val
			rowUsed[startRow+i][val] = true
			colUsed[startCol+j][val] = true
			boxUsed[(startRow/3)*3+(startCol/3)][val] = true
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

// Optimized solveSudoku using O(1) validity checks
func solveSudokuFast(grid *[9][9]int, row, col int, rowUsed, colUsed, boxUsed *[9][10]bool) bool {
	if row == 9 {
		return true
	}

	nextRow, nextCol := getNextCell(row, col)

	if grid[row][col] != 0 {
		return solveSudokuFast(grid, nextRow, nextCol, rowUsed, colUsed, boxUsed)
	}

	// Try numbers 1-9 in random order
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	shuffle(nums)
	boxIdx := (row/3)*3 + (col / 3)

	for _, num := range nums {
		if !rowUsed[row][num] && !colUsed[col][num] && !boxUsed[boxIdx][num] {
			grid[row][col] = num
			rowUsed[row][num] = true
			colUsed[col][num] = true
			boxUsed[boxIdx][num] = true

			if solveSudokuFast(grid, nextRow, nextCol, rowUsed, colUsed, boxUsed) {
				return true
			}

			grid[row][col] = 0
			rowUsed[row][num] = false
			colUsed[col][num] = false
			boxUsed[boxIdx][num] = false
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
	for i := range grid[row] {
		if grid[row][i] == num {
			return false
		}
	}

	// Check column
	for i := range grid {
		if grid[i][col] == num {
			return false
		}
	}

	// Check 3x3 box
	boxRow, boxCol := (row/3)*3, (col/3)*3
	for i := range [3]struct{}{} {
		for j := range [3]struct{}{} {
			if grid[boxRow+i][boxCol+j] == num {
				return false
			}
		}
	}

	return true
}

// Check if puzzle has unique solution (optimized version)
func hasUniqueSolution(grid [9][9]int) bool {
	var rowUsed, colUsed, boxUsed [9][10]bool

	// Initialize tracking arrays
	for i := range grid {
		for j := range grid[i] {
			if grid[i][j] != 0 {
				rowUsed[i][grid[i][j]] = true
				colUsed[j][grid[i][j]] = true
				boxUsed[(i/3)*3+(j/3)][grid[i][j]] = true
			}
		}
	}

	solutions := 0
	var testGrid [9][9]int
	copy2DArray(&testGrid, &grid)

	countSolutions(&testGrid, 0, 0, &solutions, &rowUsed, &colUsed, &boxUsed)
	return solutions == 1
}

// Count number of solutions (optimized version using tracking arrays)
func countSolutions(grid *[9][9]int, row, col int, count *int, rowUsed, colUsed, boxUsed *[9][10]bool) {
	if *count > 1 {
		return // Early exit if more than one solution found
	}

	if row == 9 {
		*count++
		return
	}

	nextRow, nextCol := getNextCell(row, col)

	if grid[row][col] != 0 {
		countSolutions(grid, nextRow, nextCol, count, rowUsed, colUsed, boxUsed)
		return
	}

	boxIdx := (row/3)*3 + (col/3)

	for num := range 9 {
		num++ // Convert 0-based to 1-based
		// O(1) validity check using tracking arrays
		if !rowUsed[row][num] && !colUsed[col][num] && !boxUsed[boxIdx][num] {
			grid[row][col] = num
			rowUsed[row][num] = true
			colUsed[col][num] = true
			boxUsed[boxIdx][num] = true

			countSolutions(grid, nextRow, nextCol, count, rowUsed, colUsed, boxUsed)

			grid[row][col] = 0
			rowUsed[row][num] = false
			colUsed[col][num] = false
			boxUsed[boxIdx][num] = false
		}
	}
}

// Helper to copy 2D array
func copy2DArray(dst, src *[9][9]int) {
	for i := range src {
		for j := range src[i] {
			dst[i][j] = src[i][j]
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

// Remove cells symmetrically to maintain puzzle quality while reducing checks
func removeCellsSymmetrically(grid *[9][9]int, targetRemoval int) {
	// Create a list of all cell positions
	type cell struct {
		row, col int
	}

	cells := make([]cell, 0, 81)
	for i := range 9 {
		for j := range 9 {
			cells = append(cells, cell{i, j})
		}
	}

	// Shuffle cells for randomness
	rand.Shuffle(len(cells), func(i, j int) {
		cells[i], cells[j] = cells[j], cells[i]
	})

	removed := 0
	attempts := 0
	maxAttempts := len(cells) * 2

	// Try to remove cells
	for removed < targetRemoval && attempts < maxAttempts {
		idx := attempts % len(cells)
		c := cells[idx]

		if grid[c.row][c.col] != 0 {
			backup := grid[c.row][c.col]
			grid[c.row][c.col] = 0

			// Check uniqueness only for first 80% of removals
			// For the last 20%, use heuristics
			if removed < int(float64(targetRemoval)*0.8) {
				if hasUniqueSolution(*grid) {
					removed++

					// Try to remove symmetric cell if possible
					symRow := 8 - c.row
					symCol := 8 - c.col
					if grid[symRow][symCol] != 0 && removed < targetRemoval {
						symBackup := grid[symRow][symCol]
						grid[symRow][symCol] = 0

						if hasUniqueSolution(*grid) {
							removed++
						} else {
							grid[symRow][symCol] = symBackup
						}
					}
				} else {
					grid[c.row][c.col] = backup
				}
			} else {
				// For last 20%, use heuristic: ensure at least 17 clues remain
				// (minimum for unique solution)
				clueCount := 0
				for i := range 9 {
					for j := range 9 {
						if grid[i][j] != 0 {
							clueCount++
						}
					}
				}

				if clueCount > 17 {
					removed++
				} else {
					grid[c.row][c.col] = backup
				}
			}
		}
		attempts++
	}
}
