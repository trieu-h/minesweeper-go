package main

import (
	"strconv"

	"github.com/charmbracelet/lipgloss"
)

func (m *Model) View() string {
	tutorial := lipgloss.JoinVertical(
		lipgloss.Left,
		"Up arrow - Move up",
		"Down arrow - Move down",
		"Left arrow - Move left",
		"Right arrow - Move right",
		"Enter - Reveal a cell",
		"Space - Flag a cell",
		"R - Restart",
	)

	leftPanel := lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.NormalBorder()).
		Align(lipgloss.Center).
		Render(tutorial)

	var allRows []string

	for x := 0; x < GRID_WIDTH; x++ {
		var thisRow []string

		for y := 0; y < GRID_HEIGHT; y++ {
			thisRow = append(thisRow, m.renderCell(x, y))
		}

		allRows = append(allRows, lipgloss.JoinHorizontal(0, thisRow...))
	}

	board := lipgloss.JoinVertical(0, allRows...)

	rightPanel := lipgloss.NewStyle().
		Width(lipgloss.Width(board)).
		Align(lipgloss.Center).
		Border(lipgloss.NormalBorder()).
		Render(board)

	game := lipgloss.JoinHorizontal(0, leftPanel, rightPanel)

	return lipgloss.Place(m.termWidth, m.termHeight, 0.6, lipgloss.Center, game)
}

func (m *Model) renderCell(x int, y int) string {
	var cell = m.cells[x][y]

	var style = lipgloss.NewStyle().
		Width(7).
		Padding(1, 0).
		Align(lipgloss.Center).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("63"))

	if &m.cells[x][y] == m.activeCell {
		style.BorderForeground(lipgloss.Color("50"))
	}

	var s string
	if cell.state == UNOPENED {
		style.Background(lipgloss.Color("63"))
		s = style.Render()
	} else if cell.state == FLAGGED {
		style.UnsetBackground()
		s = style.Render("Flagged")
	} else if cell.state == OPENED {
		style.UnsetBackground()

		if cell.value == BOMB {
			if m.status == LOSE {
				style.Background(lipgloss.Color("#FF0000"))
			}

			s = style.Render("Bomb!")
		} else if cell.value == BLANK {
			s = style.Render("")
		} else {
			s = style.Render(strconv.Itoa(int(cell.value)))
		}
	}

	return s
}
