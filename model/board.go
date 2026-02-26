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
	IsPath bool
	PathDir direction.Direction
	Box int
	Dist map[Position]int
}

type Position struct {
	X,Y int
}

type Box struct {
	X int
	Y int
	IsDead bool

	CanMove []bool
	ShallNotMove []bool
	IsChecked []bool

	XYChecked map[Position]bool

	DirBoards []*Board
}

func NewBox(x,y int) *Box {
	b := &Box{X:x,Y:y,XYChecked:make(map[Position]bool)}
	b.CanMove = make([]bool,4)
	b.ShallNotMove = make([]bool,4)
	b.IsChecked = make([]bool,4)
	b.DirBoards = make([]*Board,4)
	return b
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
				b.Boxes = append(b.Boxes,*NewBox(x,y))
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
				b.Boxes = append(b.Boxes,*NewBox(x,y))
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
		d.Boxes[i] = *NewBox(box.X,box.Y)
		copy(d.Boxes[i].CanMove,box.CanMove)
		copy(d.Boxes[i].ShallNotMove,box.ShallNotMove)
		copy(d.Boxes[i].IsChecked,box.IsChecked)
		d.Boxes[i].XYChecked = make(map[Position]bool)
	}

	d.Player = NewPlayer(b.Player.X,b.Player.Y)

	return d
}

func (b *Board) _ResetCanBoxMove() {
	for i :=0;i<len(b.Boxes);i++ {
		b.Boxes[i].IsDead = false

		copy(b.Boxes[i].CanMove,[]bool{false,false,false,false})
		copy(b.Boxes[i].ShallNotMove,[]bool{true,true,true,true})
		copy(b.Boxes[i].IsChecked,[]bool{false,false,false,false})
	}
}

func (b *Board) _CheckOneBoxIsDead(x,y int) bool {
	box := &b.Boxes[b.Get(x,y).Box]
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
	for i :=0;i<len(b.Boxes);i++ {
		box := b.Boxes[i]
		y := box.Y
		x := box.X
		c := b.Get(x,y)

		if c.TypeOf != CellTypeGoal && b._CheckOneBoxIsDead(x,y) { count++ }
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

type CellPile struct  {
	Cells []*Cell
}

func NewCellPile() *CellPile {
	return &CellPile{ Cells : make([]*Cell,0) }
}

func (c *CellPile) Push(cell *Cell) {
	c.Cells = append(c.Cells, cell)
}

func (c *CellPile) Pop() *Cell {
	if len(c.Cells) == 0 { return nil }
	cell := c.Cells[0]
	c.Cells = c.Cells[1:len(c.Cells)]
	return cell
}

// assume x,y is a box
func (b *Board) _CheckOneBoxIsStuck(x,y int, freeCells map[*Cell]bool) bool {
	cellup := b.Get(x,y-1)
	celldown := b.Get(x,y+1)
	cellleft := b.Get(x-1,y)
	cellright := b.Get(x+1,y)
	stuckup := cellup.TypeOf == CellTypeWall || (cellup.HasBox && !freeCells[cellup])
	stuckdown := celldown.TypeOf == CellTypeWall || (celldown.HasBox && !freeCells[celldown])
	stuckleft := cellleft.TypeOf == CellTypeWall || (cellleft.HasBox && !freeCells[cellleft])
	stuckright := cellright.TypeOf == CellTypeWall || (cellright.HasBox && !freeCells[cellright])

	if b.Get(x,y).HasBox && ((!stuckup && !stuckdown) || (!stuckleft && !stuckright)) {
		return false
	}
	return true
}

func (b *Board) _CheckEveryBoxIsStuck() bool {
	pile := NewCellPile()
	free := make(map[*Cell]bool)

	// pile cells
	for i :=0;i<len(b.Boxes);i++ {
		c := b.Get(b.Boxes[i].X,b.Boxes[i].Y)
		pile.Push(c)
	}
	
	// check every box as it is not free
	for {
		nextPile := NewCellPile()
		pilecount := len(pile.Cells)

		// check every box
		for current := pile.Pop(); current!=nil; current=pile.Pop() {
			box := &b.Boxes[current.Box]
			x := box.X
			y := box.Y

			if !b._CheckOneBoxIsStuck(x,y,free) {
				free[current] = true
			} else {
				nextPile.Push(current)
			}
		}
		if pilecount == len(nextPile.Cells) { break }
		pile = nextPile
	}

	// mark every stuck box as dead
	traped := false
	for i :=0;i<len(b.Boxes);i++ {
		box := &b.Boxes[i]
		cell := b.Get(box.X,box.Y)

		if !free[cell] && cell.TypeOf != CellTypeGoal {
			box.IsDead = true
			traped = true
		}
	}
	return traped
}

func (b *Board) _CheckEveryBoxIsTrap() bool {
	traped := false
	traped = b._CheckEveryBoxIsStuck() || b._CheckEveryBoxIsTrapByWall() // || b._CheckEveryBoxIsDead() // Note : Stuck includes Dead ones
	return traped
}

func getMoveDirection(dir direction.Direction) (int,int) {
	switch(dir) {
		case direction.U : return 0,1
		case direction.D : return 0,-1
		case direction.L : return 1,0
		case direction.R : return -1,0
	}
	return 0,0
}

func (b *Board) _CheckOneBoxMoveInDir(x,y, fromx,fromy int, box *Box, from, to Position, dir direction.Direction, boards map[string]*Board) {

	dx, dy := getMoveDirection(dir)

	cup := b.Get(x+dx,y+dy)
	cdown := b.Get(x-dx,y-dy)

	if !box.IsChecked[dir] {
		box.IsChecked[dir] = true
		if (cup.IsFree && cdown.TypeOf != CellTypeWall && !cdown.HasBox) {
			box.CanMove[dir] = true
			newBoard:= b.MoveBoxAndCheck(x,y,dir,boards)
			if newBoard.GetGoodBoxMoveCount() > 0 || newBoard.BestPositions[to].BestLength == 0 {
				box.ShallNotMove[dir] = false
				b.CheckEveryDist(fromx,fromy)
				if newBoard.BestPositions[to].BestLength+1+cup.Dist[from]<b.BestPositions[from].BestLength {
					b.BestPositions[from].BestLength = newBoard.BestPositions[to].BestLength+1+cup.Dist[from]
					b.BestPositions[from].BestX = x
					b.BestPositions[from].BestY = y
					b.BestPositions[from].BestDir = dir
				}
			}
		}
	} else if !box.XYChecked[from] && box.CanMove[dir] && !box.ShallNotMove[dir] {
		newBoard:= b.GetOldMoveBox(x,y,dir,boards)
		b.CheckEveryDist(fromx,fromy)
		if newBoard!= nil && newBoard.BestPositions[to].BestLength+1+cup.Dist[from]<b.BestPositions[from].BestLength {
			b.BestPositions[from].BestLength = newBoard.BestPositions[to].BestLength+1+cup.Dist[from]
			b.BestPositions[from].BestX = x
			b.BestPositions[from].BestY = y
			b.BestPositions[from].BestDir = dir
		}
	}

}

// Assume x,y got a box
func (b *Board) _CheckOneBoxMove(x,y int,boards map[string]*Board) {
	c := b.Get(x,y)
	box := &(b.Boxes[c.Box])
	
	fromx := b.Player.X
	fromy := b.Player.Y
	from := Position{X:fromx,Y:fromy}
	to := Position{X:x,Y:y}

	b._CheckOneBoxMoveInDir(x,y,fromx,fromy, box, from, to, direction.D, boards)
	b._CheckOneBoxMoveInDir(x,y,fromx,fromy, box, from, to, direction.U, boards)
	b._CheckOneBoxMoveInDir(x,y,fromx,fromy, box, from, to, direction.R, boards)
	b._CheckOneBoxMoveInDir(x,y,fromx,fromy, box, from, to, direction.L, boards)
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

func (b *Board) ResetFreeSpace(from Position) {
	for i :=0;i<len(b.Cells);i++ {
		b.Cells[i].IsFree = false
		if b.Cells[i].Dist == nil {
			b.Cells[i].Dist = make(map[Position]int)
		}
		b.Cells[i].Dist[from] = 999
	}
}

// Private Checkup every Free Space from position
func (b *Board) _CheckEveryFreeSpace(from Position, x, y, dist int) {
	c := b.Get(x,y)
	if (c.TypeOf == CellTypeWall || c.HasBox || dist >= c.Dist[from]) {
		return
	}
	c.Dist[from] = dist
	c.IsFree = true
	b._CheckEveryFreeSpace(from,x-1,y,dist+1)
	b._CheckEveryFreeSpace(from,x+1,y,dist+1)
	b._CheckEveryFreeSpace(from,x,y-1,dist+1)
	b._CheckEveryFreeSpace(from,x,y+1,dist+1)
}

func (b *Board) CheckEveryFreeSpace(x, y int) {
	from := Position{X:x,Y:y}
	c := b.Get(x,y)
	_, ok := c.Dist[from]
	if ok { return } 
	b.ResetFreeSpace(from)
	b._CheckEveryFreeSpace(from,x,y,0)
}

func (b *Board) ResetDist(from Position) {
	for i :=0;i<len(b.Cells);i++ {
		b.Cells[i].Dist[from] = 999
	}
}

// Private Checkup every Free Space from position
func (b *Board) _CheckEveryDist(from Position,x, y, dist int) {
	c := b.Get(x,y)
	if (!c.IsFree || dist >= c.Dist[from]) {
		return
	}
	c.Dist[from] = dist
	b._CheckEveryDist(from,x-1,y,dist+1)
	b._CheckEveryDist(from,x+1,y,dist+1)
	b._CheckEveryDist(from,x,y-1,dist+1)
	b._CheckEveryDist(from,x,y+1,dist+1)
}

func (b *Board) CheckEveryDist(x, y int) {
	from := Position{X:x,Y:y}
	c := b.Get(x,y)
	_, ok := c.Dist[from]
	if ok { return }

	b.ResetDist(from)
	b._CheckEveryDist(from,x,y,0)
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
	b.FindBestPath()
}


func (b *Board) MoveBox(x,y int, dir direction.Direction) {
	Pos := Position{X:b.Player.X,Y:b.Player.Y}
	b.BestPositions[Pos] = &BestPosition{BestLength:1000,BestX:-1,BestY:-1}

	lastCell := b.Get(x,y)
	b.Player.X = x
	b.Player.Y = y
	var newCell *Cell
	box := &b.Boxes[lastCell.Box]
	dx,dy := getMoveDirection(dir)
	newCell = b.Get(x-dx,y-dy)
	box.X = x-dx
	box.Y = y-dy

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
	if box.DirBoards[dir] != nil { tempBoard = box.DirBoards[dir] }
	if tempBoard != nil {
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
	box.DirBoards[dir] = newBoard
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

func (b *Board) ResetPath() {
	for i :=0;i<len(b.Cells);i++ {
		b.Cells[i].IsPath = false
		b.Cells[i].PathDir = direction.None
	}
}

func (b *Board) _FindReverseBestPath(x,y int, p Position, l int, pathDir direction.Direction) bool {
	if p.X==x && p.Y==y { 
		b.Get(x,y).PathDir = pathDir
		return true 
	}

	if b.Get(x,y).IsFree && b.Get(x,y).Dist[p]==l {
		b.Get(x,y).IsPath = true
		b.Get(x,y).PathDir = pathDir
		
		if b._FindReverseBestPath(x-1,y,p,l-1,direction.R) { return true }
		if b._FindReverseBestPath(x+1,y,p,l-1,direction.L) { return true }
		if b._FindReverseBestPath(x,y-1,p,l-1,direction.D) { return true }
		if b._FindReverseBestPath(x,y+1,p,l-1,direction.U) { return true }
	}
	return false
}

// assume best move is set up
func (b *Board) FindBestPath() {
	b.ResetPath()

	if b.GetGoodBoxMoveCount() == 0 { return }

	position := Position{X:b.Player.X,Y:b.Player.Y}
	
	bestposition := b.GetBestPosition()
	x := bestposition.BestX
	y := bestposition.BestY

	if bestposition.BestLength==0 { return }

	switch(bestposition.BestDir) {
		case direction.L :  x = x+1
		case direction.R :  x = x-1
		case direction.U :  y = y+1
		case direction.D :  y = y-1
	}
	l := b.Get(x,y).Dist[position]

	b._FindReverseBestPath(x,y,position,l,bestposition.BestDir)
}

func (b *Board) GetBoxMoveCount() int {
	count := 0
	for _, box := range b.Boxes {
		for _, canMove := range box.CanMove {
			if canMove { count = count+1 }
		}
	}
	return count
}

func (b *Board) GetGoodBoxMoveCount() int {
	count := 0
	for _, box := range b.Boxes {
		for i:=0;i<4;i++ {
			if box.CanMove[i] && !box.ShallNotMove[i] {count = count+1 }
		}
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
