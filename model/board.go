package model

import (
	"github.com/TheInvader360/sokoban-go/direction"
)

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
	ShallNotMoveDown bool
	ShallNotMoveUp bool
	ShallNotMoveLeft bool
	ShallNotMoveRight bool
}

func getCharFromNum(i int32) string {
	return string(65+i)
}

func getCellString(c Cell) string {
	var i int32
	if c.TypeOf == CellTypeNone { i = 0 
	} else if c.TypeOf == CellTypeGoal { i = 1 
	} else if c.TypeOf == CellTypeWall { i = 2 }

	if c.HasBox { i = i + 3 }
	if c.IsFree { i = i + 6 }

	return getCharFromNum(i)
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

func (b *Board) GetString() string {
	var str string
	for _, cell := range b.Cells {
		str = str + getCellString(cell)
	}
	return str
}

func (b *Board) duplicate() *Board {
	d := &Board{}
	d.Width = b.Width
	d.Height = b.Height

	d.Cells = make([]Cell, b.Width*b.Height)

	for i, cell := range b.Cells {
		d.Cells[i].TypeOf = cell.TypeOf
		d.Cells[i].HasBox = cell.HasBox
		d.Cells[i].IsFree = cell.IsFree
		d.Cells[i].CanMoveDown = cell.CanMoveDown
		d.Cells[i].CanMoveUp = cell.CanMoveUp
		d.Cells[i].CanMoveLeft = cell.CanMoveLeft
		d.Cells[i].CanMoveRight = cell.CanMoveRight
		d.Cells[i].ShallNotMoveUp = cell.ShallNotMoveUp
		d.Cells[i].ShallNotMoveDown = cell.ShallNotMoveDown
		d.Cells[i].ShallNotMoveLeft = cell.ShallNotMoveLeft
		d.Cells[i].ShallNotMoveRight = cell.ShallNotMoveRight
	}

	d.Player = NewPlayer(b.Player.X,b.Player.Y)

	return d
}

func (b *Board) copyFrom(model *Board) {
	b.Width = model.Width
	b.Height = model.Height
	b.Player = NewPlayer(model.Player.X,model.Player.Y)
	b.Cells = make([]Cell, b.Width*b.Height)

	for i, cell := range model.Cells {
		b.Cells[i].TypeOf = cell.TypeOf
		b.Cells[i].HasBox = cell.HasBox
		b.Cells[i].IsFree = cell.IsFree
		b.Cells[i].CanMoveDown = cell.CanMoveDown
		b.Cells[i].CanMoveUp = cell.CanMoveUp
		b.Cells[i].CanMoveLeft = cell.CanMoveLeft
		b.Cells[i].CanMoveRight = cell.CanMoveRight
		b.Cells[i].ShallNotMoveUp = cell.ShallNotMoveUp
		b.Cells[i].ShallNotMoveDown = cell.ShallNotMoveDown
		b.Cells[i].ShallNotMoveLeft = cell.ShallNotMoveLeft
		b.Cells[i].ShallNotMoveRight = cell.ShallNotMoveRight
	} 
}

func (b *Board) _ResetCanBoxMove() {
	for i :=0;i<len(b.Cells);i++ {
		b.Cells[i].CanMoveLeft = false
		b.Cells[i].CanMoveRight = false
		b.Cells[i].CanMoveUp = false
		b.Cells[i].CanMoveDown = false
		b.Cells[i].ShallNotMoveLeft = true
		b.Cells[i].ShallNotMoveRight = true
		b.Cells[i].ShallNotMoveUp = true
		b.Cells[i].ShallNotMoveDown = true
	}
}

func (b *Board) _CheckOneBoxMove(x,y int,boards map[string]*Board) {
	c := b.Get(x,y)
	
	if (!c.HasBox) {
		return
	}
	
	cup := b.Get(x,y-1)
	cdown := b.Get(x,y+1)
	if (cup.IsFree && cdown.TypeOf != CellTypeWall && !cdown.HasBox) {
		c.CanMoveDown = true
		newBoard := b.duplicate()
		newBoard._MoveBox(x,y,direction.D,boards)
		if newBoard.GetGoodBoxMoveCount() > 0 || newBoard.IsComplete() { c.ShallNotMoveDown = false }
	}
	if (cdown.IsFree && cup.TypeOf != CellTypeWall && !cup.HasBox) {
		c.CanMoveUp = true
		newBoard := b.duplicate()
		newBoard._MoveBox(x,y,direction.U,boards)
		if newBoard.GetGoodBoxMoveCount() > 0 || newBoard.IsComplete() { c.ShallNotMoveUp = false }
	}
	
	cleft := b.Get(x-1,y)
	cright := b.Get(x+1,y)
	if (cleft.IsFree && cright.TypeOf != CellTypeWall && !cright.HasBox) {
		c.CanMoveRight = true
		newBoard := b.duplicate()
		newBoard._MoveBox(x,y,direction.R,boards)
		
		if newBoard.GetGoodBoxMoveCount() > 0 || newBoard.IsComplete() { c.ShallNotMoveRight = false }
	}
	if (cright.IsFree && cleft.TypeOf != CellTypeWall && !cleft.HasBox) {
		c.CanMoveLeft = true
		newBoard := b.duplicate()
		newBoard._MoveBox(x,y,direction.L,boards)
		if newBoard.GetGoodBoxMoveCount() > 0 || newBoard.IsComplete() { c.ShallNotMoveLeft = false }
	}
}

func (b *Board) _CheckEveryBoxMove(boards map[string]*Board) {
	for i :=0;i<len(b.Cells);i++ {
		y := i/b.Width
		x := i%b.Width
		b._CheckOneBoxMove(x,y,boards)
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

func (b *Board) _CheckEveryFreeSpaceFromPlayer(boards map[string]*Board) {
	boardName := b.GetString()

	if boards[boardName] == nil {
		boards[boardName] = b
		b.ResetFreeSpace()
		
		b._CheckEveryFreeSpace(b.Player.X,b.Player.Y)
		b._CheckEveryBoxMove(boards)
	} else { b.copyFrom(boards[boardName]) }
}

// Checkup every Free Space from player position
func (b *Board) CheckEveryFreeSpaceFromPlayer() {
	boards := make(map[string]*Board)

	b._CheckEveryFreeSpaceFromPlayer(boards)
}

// assume it
func (b *Board) _MoveBox(x,y int, dir direction.Direction, boards map[string]*Board) {
	b.Player.X = x
	b.Player.Y = y
	lastCell := b.Get(x,y)
	var newCell *Cell
	if dir == direction.L {
		newCell = b.Get(x-1,y)
	} else if dir == direction.R {
		newCell = b.Get(x+1,y)	
	} else if dir == direction.U {
		newCell = b.Get(x,y-1)
	} else if dir == direction.D {
		newCell = b.Get(x,y+1)
	}
	lastCell.HasBox = false
	newCell.HasBox = true
	b._CheckEveryFreeSpaceFromPlayer(boards)
}

func (b *Board) GetBoxMoveNumber(number int) (int, int , direction.Direction) {
	count := 0
	for i :=0;i<len(b.Cells);i++ {
		cell := &b.Cells[i]
		y := i/b.Width
		x := i%b.Width
		if cell.CanMoveLeft { if count == number { return x, y, direction.L } 
			              count = count+1 }
		if cell.CanMoveRight { if count == number { return x, y, direction.R }
                                       count = count+1 }
		if cell.CanMoveUp { if count == number { return x, y, direction.U }
                                    count = count+1 }
		if cell.CanMoveDown { if count == number { return x, y, direction.D }
                                      count = count+1 }
	}
	return 0,0,0
}

func (b *Board) GetBoxMoveCount() int {
	count := 0
	for _, cell := range b.Cells {
		if cell.CanMoveLeft { count = count+1 }
		if cell.CanMoveRight { count = count+1 }
		if cell.CanMoveUp { count = count+1 }
		if cell.CanMoveDown { count = count+1 }
	}
	return count
}

func (b *Board) GetGoodBoxMoveCount() int {
	count := 0
	for _, cell := range b.Cells {
		if cell.CanMoveLeft && !cell.ShallNotMoveLeft { count = count+1 }
		if cell.CanMoveRight && !cell.ShallNotMoveRight { count = count+1 }
		if cell.CanMoveUp && !cell.ShallNotMoveUp { count = count+1 }
		if cell.CanMoveDown && !cell.ShallNotMoveDown { count = count+1 }
	}
	return count
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
