package ui

import (
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/jensderond/sudoku-cli/internal/game"
)

// Model for BubbleTea
type Model struct {
	Game *game.Game
	keys keyMap
	help help.Model
}

// Timer tick message
type tickMsg time.Time

// Key bindings
type keyMap struct {
	Up         key.Binding
	Down       key.Binding
	Left       key.Binding
	Right      key.Binding
	Num        key.Binding
	Delete     key.Binding
	New        key.Binding
	Quit       key.Binding
	Help       key.Binding
	Difficulty key.Binding
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
		{k.New, k.Difficulty, k.Help, k.Quit},
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

// NewModel creates a new UI model
func NewModel(g *game.Game) *Model {
	return &Model{
		Game: g,
		keys: keys,
		help: help.New(),
	}
}

// Timer command
func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Init initializes the model
func (m *Model) Init() tea.Cmd {
	return tickCmd()
}

// Update handles messages
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		m.Game.UpdateTime()
		return m, tickCmd()

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll

		case key.Matches(msg, m.keys.Difficulty):
			m.Game.SwitchDifficulty()

		case key.Matches(msg, m.keys.New):
			m.Game.Reset()

		case key.Matches(msg, m.keys.Up):
			m.Game.HandleMovement(0, -1)

		case key.Matches(msg, m.keys.Down):
			m.Game.HandleMovement(0, 1)

		case key.Matches(msg, m.keys.Left):
			m.Game.HandleMovement(-1, 0)

		case key.Matches(msg, m.keys.Right):
			m.Game.HandleMovement(1, 0)

		case key.Matches(msg, m.keys.Delete):
			m.Game.HandleClear()

		default:
			// Handle number input
			if len(msg.String()) == 1 && msg.String() >= "1" && msg.String() <= "9" {
				num, _ := strconv.Atoi(msg.String())
				m.Game.HandleNumberInput(num)
			}
		}
	}

	return m, nil
}

// View renders the UI
func (m *Model) View() string {
	view := Render(m.Game)
	view += "\n\n" + m.help.View(m.keys)
	return view
}
