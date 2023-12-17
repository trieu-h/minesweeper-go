package main

import (
	"math/rand"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type State int
type Value int
type Status int

const (
	GRID_WIDTH  = 13
	GRID_HEIGHT = 10
	CELL_WIDTH  = 3
	CELL_HEIGHT = 1
	BOMB_COUNT  = (GRID_WIDTH * GRID_HEIGHT) / 4
)

const (
	UNOPENED State = 0
	OPENED   State = 1
	FLAGGED  State = 2
)

const (
	NEW  Status = 0
	PLAY Status = 1
	WIN  Status = 2
	LOSE Status = 3
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

type Coords struct {
	x int
	y int
}

type Cell struct {
	state  State
	value  Value
	pos    Pos
	coords Coords
}

type Model struct {
	cells        [][]Cell
	activeCell   *Cell
	bombCells    []*Cell
	status       Status
	bombCounter  int
	timer        int
	timerStarted bool

	termHeight int
	termWidth  int
}

func (m *Model) Init() tea.Cmd {
	return nil
}

type TimerMsg struct{}

func timerCmd() tea.Msg {
	for {
		time.Sleep(1 * time.Second)
		return TimerMsg{}
	}
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
	case TimerMsg:
		if m.status == LOSE || m.status == NEW {
			return m, nil
		} else {
			m.timer++
			return m, timerCmd
		}

	case tea.WindowSizeMsg:
		m.termHeight, m.termWidth = msg.Height, msg.Width

	case tea.MouseMsg:
		switch msg.Button {
		case tea.MouseButtonLeft:
			switch msg.Action {
			case tea.MouseActionPress:
				for x := 0; x < GRID_HEIGHT; x++ {
					for y := 0; y < GRID_WIDTH; y++ {
						xCord := m.cells[x][y].coords.x
						yCord := m.cells[x][y].coords.y
						if xCord <= msg.X && msg.X <= xCord+CELL_WIDTH && yCord <= msg.Y && msg.Y <= yCord+CELL_HEIGHT {
							m.activeCell = &m.cells[x][y]
							return m.tryOpenActiveCell()
						}
					}
				}
			}
		case tea.MouseButtonRight:
			switch msg.Action {
			case tea.MouseActionPress:
				for x := 0; x < GRID_HEIGHT; x++ {
					for y := 0; y < GRID_WIDTH; y++ {
						xCord := m.cells[x][y].coords.x
						yCord := m.cells[x][y].coords.y
						if xCord <= msg.X && msg.X <= xCord+CELL_WIDTH && yCord <= msg.Y && msg.Y <= yCord+CELL_HEIGHT {
							m.activeCell = &m.cells[x][y]
							return m.flagCell()
						}
					}
				}
			}
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "left", "a":
			newY := min(m.activeCell.pos.y-1, 0)
			curX := m.activeCell.pos.x

			m.activeCell = &m.cells[curX][newY]
			return m, nil

		case "right", "d":
			newY := max(m.activeCell.pos.y+1, GRID_WIDTH-1)
			curX := m.activeCell.pos.x

			m.activeCell = &m.cells[curX][newY]
			return m, nil

		case "up", "w":
			newX := min(m.activeCell.pos.x-1, 0)
			curY := m.activeCell.pos.y

			m.activeCell = &m.cells[newX][curY]
			return m, nil

		case "down", "s":
			newX := max(m.activeCell.pos.x+1, GRID_HEIGHT-1)
			curY := m.activeCell.pos.y

			m.activeCell = &m.cells[newX][curY]
			return m, nil

		case "enter":
			return m.tryOpenActiveCell()

		case " ":
			return m.flagCell()

		case "r":
			m.new()
			return m, nil
		}
	}

	return m, nil
}

func (m *Model) flagCell() (tea.Model, tea.Cmd) {
	if m.status == LOSE {
		return m, nil
	}

	if m.activeCell.state == UNOPENED {
		m.activeCell.state = FLAGGED
		m.bombCounter--
	} else if m.activeCell.state == FLAGGED {
		m.activeCell.state = UNOPENED
		m.bombCounter++
	}
	return m, nil
}

func (m *Model) tryOpenActiveCell() (tea.Model, tea.Cmd) {
	if m.status == LOSE {
		return m, nil
	}

	if m.activeCell.state == UNOPENED {
		m.revealCell(m.activeCell)
	} else if m.activeCell.state == OPENED {
		x := m.activeCell.pos.x
		y := m.activeCell.pos.y

		for xc := x - 1; xc <= x+1; xc++ {
			if xc < 0 || xc > GRID_HEIGHT-1 {
				continue
			}

			for yc := y - 1; yc <= y+1; yc++ {
				if yc < 0 || yc > GRID_WIDTH-1 {
					continue
				}

				m.revealCell(&m.cells[xc][yc])
			}
		}
	}

	if !m.timerStarted {
		m.timerStarted = true
		m.status = PLAY
		return m, timerCmd
	} else {
		return m, nil
	}
}

func (m *Model) revealCell(cell *Cell) {
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

					m.revealCell(&m.cells[xc][yc])
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

func (m *Model) new() {
	var cells [][]Cell

	// Init cells
	for x := 0; x < GRID_HEIGHT; x++ {
		var row []Cell

		for y := 0; y < GRID_WIDTH; y++ {
			cell := Cell{
				state: UNOPENED,
				pos:   Pos{x, y},
				value: BLANK,
			}
			row = append(row, cell)
		}
		cells = append(cells, row)
	}

	// Start placing bombs
	// TODO: Maybe we should start placing the bombs after the player first move
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

	m.cells = cells
	m.activeCell = &cells[GRID_HEIGHT/2][GRID_WIDTH/2]
	m.status = NEW
	m.bombCells = bombCells
	m.bombCounter = BOMB_COUNT
	m.timer = 0
	m.timerStarted = false
}
