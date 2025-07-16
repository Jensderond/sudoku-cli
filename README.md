# Sudoku CLI

A terminal-based Sudoku game built with Go and Bubble Tea.

## Features

- Interactive menu system
- Terminal-based UI using Bubble Tea
- Keyboard navigation (arrow keys or vim-style j/k)

## Controls

- **Arrow Keys** or **j/k**: Navigate menu items
- **Enter** or **Space**: Select menu item
- **d**: Switch difficulty (game needs to be reloaded after switching)
- **q** or **Ctrl+C**: Quit application

## Installation

To install directly from GitHub:

```bash
go install github.com/jensderond/sudoku-cli/cmd/sudoku@latest
```

Or clone and build:

```bash
git clone https://github.com/jensderond/sudoku-cli.git
cd sudoku-cli
go build -o sudoku ./cmd/sudoku
```

## Development

To run the application in development mode:

```bash
go run ./cmd/sudoku
```

To build and install locally:

```bash
go build -o sudoku ./cmd/sudoku
mkdir -p ~/bin
cp sudoku ~/bin/sudoku
```

Then run the executable:

```bash
~/bin/sudoku
```

## Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - A powerful little TUI framework

## Project Status

This is a Sudoku CLI application. The current implementation provides:

- Basic menu structure
- Keyboard navigation
- Interactive selection system
- Different difficulty levels

Future enhancements could include:
- Game statistics
- Save/load functionality
