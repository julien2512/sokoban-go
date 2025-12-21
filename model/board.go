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
	CanMoveDown bool
	CanMoveUp bool
	CanMoveLeft bool
	CanMoveRight bool
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

func (b *Board) _ResetCanBoxMove() {
	for i :=0;i<len(b.Cells);i++ {
		b.Cells[i].CanMoveLeft = false
		b.Cells[i].CanMoveRight = false
		b.Cells[i].CanMoveUp = false
		b.Cells[i].CanMoveDown = false
	}
}

func (b *Board) _CheckOneBoxMove(x,y int) {
	c := b.Get(x,y)

	if (!c.HasBox) {
		return
	}

	cup := b.Get(x,y-1)
	cdown := b.Get(x,y+1)
	if (cup.IsFree && cdown.TypeOf != CellTypeWall && !cdown.HasBox) {
		c.CanMoveDown = true
	}
	if (cdown.IsFree && cup.TypeOf != CellTypeWall && !cup.HasBox) {
		c.CanMoveUp = true		
	}

	cleft := b.Get(x-1,y)
	cright := b.Get(x+1,y)
	if (cleft.IsFree && cright.TypeOf != CellTypeWall && !cright.HasBox) {
		c.CanMoveRight = true
	}
	if (cright.IsFree && cleft.TypeOf != CellTypeWall && !cleft.HasBox) {
		c.CanMoveLeft = true		
	}

}

func (b *Board) _CheckEveryBoxMove() {
	for i :=0;i<len(b.Cells);i++ {
		y := i/b.Width
		x := i%b.Width
		b._CheckOneBoxMove(x,y)
	}
}

func (b *Board) ResetFreeSpace() {
	for i :=0;i<len(b.Cells);i++ {
		b.Cells[i].IsFree = false
	}
	b._ResetCanBoxMove()
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
	b._CheckEveryBoxMove()
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
