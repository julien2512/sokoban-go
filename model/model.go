package model

type state int

const (
	StatePlaying state = iota
	StateLevelComplete
	StateGameComplete
)

const (
	LayerPlay = iota
	LayerBox
	LayerPath
)

type Model struct {
	LM              *LevelManager
	Board           *Board
	State           state
	TickAccumulator int
	Layer           int
	BoxLayer	int
}

// NewModel - Creates a model
func NewModel() *Model {
	m := Model{
		LM: NewLevelManager(false),
		Layer: LayerPlay,
	}

	return &m
}

// Update - Updates the model's current state (called once per main game loop iteration)
func (m *Model) Update() {
	m.TickAccumulator++
	if m.TickAccumulator > 20 {
		m.TickAccumulator = 0
	}
}
