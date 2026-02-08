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
	Dist int
}

type Position struct {
	X,Y int
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

	XYChecked map[Position]bool

	UpBoard *Board
	DownBoard *Board
	LeftBoard *Board
	RightBoard *Board
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

type BestPosition struct {
	BestDir direction.Direction
	BestX, BestY int
	BestLength int
}

type Board struct {
	Width, Height int
	Cells         []Cell
	Boxes         []Box
	Player        *Player

	BestPositions map[Position]*BestPosition
}

func (b *Board) GetBestPosition() *BestPosition {
	Pos := Position{X:b.Player.X,Y:b.Player.Y}
	return b.BestPositions[Pos]
}

// NewBoard - Creates a board (map data encoding: Player "@", Box "$", Goal ".", Wall "#", Goal+Player "+", Goal+Box "*")
func NewBoard(mapData string, boardWidth, boardHeight int) *Board {
	b := Board{}

	b.Width = boardWidth
	b.Height = boardHeight

	b.Cells = make([]Cell, b.Width*b.Height)

	box := 0
	b.Boxes = make([]Box, 0)

	b.BestPositions = make(map[Position]*BestPosition)

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
				b.Boxes = append(b.Boxes,Box{X:x,Y:y,XYChecked:make(map[Position]bool)})
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
				b.Boxes = append(b.Boxes,Box{X:x,Y:y,XYChecked:make(map[Position]bool)})
			}
			b.Cells[(y*b.Width)+x] = cell
		}
	}

	b._ResetCanBoxMove()

	// assume max length
	b.BestPositions[Position{X:b.Player.X,Y:b.Player.Y}] = &BestPosition{BestLength:1000,BestX:-1,BestY:-1}

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
	d.BestPositions = make(map[Position]*BestPosition)

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
		d.Boxes[i].XYChecked = make(map[Position]bool)
	}

	d.Player = NewPlayer(b.Player.X,b.Player.Y)

	return d
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

// Assume x,y got a box
func (b *Board) _CheckOneBoxMove(x,y int,boards map[string]*Board) {
	c := b.Get(x,y)
	box := &(b.Boxes[c.Box])
	
	fromx := b.Player.X
	fromy := b.Player.Y
	from := Position{X:fromx,Y:fromy}
	to := Position{X:x,Y:y}
	var XYCheckedAlready bool
	if box.XYChecked[from] {  // assume every direction is checked for that box
		XYCheckedAlready = true
	} 
	

	cup := b.Get(x,y-1)
	cdown := b.Get(x,y+1)
	if !box.IsCheckedDown {
		box.IsCheckedDown = true
		if (cup.IsFree && cdown.TypeOf != CellTypeWall && !cdown.HasBox) {
			box.CanMoveDown = true
			newBoard:= b.MoveBoxAndCheck(x,y,direction.D,boards)
			if newBoard.GetGoodBoxMoveCount() > 0 || newBoard.BestPositions[to].BestLength == 0 {
				box.ShallNotMoveDown = false
				b.CheckEveryDist(fromx,fromy)
				if newBoard.BestPositions[to].BestLength+1+cup.Dist<b.BestPositions[from].BestLength {
					b.BestPositions[from].BestLength = newBoard.BestPositions[to].BestLength+1+cup.Dist
					b.BestPositions[from].BestX = x
					b.BestPositions[from].BestY = y
					b.BestPositions[from].BestDir = direction.D
				}
			}
		}
	} else if !XYCheckedAlready && box.CanMoveDown && !box.ShallNotMoveDown {
		newBoard:= b.GetOldMoveBox(x,y,direction.D,boards)
		b.CheckEveryDist(fromx,fromy)
		if newBoard!= nil && newBoard.BestPositions[to].BestLength+1+cup.Dist<b.BestPositions[from].BestLength {
			b.BestPositions[from].BestLength = newBoard.BestPositions[to].BestLength+1+cup.Dist
			b.BestPositions[from].BestX = x
			b.BestPositions[from].BestY = y
			b.BestPositions[from].BestDir = direction.D
		}
	}
	if !box.IsCheckedUp {
		box.IsCheckedUp = true
		if (cdown.IsFree && cup.TypeOf != CellTypeWall && !cup.HasBox) {
			box.CanMoveUp = true
			newBoard:= b.MoveBoxAndCheck(x,y,direction.U,boards)
			if newBoard.GetGoodBoxMoveCount() > 0 || newBoard.BestPositions[to].BestLength == 0 {
				box.ShallNotMoveUp = false
				b.CheckEveryDist(fromx,fromy)
				if newBoard.BestPositions[to].BestLength+1+cdown.Dist<b.BestPositions[from].BestLength {
					b.BestPositions[from].BestLength = newBoard.BestPositions[to].BestLength+1+cdown.Dist
					b.BestPositions[from].BestX = x
					b.BestPositions[from].BestY = y
					b.BestPositions[from].BestDir = direction.U
				}
			}
		}
	} else if !XYCheckedAlready && box.CanMoveUp && !box.ShallNotMoveUp {
		newBoard:= b.GetOldMoveBox(x,y,direction.U,boards)
		b.CheckEveryDist(fromx,fromy)
		if newBoard!= nil && newBoard.BestPositions[to].BestLength+1+cdown.Dist<b.BestPositions[from].BestLength {
			b.BestPositions[from].BestLength = newBoard.BestPositions[to].BestLength+1+cdown.Dist
			b.BestPositions[from].BestX = x
			b.BestPositions[from].BestY = y
			b.BestPositions[from].BestDir = direction.U
		}
	}
	
	cleft := b.Get(x-1,y)
	cright := b.Get(x+1,y)
	if !box.IsCheckedRight {
		box.IsCheckedRight = true
		if (cleft.IsFree && cright.TypeOf != CellTypeWall && !cright.HasBox) {
			box.CanMoveRight = true
			newBoard:= b.MoveBoxAndCheck(x,y,direction.R,boards)
			if newBoard.GetGoodBoxMoveCount() > 0 || newBoard.BestPositions[to].BestLength == 0 {
				box.ShallNotMoveRight = false
				b.CheckEveryDist(fromx,fromy)
				if newBoard.BestPositions[to].BestLength+1+cleft.Dist<b.BestPositions[from].BestLength {
					b.BestPositions[from].BestLength = newBoard.BestPositions[to].BestLength+1+cleft.Dist
					b.BestPositions[from].BestX = x
					b.BestPositions[from].BestY = y
					b.BestPositions[from].BestDir = direction.R
				}
			}
		}
	} else if !XYCheckedAlready && box.CanMoveRight && !box.ShallNotMoveRight {
		newBoard:= b.GetOldMoveBox(x,y,direction.R,boards)
		b.CheckEveryDist(fromx,fromy)
		if newBoard!= nil && newBoard.BestPositions[to].BestLength+1+cleft.Dist<b.BestPositions[from].BestLength {
			b.BestPositions[from].BestLength = newBoard.BestPositions[to].BestLength+1+cleft.Dist
			b.BestPositions[from].BestX = x
			b.BestPositions[from].BestY = y
			b.BestPositions[from].BestDir = direction.R
		}
	}
	if !box.IsCheckedLeft {
		box.IsCheckedLeft = true
		if (cright.IsFree && cleft.TypeOf != CellTypeWall && !cleft.HasBox) {
			box.CanMoveLeft = true
			newBoard:= b.MoveBoxAndCheck(x,y,direction.L,boards)
			if newBoard.GetGoodBoxMoveCount() > 0 || newBoard.BestPositions[to].BestLength == 0 { 
				box.ShallNotMoveLeft = false
				b.CheckEveryDist(fromx,fromy)
				if newBoard.BestPositions[to].BestLength+1+cright.Dist<b.BestPositions[from].BestLength {
					b.BestPositions[from].BestLength = newBoard.BestPositions[to].BestLength+1+cright.Dist
					b.BestPositions[from].BestX = x
					b.BestPositions[from].BestY = y
					b.BestPositions[from].BestDir = direction.L
				}
			}
		}
	} else if !XYCheckedAlready && box.CanMoveLeft && !box.ShallNotMoveLeft {
		newBoard:= b.GetOldMoveBox(x,y,direction.L,boards)
		b.CheckEveryDist(fromx,fromy)
		if newBoard!= nil && newBoard.BestPositions[to].BestLength+1+cright.Dist<b.BestPositions[from].BestLength {
			b.BestPositions[from].BestLength = newBoard.BestPositions[to].BestLength+1+cright.Dist
			b.BestPositions[from].BestX = x
			b.BestPositions[from].BestY = y
			b.BestPositions[from].BestDir = direction.L
		}
	}

	box.XYChecked[from] = true
}

func (b *Board) _CheckEveryBoxMove(boards map[string]*Board) {
	PlayerPos := Position{X:b.Player.X,Y:b.Player.Y}
	if b.BestPositions[PlayerPos] == nil {
		b.BestPositions[PlayerPos] = &BestPosition{BestLength:1000,BestX:-1,BestY:-1}
	}
	
	for i :=0;i<len(b.Boxes);i++ {
		// it loses player pos because of oldmovebox uses
		b.Player.X = PlayerPos.X
		b.Player.Y = PlayerPos.Y
		y := b.Boxes[i].Y
		x := b.Boxes[i].X
		b._CheckOneBoxMove(x,y,boards)
	}
}

func (b *Board) ResetFreeSpace() {
	for i :=0;i<len(b.Cells);i++ {
		b.Cells[i].IsFree = false
		b.Cells[i].Dist = 999
	}
}

// Private Checkup every Free Space from position
func (b *Board) _CheckEveryFreeSpace(x, y, dist int) {
	c := b.Get(x,y)
	if (c.TypeOf == CellTypeWall || c.HasBox || c.IsFree && dist >= c.Dist) {
		return
	}
	c.Dist = dist
	c.IsFree = true
	b._CheckEveryFreeSpace(x-1,y,dist+1)
	b._CheckEveryFreeSpace(x+1,y,dist+1)
	b._CheckEveryFreeSpace(x,y-1,dist+1)
	b._CheckEveryFreeSpace(x,y+1,dist+1)
}

func (b *Board) CheckEveryFreeSpace(x, y int) {
	b.ResetFreeSpace()
	b._CheckEveryFreeSpace(x,y,0)
}

func (b *Board) ResetDist() {
	for i :=0;i<len(b.Cells);i++ {
		b.Cells[i].Dist = 999
	}
}

// Private Checkup every Free Space from position
func (b *Board) _CheckEveryDist(x, y, dist int) {
	c := b.Get(x,y)
	if (!c.IsFree || dist >= c.Dist) {
		return
	}
	c.Dist = dist
	b._CheckEveryDist(x-1,y,dist+1)
	b._CheckEveryDist(x+1,y,dist+1)
	b._CheckEveryDist(x,y-1,dist+1)
	b._CheckEveryDist(x,y+1,dist+1)
}

func (b *Board) CheckEveryDist(x, y int) {
	b.ResetDist()
	b._CheckEveryDist(x,y,0)
}

func (b *Board) _CheckEveryBoxMoveFromPlayer(boards map[string]*Board) {
	Pos := Position{X:b.Player.X,Y:b.Player.Y}

	if b.BestPositions[Pos] == nil {
		b.BestPositions[Pos] = &BestPosition{BestLength:1000,BestX:-1,BestY:-1}
	}

	if b.BestPositions[Pos].BestLength==999 || b.BestPositions[Pos].BestLength ==0 { return }

	if b._CheckEveryBoxIsTrap() {
		b.BestPositions[Pos].BestLength = 999
	} else if b.IsComplete() {
		b.BestPositions[Pos].BestLength = 0
	} else {
		b._CheckEveryBoxMove(boards)
	}
}

// Checkup every Free Space from player position
func (b *Board) CheckEveryBoxMoveFromPlayer(boards map[string]*Board) {
	X := b.Player.X
	Y := b.Player.Y
	Pos := Position{X:X,Y:Y}
	//fmt.Println("-----")
	//b.Print()
	if b.BestPositions[Pos] == nil {
		b.BestPositions[Pos] = &BestPosition{BestLength:1000,BestX:-1,BestY:-1}
	} else {
		b.GetBestPosition().BestLength = 1000
	}
	b._ResetCanBoxMove()
	b.CheckEveryFreeSpace(b.Player.X,b.Player.Y)
	b._CheckEveryBoxMoveFromPlayer(boards)
	b.Player = NewPlayer(X,Y)
	b.CheckEveryDist(X,Y)
}


func (b *Board) MoveBox(x,y int, dir direction.Direction) {
	Pos := Position{X:b.Player.X,Y:b.Player.Y}
	b.BestPositions[Pos] = &BestPosition{BestLength:1000,BestX:-1,BestY:-1}

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
}

func (b *Board) GetBoard(boards map[string]*Board) *Board {
	b.CheckEveryFreeSpace(b.Player.X,b.Player.Y)

	newBoard := b
	boardName := newBoard.GetString()
	tempBoard := boards[boardName]
	if tempBoard == nil {
		boards[boardName] = newBoard
		newBoard._ResetCanBoxMove()
	} else {
		if tempBoard.Player.X != newBoard.Player.X || tempBoard.Player.Y != newBoard.Player.Y {
			tempBoard.Player.X = newBoard.Player.X
			tempBoard.Player.Y = newBoard.Player.Y
			tempBoard.CheckEveryDist(tempBoard.Player.X,tempBoard.Player.Y)
		}
		newBoard = tempBoard 
	}

	return newBoard
}

func (b *Board) GetOldMoveBox(x,y int, dir direction.Direction, boards map[string]*Board) *Board {
	box := &b.Boxes[b.Get(x,y).Box]
	var tempBoard *Board
	switch(dir) {
		case direction.L : if box.LeftBoard != nil { tempBoard = box.LeftBoard }
		case direction.R : if box.RightBoard != nil { tempBoard = box.RightBoard }
		case direction.U : if box.UpBoard != nil { tempBoard = box.UpBoard }
		case direction.D : if box.DownBoard != nil { tempBoard = box.DownBoard }
	}
	//fmt.Println("check old move box")
	if tempBoard != nil {
		//fmt.Println("get old move box")
		tempBoard.Player.X = x
		tempBoard.Player.Y = y
		tempBoard.CheckEveryDist(tempBoard.Player.X,tempBoard.Player.Y)
		return tempBoard
	}
	return nil
}

// assume x,y is a box
func (b *Board) MakeMoveBox(x,y int, dir direction.Direction, boards map[string]*Board) *Board {
	box := &b.Boxes[b.Get(x,y).Box]
	newBoard := b.Duplicate()
	newBoard.MoveBox(x,y,dir)
	newBoard = newBoard.GetBoard(boards)
	switch(dir) {
		case direction.L : box.LeftBoard = newBoard
		case direction.R : box.RightBoard = newBoard
		case direction.U : box.UpBoard = newBoard
		case direction.D : box.DownBoard = newBoard
	}
	return newBoard
}

// assume it
func (b *Board) MoveBoxAndCheck(x,y int, dir direction.Direction, boards map[string]*Board) *Board {
	tempboard := b.GetOldMoveBox(x,y,dir,boards)
	if tempboard != nil { return tempboard }

	newboard := b.MakeMoveBox(x,y,dir,boards)
	newboard._CheckEveryBoxMoveFromPlayer(boards)

	return newboard
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
