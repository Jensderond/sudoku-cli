package sudoku

import "math/rand"

// Generate a complete valid Sudoku grid
func generateCompleteGrid(grid *[9][9]int) {
	// Fill diagonal 3x3 boxes first (they don't affect each other)
	for _, i := range []int{0, 3, 6} {
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
	for i := range [3]struct{}{} {
		for j := range [3]struct{}{} {
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

// Check if puzzle has unique solution
func hasUniqueSolution(grid [9][9]int) bool {
	solutions := 0
	var testGrid [9][9]int

	// Copy grid
	for i := range testGrid {
		for j := range testGrid[i] {
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
