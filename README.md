# Sudoku CLI

A terminal-based Sudoku game built with Go and Bubble Tea.

## Features

- Interactive menu system
- Terminal-based UI using Bubble Tea
- Keyboard navigation (arrow keys or vim-style j/k)

## Controls

- **Arrow Keys** or **j/k**: Navigate menu items
- **Enter** or **Space**: Select menu item
- **q** or **Ctrl+C**: Quit application

## Installation

1. Clone or download this project
2. Navigate to the project directory
3. Run the application:

```bash
go run main.go
```

## Build

To build the application:

```bash
go build -o sudoku-cli
```

Then run the executable:

```bash
./sudoku-cli
```

## Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - A powerful little TUI framework

## Development

This is a starting point for a Sudoku CLI application. The current implementation provides:

- Basic menu structure
- Keyboard navigation
- Interactive selection system

Future enhancements could include:
- Actual Sudoku game logic
- Save/load functionality
- Different difficulty levels
- Game statistics
