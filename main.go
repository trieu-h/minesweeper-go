package main

import (
	"fmt"
	"math/rand"
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
	BOMB_COUNT  = 4
)

const (
	UNOPENED State = 0
	OPENED   State = 1
)

const (
	BLANK Value = 0
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
				m.activeCell.state = OPENED
			}

		default:
			fmt.Printf("%s key is not handled yet!\n", msg.String())
		}
	}

	return m, nil
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
	} else if cell.state == OPENED {
		style.UnsetBackground()

		if cell.value == BOMB {
			s = style.Render("Bomb!")
		} else {
			s = style.Render(strconv.Itoa(int(cell.value)))
		}
	}

	return s
}

func initGame() *Model {
	var cells [][]Cell

	// Init cells
	for x := 0; x < GRID_HEIGHT; x++ {
		var row []Cell

		for y := 0; y < GRID_WIDTH; y++ {
			cell := Cell{
				state: UNOPENED,
				pos:   Pos{x: x, y: y},
				value: BLANK,
			}
			row = append(row, cell)
		}
		cells = append(cells, row)
	}

	// Start placing bombs
	bombCount := BOMB_COUNT
	for bombCount > 0 {
		rx := rand.Intn(GRID_HEIGHT)
		ry := rand.Intn(GRID_WIDTH)

		if cells[rx][ry].value != BOMB {
			cells[rx][ry].value = BOMB
			bombCount = bombCount - 1
		}
	}

	// Calculate the value of each cells
	for x := 0; x < GRID_HEIGHT; x++ {
		for y := 0; y < GRID_WIDTH; y++ {
			if cells[x][y].value == BOMB {
				continue
			}

			numberOfBomb := 0
			for xc := x - 1; xc <= x+1; xc++ {
				// Check for out of bounds
				if xc < 0 || xc > GRID_HEIGHT-1 {
					continue
				}
				for yc := y - 1; yc <= y+1; yc++ {
					// Check for out of bounds
					if yc < 0 || yc > GRID_WIDTH-1 {
						continue
					}

					// Exclude current cell
					if xc == x && yc == y {
						continue
					}

					if cells[xc][yc].value == BOMB {
						numberOfBomb++
					}
				}
			}

			cells[x][y].value = Value(numberOfBomb)
		}
	}

	return &Model{
		cells:      cells,
		activeCell: &cells[1][1],
	}
}

func main() {
	g := initGame()
	p := tea.NewProgram(g)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
