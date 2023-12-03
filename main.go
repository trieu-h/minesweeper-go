package main

import (
	"fmt"
	"os"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type State int
type Value int

const (
	GRID_WIDTH  = 4
	GRID_HEIGHT = 4
	CELL_WIDTH  = 5
	CELL_HEIGHT = CELL_WIDTH / 2
)

const (
	UNOPENED        State = 0
	NUMBERED        State = 1
	BLANK           State = 2
	FLAGGED         State = 3
	QUESTION_MARKED State = 4
)

const (
	ONE   Value = 1
	TWO   Value = 2
	THREE Value = 3
	FOUR  Value = 4
	FIVE  Value = 5
	SIX   Value = 6
	SEVEN Value = 7
	EIGHT Value = 8
	BOMB  Value = 9
)

type Pos struct {
	x int
	y int
}

type Cell struct {
	state State
	value Value
	pos   Pos
}

type Model struct {
	cells      [][]Cell
	activeCell *Cell
}

func (m Model) Init() tea.Cmd {
	return nil
}

func min(a int, b int) int {
	if a <= b {
		return b
	}
	return a
}

func max(a int, b int) int {
	if a >= b {
		return b
	}
	return a
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "left":
			newY := min(m.activeCell.pos.y-1, 0)
			curX := m.activeCell.pos.x

			m.activeCell = &m.cells[curX][newY]

		case "right":
			newY := max(m.activeCell.pos.y+1, GRID_WIDTH-1)
			curX := m.activeCell.pos.x

			m.activeCell = &m.cells[curX][newY]

		case "up":
			newX := min(m.activeCell.pos.x-1, 0)
			curY := m.activeCell.pos.y

			m.activeCell = &m.cells[newX][curY]

		case "down":
			newX := max(m.activeCell.pos.x+1, GRID_HEIGHT-1)
			curY := m.activeCell.pos.y

			m.activeCell = &m.cells[newX][curY]

		case "enter":
			if m.activeCell.state == UNOPENED {
				if !isBomb(m.activeCell) {
					m.activeCell.state = NUMBERED
				}
			}

		default:
			panic("Not handled yet")
		}
	}

	return m, nil
}

func isBomb(c *Cell) bool {
	if 1 <= c.value && c.value <= 8 {
		return false
	}
	return true
}

func (m *Model) View() string {
	var allRows []string

	for x := 0; x < GRID_WIDTH; x++ {
		var thisRow []string

		for y := 0; y < GRID_HEIGHT; y++ {
			thisRow = append(thisRow, renderCell(m, x, y))
		}

		allRows = append(allRows, lipgloss.JoinHorizontal(0, thisRow...))
	}

	return lipgloss.JoinVertical(0, allRows...)
}

func renderCell(m *Model, x int, y int) string {
	var cell = &m.cells[x][y]

	var style = lipgloss.NewStyle().
		Width(CELL_WIDTH).
		Height(CELL_HEIGHT).
		UnsetPadding().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("63"))

	if &m.cells[x][y] == m.activeCell {
		style.BorderForeground(lipgloss.Color("50"))
	}

	var s string
	if cell.state == UNOPENED {
		style.Background(lipgloss.Color("63"))
		s = style.Render()
	} else if cell.state == NUMBERED {
		style.UnsetBackground()
		s = style.Render(strconv.Itoa(int(cell.value)))
	}

	return s
}

func initState(w int, h int) *Model {
	var cells [][]Cell

	for x := 0; x < w; x++ {
		var row []Cell

		for y := 0; y < h; y++ {
			cell := Cell{
				state: UNOPENED,
				value: ONE,
				pos:   Pos{x: x, y: y},
			}
			row = append(row, cell)
		}
		cells = append(cells, row)
	}

	return &Model{
		cells:      cells,
		activeCell: &cells[1][1],
	}
}

func main() {
	p := tea.NewProgram(initState(GRID_WIDTH, GRID_HEIGHT))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
