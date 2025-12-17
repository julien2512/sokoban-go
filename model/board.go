package model

type cellType int

const (
	CellTypeNone cellType = iota
	CellTypeGoal
	CellTypeWall
)

type Cell struct {
	TypeOf cellType
	HasBox bool
	IsFree bool
}

type Board struct {
	Width, Height int
	Cells         []Cell
	LastMove      *LastMove
	Player        *Player
}

// NewBoard - Creates a board (map data encoding: Player "@", Box "$", Goal ".", Wall "#", Goal+Player "+", Goal+Box "*")
func NewBoard(mapData string, boardWidth, boardHeight int) *Board {
	b := Board{}

	b.Width = boardWidth
	b.Height = boardHeight

	b.Cells = make([]Cell, b.Width*b.Height)

	for y := 0; y < b.Height; y++ {
		for x := 0; x < b.Width; x++ {
			code := string(mapData[(y*b.Width)+x])
			cell := Cell{}
			switch code {
			case "@":
				b.Player = NewPlayer(x, y)
			case "$":
				cell.HasBox = true
			case ".":
				cell.TypeOf = CellTypeGoal
			case "#":
				cell.TypeOf = CellTypeWall
			case "+":
				cell.TypeOf = CellTypeGoal
				b.Player = NewPlayer(x, y)
			case "*":
				cell.TypeOf = CellTypeGoal
				cell.HasBox = true
			}
			b.Cells[(y*b.Width)+x] = cell
		}
	}

	return &b
}

func (b *Board) ResetFreeSpace() {
	for i :=0;i<len(b.Cells);i++ {
		b.Cells[i].IsFree = false
	}
}

// Private Checkup every Free Space from position
func (b *Board) _CheckEveryFreeSpace(x, y int) {
	c := b.Get(x,y)
	if (c.TypeOf == CellTypeWall || c.HasBox || c.IsFree) {
		return
	}
	c.IsFree = true
	b._CheckEveryFreeSpace(x-1,y)
	b._CheckEveryFreeSpace(x+1,y)
	b._CheckEveryFreeSpace(x,y-1)
	b._CheckEveryFreeSpace(x,y+1)
}

// Checkup every Free Space from player position
func (b *Board) CheckEveryFreeSpaceFromPlayer() {
	b.ResetFreeSpace()

	b._CheckEveryFreeSpace(b.Player.X,b.Player.Y)
}

// Get - Returns the cell at the given location
func (b *Board) Get(x, y int) *Cell {
	return &b.Cells[(y*b.Width)+x]
}

// IsComplete - Returns true if every goal cell on the board has a box
func (b *Board) IsComplete() bool {
	for _, cell := range b.Cells {
		if cell.TypeOf == CellTypeGoal && !cell.HasBox {
			return false
		}
	}
	return true
}
