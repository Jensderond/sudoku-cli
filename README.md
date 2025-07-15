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

To build and install the application:

```bash
go build -o sudoku .
mkdir -p ~/bin
cp sudoku ~/bin/sudoku
```

Then run the executable:

```bash
~/bin/sudoku
```

## Installation

To install directly from GitHub:

```bash
go install github.com/jensderond/sudoku-cli@latest
```

Or clone and build:

```bash
git clone https://github.com/jensderond/sudoku-cli.git
cd sudoku-cli
go build -o sudoku .
```

## Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - A powerful little TUI framework

## Development

This is a starting point for a Sudoku CLI application. The current implementation provides:

- Basic menu structure
- Keyboard navigation
- Interactive selection system
- Different difficulty levels

Future enhancements could include:
- Game statistics
