package blocks

import (
	"encoding/json"
	"math/rand"
	"time"
)

const (
	numRows      = 22
	numCols      = 10
	overflow     = 2
	tickerMillis = 17
	levelSecs    = 30
	speedBump    = 1.1
)

// Meta ...
type Meta struct {
	Rows     int `json:"rows"`
	Cols     int `json:"cols"`
	Overflow int `json:"overflow"`
}

// State ...
type State struct {
	Tiles    []int `json:"tiles"`
	Score    int   `json:"score"`
	Level    int   `json:"level"`
	GameOver bool  `json:"gameOver"`
}

const fastSpeed = 1
const slowSpeed = 50

// Model ...
type Model struct {
	commands   chan string
	updates    chan []byte
	done       chan bool
	curRow     int
	curCol     int
	board      *Arr
	curPiece   *Arr
	levelSpeed int
	curSpeed   int
	elapsed    int
	level      int
	levelStart int64
	score      int
	gameOver   bool
}

// NewModel ...
func NewModel() *Model {
	model := Model{
		commands:   make(chan string, 64),
		updates:    make(chan []byte),
		done:       make(chan bool, 1),
		level:      1,
		score:      0,
		levelStart: time.Now().Unix(),
	}
	model.board = NewArr(numRows, numCols)
	return &model
}

// Metadata ...
func Metadata() *Meta {
	return &Meta{
		Rows:     numRows,
		Cols:     numCols,
		Overflow: overflow,
	}
}

// Run ...
func (m *Model) Run() {
	m.levelSpeed = slowSpeed
	m.curSpeed = m.levelSpeed
	m.spawn()
	m.down()
	rand.Seed(time.Now().UnixNano())
	ticker := time.NewTicker(tickerMillis * time.Millisecond)
	for {
		select {
		case <-m.done:
			return
		case <-ticker.C:
			m.elapsed++
			if m.elapsed > m.curSpeed {
				m.elapsed = 0
				m.down()
			}
			if (time.Now().Unix() - m.levelStart) >= levelSecs {
				m.level++
				m.levelStart = time.Now().Unix()
				newSpeed := float64(m.levelSpeed) / speedBump
				m.levelSpeed = int(newSpeed)
			}
			state := State{
				Tiles:    m.board.buf,
				Level:    m.level,
				Score:    m.score,
				GameOver: m.gameOver,
			}
			b, _ := json.Marshal(state)
			m.updates <- b
		case cmd := <-m.commands:
			m.handleCommand(cmd)
		}
	}
}

func (m *Model) left() {
	clone := m.board.Clone()
	clone.Remove(m.curPiece, m.curRow, m.curCol)
	if clone.CanPlace(m.curPiece, m.curRow, m.curCol-1) {
		m.curCol--
		clone.Place(m.curPiece, m.curRow, m.curCol)
		m.board = clone
	}
}

func (m *Model) right() {
	clone := m.board.Clone()
	clone.Remove(m.curPiece, m.curRow, m.curCol)
	if clone.CanPlace(m.curPiece, m.curRow, m.curCol+1) {
		m.curCol++
		clone.Place(m.curPiece, m.curRow, m.curCol)
		m.board = clone
	}
}

func (m *Model) up() {
	clone := m.board.Clone()
	clone.Remove(m.curPiece, m.curRow, m.curCol)
	rotated := m.curPiece.RotateCounterClockwise()
	if clone.CanPlace(rotated, m.curRow, m.curCol) {
		clone.Place(rotated, m.curRow, m.curCol)
		m.curPiece = rotated
		m.board = clone
	}
}

func (m *Model) down() {
	clone := m.board.Clone()
	clone.Remove(m.curPiece, m.curRow, m.curCol)
	if clone.CanPlace(m.curPiece, m.curRow+1, m.curCol) {
		m.curRow++
		clone.Place(m.curPiece, m.curRow, m.curCol)
		m.board = clone
	} else {
		scoreMultipliers := [...]int{0, 40, 100, 300, 1200}
		removed := m.board.RemoveFullRows()
		m.score += scoreMultipliers[removed] * m.level
		m.spawn()
		m.curSpeed = m.levelSpeed
	}
}

func (m *Model) space() {
	m.curSpeed = fastSpeed
}

func (m *Model) handleCommand(cmd string) {
	switch cmd {
	case "left":
		m.left()
	case "right":
		m.right()
	case "up":
		m.up()
	case "down":
		m.down()
	case "space":
		m.space()
	}

}

func (m *Model) spawn() {
	which := rand.Intn(len(pieces))
	m.curPiece = pieces[which].Clone()
	m.curCol = 4
	m.curRow = 0
	if !m.board.CanPlace(m.curPiece, m.curRow, m.curCol) {
		m.gameOver = true
	}
}
