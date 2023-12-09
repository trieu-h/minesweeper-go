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
type Status int

const (
	GRID_WIDTH  = 8
	GRID_HEIGHT = 8
	BOMB_COUNT  = 8
)

const (
	UNOPENED State = 0
	OPENED   State = 1
	FLAGGED  State = 2
)

const (
	PLAYING Status = 0
	WIN     Status = 1
	LOSE    Status = 2
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
	bombCells  []*Cell
	status     Status

	termHeight int
	termWidth  int
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
	case tea.WindowSizeMsg:
		m.termHeight = msg.Height
		m.termWidth = msg.Width

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
			m.revealOneCell(m.activeCell)

		case " ":
			if m.activeCell.state == UNOPENED {
				m.activeCell.state = FLAGGED
			} else if m.activeCell.state == FLAGGED {
				m.activeCell.state = UNOPENED
			}
		}
	}

	return m, nil
}

func (m *Model) revealOneCell(cell *Cell) {
	var x = cell.pos.x
	var y = cell.pos.y

	if cell.state == UNOPENED {
		cell.state = OPENED

		if cell.value == BLANK {
			for xc := x - 1; xc <= x+1; xc++ {
				if xc < 0 || xc > GRID_HEIGHT-1 {
					continue
				}

				for yc := y - 1; yc <= y+1; yc++ {
					if yc < 0 || yc > GRID_WIDTH-1 {
						continue
					}

					m.revealOneCell(&m.cells[xc][yc])
				}
			}
		}

		if cell.value == BOMB {
			m.status = LOSE

			for _, bomb := range m.bombCells {
				bomb.state = OPENED
			}
		}
	}
}

func (m *Model) View() string {
	var allRows []string

	for x := 0; x < GRID_WIDTH; x++ {
		var thisRow []string

		for y := 0; y < GRID_HEIGHT; y++ {
			thisRow = append(thisRow, m.renderCell(x, y))
		}

		allRows = append(allRows, lipgloss.JoinHorizontal(0, thisRow...))
	}

	return lipgloss.JoinVertical(0, allRows...)
}

func (m *Model) renderCell(x int, y int) string {
	var cell = m.cells[x][y]

	var style = lipgloss.NewStyle().
		Width(7).
		PaddingTop(1).
		PaddingBottom(1).
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
	bombCells := []*Cell{}

	for bombCount > 0 {
		rx := rand.Intn(GRID_HEIGHT)
		ry := rand.Intn(GRID_WIDTH)

		if cells[rx][ry].value != BOMB {
			cells[rx][ry].value = BOMB
			bombCells = append(bombCells, &cells[rx][ry])
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

	fmt.Println(bombCells)
	return &Model{
		cells:      cells,
		activeCell: &cells[1][1],
		status:     PLAYING,
		bombCells:  bombCells,
	}
}

func main() {
	g := initGame()
	p := tea.NewProgram(g, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
