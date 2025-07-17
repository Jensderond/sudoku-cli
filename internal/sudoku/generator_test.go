package sudoku

import (
	"testing"
	"time"
)

// Helper: check if a grid is a valid Sudoku solution
func isValidSudokuGrid(grid *[9][9]int) bool {
	var row, col, box [9][10]bool
	for i := range grid {
		for j := range grid[i] {
			num := grid[i][j]
			if num < 1 || num > 9 {
				return false
			}
			if row[i][num] || col[j][num] || box[(i/3)*3+(j/3)][num] {
				return false
			}
			row[i][num] = true
			col[j][num] = true
			box[(i/3)*3+(j/3)][num] = true
		}
	}
	return true
}

// --- Old generator code for benchmarking ---
func oldFillBox(grid *[9][9]int, startRow, startCol int) {
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

func oldSolveSudoku(grid *[9][9]int, row, col int) bool {
	if row == 9 {
		return true
	}
	nextRow, nextCol := getNextCell(row, col)
	if grid[row][col] != 0 {
		return oldSolveSudoku(grid, nextRow, nextCol)
	}
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	shuffle(nums)
	for _, num := range nums {
		if isValid(grid, row, col, num) {
			grid[row][col] = num
			if oldSolveSudoku(grid, nextRow, nextCol) {
				return true
			}
			grid[row][col] = 0
		}
	}
	return false
}

// Test basic functionality
func TestNewGenerator(t *testing.T) {
	var grid [9][9]int
	var rowUsed, colUsed, boxUsed [9][10]bool
	for _, i := range []int{0, 3, 6} {
		fillBox(&grid, i, i, &rowUsed, &colUsed, &boxUsed)
	}
	if !solveSudokuFast(&grid, 0, 0, &rowUsed, &colUsed, &boxUsed) {
		t.Fatal("New generator failed to generate a grid")
	}
	if !isValidSudokuGrid(&grid) {
		t.Fatal("New generator produced invalid grid")
	}
}

func TestOldGenerator(t *testing.T) {
	var grid [9][9]int
	for _, i := range []int{0, 3, 6} {
		oldFillBox(&grid, i, i)
	}
	if !oldSolveSudoku(&grid, 0, 0) {
		t.Fatal("Old generator failed to generate a grid")
	}
	if !isValidSudokuGrid(&grid) {
		t.Fatal("Old generator produced invalid grid")
	}
}

// Benchmark complete grid generation
func BenchmarkOldGenerator(b *testing.B) {
	for n := 0; n < b.N; n++ {
		var grid [9][9]int
		for _, i := range []int{0, 3, 6} {
			oldFillBox(&grid, i, i)
		}
		if !oldSolveSudoku(&grid, 0, 0) {
			b.Fatal("Old generator failed to generate a grid")
		}
	}
}

func BenchmarkNewGenerator(b *testing.B) {
	for n := 0; n < b.N; n++ {
		var grid [9][9]int
		var rowUsed, colUsed, boxUsed [9][10]bool
		for _, i := range []int{0, 3, 6} {
			fillBox(&grid, i, i, &rowUsed, &colUsed, &boxUsed)
		}
		if !solveSudokuFast(&grid, 0, 0, &rowUsed, &colUsed, &boxUsed) {
			b.Fatal("New generator failed to generate a grid")
		}
	}
}

// Benchmark the generateCompleteGrid function
func BenchmarkGenerateCompleteGrid(b *testing.B) {
	for n := 0; n < b.N; n++ {
		var grid [9][9]int
		generateCompleteGrid(&grid)
	}
}

// Benchmark isValid function (critical for performance)
func BenchmarkIsValid(b *testing.B) {
	var grid [9][9]int
	generateCompleteGrid(&grid)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = isValid(&grid, 4, 4, 5)
	}
}

// Benchmark hasUniqueSolution - THE MAIN BOTTLENECK
func BenchmarkHasUniqueSolution(b *testing.B) {
	var grid [9][9]int
	generateCompleteGrid(&grid)

	// Remove some cells to create a partial puzzle
	removed := 0
	for i := 0; i < 9 && removed < 40; i++ {
		for j := 0; j < 9 && removed < 40; j++ {
			if i*9+j < 40 {
				grid[i][j] = 0
				removed++
			}
		}
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = hasUniqueSolution(grid)
	}
}

// Benchmark solution counting for different numbers of empty cells
func BenchmarkCountSolutions20Empty(b *testing.B) {
	benchmarkCountSolutionsWithEmpty(b, 20)
}

func BenchmarkCountSolutions40Empty(b *testing.B) {
	benchmarkCountSolutionsWithEmpty(b, 40)
}

func BenchmarkCountSolutions60Empty(b *testing.B) {
	benchmarkCountSolutionsWithEmpty(b, 60)
}

func benchmarkCountSolutionsWithEmpty(b *testing.B, emptyCells int) {
	var grid [9][9]int
	generateCompleteGrid(&grid)

	// Remove cells
	removed := 0
	for i := range 9 {
		for j := range 9 {
			if removed < emptyCells {
				grid[i][j] = 0
				removed++
			}
		}
		if removed >= emptyCells {
			break
		}
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		var testGrid [9][9]int
		copy2DArray(&testGrid, &grid)
		solutions := 0
		// Use the new optimized countSolutions function
		var rowUsed, colUsed, boxUsed [9][10]bool
		// Initialize tracking arrays
		for i := range testGrid {
			for j := range testGrid[i] {
				if testGrid[i][j] != 0 {
					rowUsed[i][testGrid[i][j]] = true
					colUsed[j][testGrid[i][j]] = true
					boxUsed[(i/3)*3+(j/3)][testGrid[i][j]] = true
				}
			}
		}
		countSolutions(&testGrid, 0, 0, &solutions, &rowUsed, &colUsed, &boxUsed)
	}
}


// Benchmark puzzle generation by difficulty
func BenchmarkGeneratePuzzleEasy(b *testing.B) {
	benchmarkGeneratePuzzle(b, Easy)
}

func BenchmarkGeneratePuzzleMedium(b *testing.B) {
	benchmarkGeneratePuzzle(b, Medium)
}

func BenchmarkGeneratePuzzleHard(b *testing.B) {
	benchmarkGeneratePuzzle(b, Hard)
}

func BenchmarkGeneratePuzzleExpert(b *testing.B) {
	benchmarkGeneratePuzzle(b, Expert)
}

// Mock the Generate function to benchmark puzzle generation
func benchmarkGeneratePuzzle(b *testing.B, difficulty Difficulty) {
	b.StopTimer()

	for n := 0; n < b.N; n++ {
		var grid [9][9]int
		generateCompleteGrid(&grid)

		cellsToRemove := getCellsToRemove(difficulty)

		b.StartTimer()

		// Use optimized cell removal strategy
		removeCellsSymmetrically(&grid, cellsToRemove)

		b.StopTimer()
	}
}

// Test to measure actual generation times
func TestGenerationTimes(t *testing.T) {
	difficulties := []struct {
		name string
		diff Difficulty
	}{
		{"Easy", Easy},
		{"Medium", Medium},
		{"Hard", Hard},
		{"Expert", Expert},
	}

	for _, d := range difficulties {
		t.Run(d.name, func(t *testing.T) {
			times := make([]time.Duration, 5)

			for i := 0; i < 5; i++ {
				var grid [9][9]int
				generateCompleteGrid(&grid)

				start := time.Now()

				cellsToRemove := getCellsToRemove(d.diff)
				removeCellsSymmetrically(&grid, cellsToRemove)

				times[i] = time.Since(start)
			}

			// Calculate average
			var total time.Duration
			for _, t := range times {
				total += t
			}
			avg := total / 5

			t.Logf("%s difficulty average generation time: %v", d.name, avg)
			for i, dur := range times {
				t.Logf("  Run %d: %v", i+1, dur)
			}
		})
	}
}

// Benchmark the optimized validity check using tracking arrays
func BenchmarkOptimizedValidityCheck(b *testing.B) {
	var grid [9][9]int
	var rowUsed, colUsed, boxUsed [9][10]bool

	// Setup a partial grid
	generateCompleteGrid(&grid)

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

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		// Check if we can place 5 at position (4,4)
		row, col, num := 4, 4, 5
		boxIdx := (row/3)*3 + (col/3)
		_ = !rowUsed[row][num] && !colUsed[col][num] && !boxUsed[boxIdx][num]
	}
}

// Compare old vs new validity check
func BenchmarkCompareValidityChecks(b *testing.B) {
	var grid [9][9]int
	generateCompleteGrid(&grid)
	grid[4][4] = 0 // Clear one cell

	b.Run("Traditional", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			_ = isValid(&grid, 4, 4, 5)
		}
	})

	b.Run("Optimized", func(b *testing.B) {
		var rowUsed, colUsed, boxUsed [9][10]bool
		// Initialize tracking arrays
		for i := 0; i < 9; i++ {
			for j := 0; j < 9; j++ {
				if grid[i][j] != 0 {
					rowUsed[i][grid[i][j]] = true
					colUsed[j][grid[i][j]] = true
					boxUsed[(i/3)*3+(j/3)][grid[i][j]] = true
				}
			}
		}

		b.ResetTimer()
		for n := 0; n < b.N; n++ {
			row, col, num := 4, 4, 5
			boxIdx := (row/3)*3 + (col/3)
			_ = !rowUsed[row][num] && !colUsed[col][num] && !boxUsed[boxIdx][num]
		}
	})
}
