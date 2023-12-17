package main

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

type Color = string

const (
	BackgroundColor         = lipgloss.Color("#282828")
	PanelBorderColor        = lipgloss.Color("#928374")
	PrimaryTextColor        = lipgloss.Color("#ebdbb2")
	SecondaryTextColor      = lipgloss.Color("#a89984")
	ActiveBorderColor       = lipgloss.Color("#8ec07c")
	CellEmptyColor          = lipgloss.Color("#1d2021")
	CellBackgroundColor     = lipgloss.Color("#bdae93")
	CellBorderColor         = lipgloss.Color("#d5c4a1")
	BombCellBackgroundColor = lipgloss.Color("#cc241d")
	BombColor               = lipgloss.Color("#282828")
	FlagColor               = lipgloss.Color("#fb4934")
)

var (
	tutorial = tutorialText("↑/W", "Move up") +
		"\n\n" +
		tutorialText("↓/S", "Move down") +
		"\n\n" +
		tutorialText("←/A", "Move left") +
		"\n\n" +
		tutorialText("→/D", "Move right") +
		"\n\n" +
		tutorialText("ENTER/LeftM", "Reveal cell") +
		"\n\n" +
		tutorialText("SPACE/RightM", "Flag cell") +
		"\n\n" +
		tutorialText("R", "Restart") +
		"\n\n" +
		tutorialText("Q", "Quit")
)

var (
	baseStyle = lipgloss.NewStyle().
			Background(BackgroundColor)

	primaryTextStyle = lipgloss.NewStyle().
				Inherit(baseStyle).
				Foreground(PrimaryTextColor)

	secondaryTextStyle = lipgloss.NewStyle().
				Inherit(baseStyle).
				Foreground(SecondaryTextColor)

	scoreTextStyle = lipgloss.NewStyle().
			Inherit(baseStyle).
			Foreground(lipgloss.Color("#d3869b"))

	panelStyle = lipgloss.NewStyle().
			Inherit(baseStyle).
			MarginBackground(BackgroundColor).
			BorderBackground(BackgroundColor).
			BorderForeground(PanelBorderColor).
			BorderStyle(lipgloss.NormalBorder())

	cellStyle = lipgloss.NewStyle().
			Width(3).
			Height(1).
			AlignHorizontal(lipgloss.Center).
			AlignVertical(lipgloss.Center).
			BorderStyle(lipgloss.ThickBorder()).
			BorderForeground(CellBorderColor).
			BorderBackground(BackgroundColor).
			Background(CellBackgroundColor).
			Bold(true)

	termStyle = lipgloss.NewStyle().
			Inherit(baseStyle).
			AlignVertical(lipgloss.Center).
			AlignHorizontal(lipgloss.Center)
)

func makeGap(n int) string {
	s := ""

	for i := 0; i < n; i++ {
		s += " "
	}

	return s
}

func (m *Model) View() string {
	w := lipgloss.Width
	h := lipgloss.Height
	leftPanelMarginRight := 1
	leftPanelStyle := panelStyle.Copy().Padding(1, 2).MarginRight(leftPanelMarginRight)

	// Tutorial panel
	tutorialPanel := leftPanelStyle.Render(tutorial)

	// Score panel
	bombText := fmt.Sprintf("%d", m.bombCounter)

	clockText := fmt.Sprintf("%d", m.timer)

	gap := primaryTextStyle.Render(makeGap(w(tutorial) - w(clockText) - w(bombText)))

	scorePanel := leftPanelStyle.Render(scoreTextStyle.Render(bombText) + gap + scoreTextStyle.Render(clockText))

	// Left panel
	leftPanel := lipgloss.JoinVertical(lipgloss.Top, tutorialPanel, scorePanel)

	// Board
	board := m.board()

	// Right panel
	rightPanelPadding := 1
	rightPanel := panelStyle.Copy().Padding(0, rightPanelPadding).Render(board)

	// Calculate coords for each cell
	// TODO: Check if can move to a different function
	paddingHorz := (m.termWidth - (w(leftPanel) + w(rightPanel))) / 2

	paddingVer := (m.termHeight - (h(rightPanel))) / 2

	// Top left cell X coordinate, to get the top left cell border we need to add 1
	X := paddingHorz + w(leftPanel) + leftPanelMarginRight + rightPanelPadding + 1

	// Top right cell Y coordinate, panel border accounts for 1 and top left cell border accounts for the other 1
	Y := paddingVer + 2

	for x := 0; x < GRID_HEIGHT; x++ {
		for y := 0; y < GRID_WIDTH; y++ {
			// [W] -> A cell will have its own width + 2 for each border
			m.cells[x][y].coords.x = X + ((CELL_WIDTH + 2) * y)
			m.cells[x][y].coords.y = Y + ((CELL_HEIGHT + 2) * x)
		}
	}

	ui := lipgloss.JoinHorizontal(0, leftPanel, rightPanel)

	term := termStyle.Width(m.termWidth).Height(m.termHeight)

	return term.Render(ui)
}

func (m *Model) board() string {
	var allRows []string

	for x := 0; x < GRID_HEIGHT; x++ {
		var thisRow []string

		for y := 0; y < GRID_WIDTH; y++ {
			thisRow = append(thisRow, m.renderCell(x, y))
		}

		allRows = append(allRows, lipgloss.JoinHorizontal(0, thisRow...))
	}

	return lipgloss.JoinVertical(0, allRows...)
}

func tutorialText(key string, instruction string) string {
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		primaryTextStyle.Copy().Bold(true).Render(key),
		primaryTextStyle.Copy().Render(" "),
		secondaryTextStyle.Copy().Italic(true).Render(instruction),
	)
}

func (m *Model) renderCell(x int, y int) string {
	cell := m.cells[x][y]
	cellStyle := cellStyle.Copy()
	content := ""

	if &m.cells[x][y] == m.activeCell {
		cellStyle.BorderForeground(ActiveBorderColor)
	}

	if cell.state == FLAGGED {
		content = "⚑"
		cellStyle.Foreground(FlagColor)
		cellStyle.Bold(true)
	}

	if cell.state == OPENED {
		cellContent, cellColor := getCell(cell.value)

		content = cellContent
		cellStyle.Background(CellEmptyColor)
		cellStyle.Foreground(lipgloss.Color(cellColor)).Bold(true)

		if cell.value == BOMB {
			cellStyle.Background(BombCellBackgroundColor).Foreground(BombColor)
		}
	}

	return cellStyle.Render(content)
}

func getCell(v Value) (string, string) {
	switch v {
	case BLANK:
		return "", ""
	case ONE:
		return "1", "#fb4934"
	case TWO:
		return "2", "#b8bb26"
	case THREE:
		return "3", "#fabd2f"
	case FOUR:
		return "4", "#83a598"
	case FIVE:
		return "5", "#d3869b"
	case SIX:
		return "6", "#8ec07c"
	case SEVEN:
		return "7", "#fe8019"
	case EIGHT:
		return "8", "#b16286"
	case BOMB:
		return "💣", "#000000"
	default:
		panic("Should not happen")
	}
}
