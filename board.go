package main

import (
	"github.com/charmbracelet/lipgloss"
)

const (
	PanelBorderColor        = lipgloss.Color("#ffffff")
	TextColor               = lipgloss.Color("#ffffff")
	ActiveBorderColor       = lipgloss.Color("#00ff00")
	CellBackgroundColor     = lipgloss.Color("#bdb2ff")
	BombCellBackgroundColor = lipgloss.Color("#ff0000")
	BombColor               = lipgloss.Color("#000000")
	FlagColor               = lipgloss.Color("#ff0000")
)

var (
	textStyle = lipgloss.NewStyle().
			Foreground(TextColor)

	panelStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(PanelBorderColor).
			Align(lipgloss.Center)

	cellStyle = lipgloss.NewStyle().
			Width(7).
			Padding(1, 0).
			Align(lipgloss.Center).
			BorderStyle(lipgloss.ThickBorder()).
			Bold(true)
)

func (m *Model) View() string {
	tutorial := lipgloss.JoinVertical(
		lipgloss.Left,
		renderTutorialText("Up arrow", "Move up"),
		renderTutorialText("Down arrow", "Move down"),
		renderTutorialText("Left arrow", "Move left"),
		renderTutorialText("Right arrow", "Move right"),
		renderTutorialText("Enter", "Reveal a cell"),
		renderTutorialText("Space", "Flag a cell"),
		renderTutorialText("R", "Restart"),
	)

	leftPanel := panelStyle.Copy().Padding(1, 2).MarginRight(3).Render(tutorial)

	var allRows []string

	for x := 0; x < GRID_HEIGHT; x++ {
		var thisRow []string

		for y := 0; y < GRID_WIDTH; y++ {
			thisRow = append(thisRow, m.renderCell(x, y))
		}

		allRows = append(allRows, lipgloss.JoinHorizontal(0, thisRow...))
	}

	board := lipgloss.JoinVertical(0, allRows...)

	rightPanel := panelStyle.Copy().Width(lipgloss.Width(board)).Render(board)

	gui := lipgloss.JoinHorizontal(0, leftPanel, rightPanel)

	return lipgloss.Place(m.termWidth, m.termHeight, 0.6, lipgloss.Center, gui)
}

func renderTutorialText(key string, instruction string) string {
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		textStyle.Copy().Bold(true).Render(key),
		textStyle.Copy().Render(" - "),
		textStyle.Copy().Italic(true).Render(instruction),
	)
}

func (m *Model) renderCell(x int, y int) string {
	var cell = m.cells[x][y]
	var cellStyle = cellStyle.Copy()
	var string = ""

	if &m.cells[x][y] == m.activeCell {
		cellStyle.BorderForeground(ActiveBorderColor)
	}

	if cell.state == UNOPENED {
		cellStyle.Background(CellBackgroundColor)
	}

	if cell.state == FLAGGED {
		string = "🚩"
		cellStyle.UnsetBackground().Foreground(FlagColor)
	}

	if cell.state == OPENED {
		string = convertValueToText(cell.value)

		if cell.value == BOMB {
			cellStyle.Background(BombCellBackgroundColor).Foreground(BombColor)
		} else {
			cellStyle.UnsetBackground()
		}
	}

	return cellStyle.Render(string)
}

func convertValueToText(v Value) string {
	switch v {
	case BLANK:
		return ""
	case ONE:
		return "1"
	case TWO:
		return "2"
	case THREE:
		return "3"
	case FOUR:
		return "4"
	case FIVE:
		return "5"
	case SIX:
		return "6"
	case SEVEN:
		return "7"
	case EIGHT:
		return "8"
	case BOMB:
		return "💣"
	default:
		panic("Should not happen")
	}
}
