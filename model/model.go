package model

import (
	"time"
)

type state int

const (
	StatePlaying state = iota
	StateLevelComplete
	StateGameComplete
	StateAutoplay
)

type Model struct {
	LM             *LevelManager
	Board          *Board
	Boards		map[string]*Board
	LastMove       *LastMove
	State           state
	TickAccumulator int
	Moves		int
	BestMoves	int
	SolveDuration	time.Duration
}

// NewModel - Creates a model
func NewModel() *Model {
	m := Model{
		LM: NewLevelManager(false),
		Boards: make(map[string]*Board)	}

	return &m
}

// Update - Updates the model's current state (called once per main game loop iteration)
func (m *Model) Update() {
	m.TickAccumulator++
	if m.TickAccumulator > 20 {
		m.TickAccumulator = 0
	}
}
