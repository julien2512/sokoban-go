package model

import (
	"fmt"
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
	Box int
}

type Box struct {
	X int
	Y int
	IsDead bool
	CanMoveDown bool
	CanMoveUp bool
	CanMoveLeft bool
	CanMoveRight bool
	ShallNotMoveDown bool
	ShallNotMoveUp bool
	ShallNotMoveLeft bool
	ShallNotMoveRight bool

	IsCheckedUp bool
	IsCheckedDown bool
	IsCheckedRight bool
	IsCheckedLeft bool
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
	Boxes         []Box
	Player        *Player

	BestDir direction.Direction
	BestX, BestY int
	BestLength int
}

// NewBoard - Creates a board (map data encoding: Player "@", Box "$", Goal ".", Wall "#", Goal+Player "+", Goal+Box "*")
func NewBoard(mapData string, boardWidth, boardHeight int) *Board {
	b := Board{}

	b.Width = boardWidth
	b.Height = boardHeight

	b.Cells = make([]Cell, b.Width*b.Height)

	box := 0
	b.Boxes = make([]Box, 0)

	for y := 0; y < b.Height; y++ {
		for x := 0; x < b.Width; x++ {
			code := string(mapData[(y*b.Width)+x])
			cell := Cell{}
			switch code {
			case "@":
				b.Player = NewPlayer(x, y)
			case "$":
				cell.HasBox = true
				cell.Box = box
				box++
				b.Boxes = append(b.Boxes,Box{X:x,Y:y})
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
				cell.Box = box
				box++
				b.Boxes = append(b.Boxes,Box{X:x,Y:y})
			}
			b.Cells[(y*b.Width)+x] = cell
		}
	}

	b.BestX = -1
	b.BestY = -1
	b.BestLength = 1000 // assume max length

	return &b
}

func (b *Board) Print() {
	for y :=0;y<b.Height;y++ {
		for x :=0;x<b.Width;x++ {
			c := b.Get(x,y)
			
			if c.HasBox && c.TypeOf == CellTypeGoal && b.Boxes[c.Box].X==x && b.Boxes[c.Box].Y==y { fmt.Print("*")
			} else if x == b.Player.X && y == b.Player.Y && c.TypeOf == CellTypeGoal { fmt.Print("+")
			} else if c.TypeOf == CellTypeWall { fmt.Print("#")
			} else if c.TypeOf == CellTypeGoal { fmt.Print(".")
			} else if c.HasBox && c.TypeOf != CellTypeGoal && b.Boxes[c.Box].X==x && b.Boxes[c.Box].Y==y { fmt.Print("$")
			} else if x == b.Player.X && y == b.Player.Y && c.TypeOf != CellTypeGoal { fmt.Print("@")
			} else { fmt.Print(" ") }
		}
		fmt.Println()
	}
}

func (b *Board) GetString() string {
	var str string
	for _, cell := range b.Cells {
		str = str + getCellString(cell)
	}
	return str
}

func (b *Board) Duplicate() *Board {
	d := &Board{}
	d.Width = b.Width
	d.Height = b.Height

	d.Cells = make([]Cell, b.Width*b.Height)
	d.Boxes = make([]Box, len(b.Boxes))

	for i, cell := range b.Cells {
		d.Cells[i].TypeOf = cell.TypeOf
		d.Cells[i].HasBox = cell.HasBox
		d.Cells[i].IsFree = cell.IsFree
		d.Cells[i].Box = cell.Box
	}
	for i, box := range b.Boxes {
		d.Boxes[i].X = box.X
		d.Boxes[i].Y = box.Y
		d.Boxes[i].IsDead = box.IsDead
		d.Boxes[i].CanMoveDown = box.CanMoveDown
		d.Boxes[i].CanMoveUp = box.CanMoveUp
		d.Boxes[i].CanMoveLeft = box.CanMoveLeft
		d.Boxes[i].CanMoveRight = box.CanMoveRight
		d.Boxes[i].ShallNotMoveUp = box.ShallNotMoveUp
		d.Boxes[i].ShallNotMoveDown = box.ShallNotMoveDown
		d.Boxes[i].ShallNotMoveLeft = box.ShallNotMoveLeft
		d.Boxes[i].ShallNotMoveRight = box.ShallNotMoveRight
		d.Boxes[i].IsCheckedDown = box.IsCheckedDown
		d.Boxes[i].IsCheckedUp = box.IsCheckedUp
		d.Boxes[i].IsCheckedLeft = box.IsCheckedLeft
		d.Boxes[i].IsCheckedRight = box.IsCheckedRight
	}

	d.BestDir = b.BestDir
	d.BestX = b.BestX
	d.BestY = b.BestY
	d.BestLength = b.BestLength

	d.Player = NewPlayer(b.Player.X,b.Player.Y)

	return d
}

func (b *Board) copyFrom(model *Board) {
	b.Width = model.Width
	b.Height = model.Height
	b.Player = NewPlayer(model.Player.X,model.Player.Y)
	b.Cells = make([]Cell, b.Width*b.Height)
	b.Boxes = make([]Box, len(model.Boxes))

	for i, cell := range model.Cells {
		b.Cells[i].TypeOf = cell.TypeOf
		b.Cells[i].HasBox = cell.HasBox
		b.Cells[i].IsFree = cell.IsFree
		b.Cells[i].Box = cell.Box
	}
	for i, box := range model.Boxes {
		b.Boxes[i].X = box.X
		b.Boxes[i].Y = box.Y
		b.Boxes[i].IsDead = box.IsDead
		b.Boxes[i].CanMoveDown = box.CanMoveDown
		b.Boxes[i].CanMoveUp = box.CanMoveUp
		b.Boxes[i].CanMoveLeft = box.CanMoveLeft
		b.Boxes[i].CanMoveRight = box.CanMoveRight
		b.Boxes[i].ShallNotMoveUp = box.ShallNotMoveUp
		b.Boxes[i].ShallNotMoveDown = box.ShallNotMoveDown
		b.Boxes[i].ShallNotMoveLeft = box.ShallNotMoveLeft
		b.Boxes[i].ShallNotMoveRight = box.ShallNotMoveRight
		b.Boxes[i].IsCheckedDown = box.IsCheckedDown
		b.Boxes[i].IsCheckedUp = box.IsCheckedUp
		b.Boxes[i].IsCheckedLeft = box.IsCheckedLeft
		b.Boxes[i].IsCheckedRight = box.IsCheckedRight
	}
	b.BestDir = model.BestDir
	b.BestX = model.BestX
	b.BestY = model.BestY
	b.BestLength = model.BestLength
}

func (b *Board) _ResetCanBoxMove() {
	for i :=0;i<len(b.Boxes);i++ {
		b.Boxes[i].IsDead = false
		b.Boxes[i].CanMoveLeft = false
		b.Boxes[i].CanMoveRight = false
		b.Boxes[i].CanMoveUp = false
		b.Boxes[i].CanMoveDown = false
		b.Boxes[i].ShallNotMoveLeft = true
		b.Boxes[i].ShallNotMoveRight = true
		b.Boxes[i].ShallNotMoveUp = true
		b.Boxes[i].ShallNotMoveDown = true
		b.Boxes[i].IsCheckedDown = false
		b.Boxes[i].IsCheckedUp = false
		b.Boxes[i].IsCheckedLeft = false
		b.Boxes[i].IsCheckedRight = false
	}
}

func (b *Board) _CheckOneBoxIsDead(x,y int) bool {
	c := b.Get(x,y)

	if !c.HasBox || c.TypeOf == CellTypeGoal {
		return false
	}
	box := &b.Boxes[c.Box]
	cup := b.Get(x,y-1)
	cdown := b.Get(x,y+1)
	cleft := b.Get(x-1,y)
	cright := b.Get(x+1,y)

	if cup.TypeOf == CellTypeWall && cleft.TypeOf == CellTypeWall { box.IsDead = true; return true }
	if cup.TypeOf == CellTypeWall && cright.TypeOf == CellTypeWall { box.IsDead = true; return true }
	if cdown.TypeOf == CellTypeWall && cleft.TypeOf == CellTypeWall { box.IsDead = true; return true }
	if cdown.TypeOf == CellTypeWall && cright.TypeOf == CellTypeWall { box.IsDead = true; return true }

	return false
}

func (b *Board) _CheckEveryBoxIsDead() bool {
	count := 0
	for i :=0;i<len(b.Cells);i++ {
		y := i/b.Width
		x := i%b.Width

		if b._CheckOneBoxIsDead(x,y) { count++ }
	}
	return count > 0
}

func (b *Board) _CheckOneBoxIsTrapByDirWall(x,y int, dirx, diry int) bool {
	c := b.Get(x,y)
	if !c.HasBox { return false }
	box := &b.Boxes[c.Box]
	if box.IsDead { return true }
	if c.TypeOf == CellTypeGoal { return false }
	GoalCount := 0
	BoxCount := 1
	cUp := b.Get(x+dirx,y+diry)
	if cUp.TypeOf != CellTypeWall { return false }

	xRight := x-1*diry
	yRight := y-1*dirx
	for {
		cRight := b.Get(xRight,yRight)
		if cRight.TypeOf == CellTypeWall { break; }
		if cRight.HasBox { BoxCount++ }
		if cRight.TypeOf == CellTypeGoal { GoalCount++ }

		cUp = b.Get(xRight+dirx,yRight+diry)
		if cUp.TypeOf != CellTypeWall { return false }
		xRight += -1*diry
		yRight += -1*dirx
	}
	xLeft := x+1*diry
	yLeft := y+1*dirx
	for {
		cLeft := b.Get(xLeft,yLeft)
		if cLeft.TypeOf == CellTypeWall { break; }
		if cLeft.HasBox { BoxCount++ }
		if cLeft.TypeOf == CellTypeGoal { GoalCount++ }

		cUp = b.Get(xLeft+dirx,yLeft+diry)
		if cUp.TypeOf != CellTypeWall { return false }
		xLeft += 1*diry
		yLeft += 1*dirx
	}
	
	if BoxCount > GoalCount {
		box.IsDead = true
		return true
	} else { return false }
}

func (b *Board) _CheckEveryBoxIsTrapByWall() bool {
	count := 0
	for i :=0;i<len(b.Boxes);i++ {
		box := &b.Boxes[i]
		y := box.Y
		x := box.X
		
		if b._CheckOneBoxIsTrapByDirWall(x,y,0,-1) { count++ }
		if b._CheckOneBoxIsTrapByDirWall(x,y,0,1) { count++ }
		if b._CheckOneBoxIsTrapByDirWall(x,y,-1,0) { count++ }
		if b._CheckOneBoxIsTrapByDirWall(x,y,1,0) { count++ }
	}	
	return count > 0
}

func (b *Board) _CheckEveryBoxIsTrap() bool {
	traped := false
	traped = traped || b._CheckEveryBoxIsTrapByWall()
	return traped
}

func (b *Board) _CheckOneBoxMove(x,y int,boards map[string]*Board) {
	c := b.Get(x,y)
	box := &b.Boxes[c.Box]
	
	if (!c.HasBox) {
		return
	}
	
	cup := b.Get(x,y-1)
	cdown := b.Get(x,y+1)
	if !box.IsCheckedDown {
		box.IsCheckedDown = true
		if (cup.IsFree && cdown.TypeOf != CellTypeWall && !cdown.HasBox) {
			box.CanMoveDown = true
			newBoard := b.Duplicate()
			newBoard.BestLength = 1000
			newBoard.BestX = -1
			newBoard.BestY = -1
			newBoard._MoveBox(x,y,direction.D,boards)
			if newBoard.GetGoodBoxMoveCount() > 0 || newBoard.BestLength == 0 {
				box.ShallNotMoveDown = false
				if newBoard.BestLength+1<b.BestLength {
					b.BestLength = newBoard.BestLength+1
					b.BestX = x
					b.BestY = y
					b.BestDir = direction.D
				}
			}
		}
	}
	if !box.IsCheckedUp {
		box.IsCheckedUp = true
		if (cdown.IsFree && cup.TypeOf != CellTypeWall && !cup.HasBox) {
			box.CanMoveUp = true
			newBoard := b.Duplicate()
			newBoard.BestLength = 1000
			newBoard.BestX = -1
			newBoard.BestY = -1
			newBoard._MoveBox(x,y,direction.U,boards)
			if newBoard.GetGoodBoxMoveCount() > 0 || newBoard.BestLength == 0 {
				box.ShallNotMoveUp = false
				if newBoard.BestLength+1<b.BestLength {
					b.BestLength = newBoard.BestLength+1
					b.BestX = x
					b.BestY = y
					b.BestDir = direction.U
				}
			}
		}
	}
	
	cleft := b.Get(x-1,y)
	cright := b.Get(x+1,y)
	if !box.IsCheckedRight {
		box.IsCheckedRight = true
		if (cleft.IsFree && cright.TypeOf != CellTypeWall && !cright.HasBox) {
			box.CanMoveRight = true
			newBoard := b.Duplicate()
			newBoard.BestLength = 1000
			newBoard.BestX = -1
			newBoard.BestY = -1
			newBoard._MoveBox(x,y,direction.R,boards)
			if newBoard.GetGoodBoxMoveCount() > 0 || newBoard.BestLength == 0 {
				box.ShallNotMoveRight = false
				if newBoard.BestLength+1<b.BestLength {
					b.BestLength = newBoard.BestLength+1
					b.BestX = x
					b.BestY = y
					b.BestDir = direction.R
				}
			}
		}
	}
	if !box.IsCheckedLeft {
		box.IsCheckedLeft = true
		if (cright.IsFree && cleft.TypeOf != CellTypeWall && !cleft.HasBox) {
			box.CanMoveLeft = true
			newBoard := b.Duplicate()
			newBoard.BestLength = 1000
			newBoard.BestX = -1
			newBoard.BestY = -1
			newBoard._MoveBox(x,y,direction.L,boards)
			if newBoard.GetGoodBoxMoveCount() > 0 || newBoard.BestLength == 0 { 
				box.ShallNotMoveLeft = false
				if newBoard.BestLength+1<b.BestLength {
					b.BestLength = newBoard.BestLength+1
					b.BestX = x
					b.BestY = y
					b.BestDir = direction.L
				}
			}
		}
	}
}

func (b *Board) _CheckEveryBoxMove(boards map[string]*Board) {
	for i :=0;i<len(b.Boxes);i++ {
		y := b.Boxes[i].Y
		x := b.Boxes[i].X
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
	b.ResetFreeSpace()
	b._CheckEveryFreeSpace(b.Player.X,b.Player.Y)

	boardName := b.GetString()

	if boards[boardName] == nil {
		boards[boardName] = b
		if b._CheckEveryBoxIsTrap() {
			b.BestLength = 999
		} else if b.IsComplete() {
			b.BestLength = 0
		} else {
			b._CheckEveryBoxMove(boards)
		}
	} else {
                b.copyFrom(boards[boardName])
		if b.BestLength!=999 && b.BestLength!=0 {
			//b._CheckEveryBoxMove(boards)  // there is an infinite loop inside too solve
		}
	}
}

// Checkup every Free Space from player position
func (b *Board) CheckEveryFreeSpaceFromPlayer(boards map[string]*Board) {
	X := b.Player.X
	Y := b.Player.Y
	b.BestX = -1
	b.BestY = -1
	b.BestLength = 1000
	b._CheckEveryFreeSpaceFromPlayer(boards)
	b.Player = NewPlayer(X,Y)
}

// assume it
func (b *Board) _MoveBox(x,y int, dir direction.Direction, boards map[string]*Board) {
	lastCell := b.Get(x,y)
	b.Player.X = x
	b.Player.Y = y
	var newCell *Cell
	box := &b.Boxes[lastCell.Box]
	if dir == direction.L {
		newCell = b.Get(x-1,y)
		box.X = x-1
	} else if dir == direction.R {
		newCell = b.Get(x+1,y)	
		box.X = x+1
	} else if dir == direction.U {
		newCell = b.Get(x,y-1)
		box.Y = y-1
	} else if dir == direction.D {
		newCell = b.Get(x,y+1)
		box.Y = y+1
	}
	lastCell.HasBox = false
	newCell.HasBox = true
	newCell.Box = lastCell.Box
	
	b._CheckEveryFreeSpaceFromPlayer(boards)
}

func (b *Board) GetBoxMoveCount() int {
	count := 0
	for _, box := range b.Boxes {
		if box.CanMoveLeft { count = count+1 }
		if box.CanMoveRight { count = count+1 }
		if box.CanMoveUp { count = count+1 }
		if box.CanMoveDown { count = count+1 }
	}
	return count
}

func (b *Board) GetGoodBoxMoveCount() int {
	count := 0
	for _, box := range b.Boxes {
		if box.CanMoveLeft && !box.ShallNotMoveLeft { count = count+1 }
		if box.CanMoveRight && !box.ShallNotMoveRight { count = count+1 }
		if box.CanMoveUp && !box.ShallNotMoveUp { count = count+1 }
		if box.CanMoveDown && !box.ShallNotMoveDown { count = count+1 }
	}
	return count
}

// Get - Returns the cell at the given location
func (b *Board) Get(x, y int) *Cell {
	return &b.Cells[(y*b.Width)+x]
}

// IsComplete - Returns true if every goal cell on the board has a box
func (b *Board) IsComplete() bool {
	for _, box := range b.Boxes {
		cell := b.Get(box.X,box.Y)
		if cell.TypeOf != CellTypeGoal {
			return false
		}
	}
	return true
}
