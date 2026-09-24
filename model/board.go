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
	X,Y int
	TypeOf cellType
	HasBox bool
	Box    int

	IsFree bool
	IsDead bool
	CanMove []bool
}

type Board struct {
	Width, Height int
	Boxes []int
	BestBoxes []int
	Goals []int
	Distances [][]int
	MaxMoves int
	FreeCells []int

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
	boxes  := 0

	for y := 0; y < b.Height; y++ {
		for x := 0; x < b.Width; x++ {
			index := (y*b.Width)+x
			code := string(mapData[index])
			cell := Cell{CanMove:make([]bool,4),X:x,Y:y}
			switch code {
			case "@":
				b.Player = NewPlayer(x, y)
			case "$":
				cell.HasBox = true
				cell.Box    = boxes
				boxes++
				b.Boxes = append(b.Boxes,index)
			case ".":
				cell.TypeOf = CellTypeGoal
				b.Goals = append(b.Goals,index)
			case "#":
				cell.TypeOf = CellTypeWall
			case "+":
				cell.TypeOf = CellTypeGoal
				b.Player = NewPlayer(x, y)
				b.Goals = append(b.Goals,index)
			case "*":
				cell.TypeOf = CellTypeGoal
				cell.HasBox = true
				cell.Box    = boxes
				boxes++
				b.Boxes = append(b.Boxes,index)
				b.Goals = append(b.Goals,index)
			}
			b.Cells[index] = cell
		}
	}

	b.Distances = make([][]int,len(b.Goals))
	for i := 0; i<len(b.Goals); i++ {
		y := b.Goals[i]/b.Width
		x := b.Goals[i] - y*b.Width
		b.Distances[i] = b.GetDistancesFrom(x,y)
	}

	return &b
}

// Get - Returns the cell at the given location
func (b *Board) Get(x, y int) *Cell {
	return &b.Cells[(y*b.Width)+x]
}

// Get - Returns the cell at the given location
func (b *Board) GetBox(i int) *Cell {
	return &b.Cells[b.Boxes[i]]
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

func (b *Board) GetDistancesFromRecursive(x,y,d int, distances *[]int) {
	index := (y*b.Width)+x
	if b.Cells[index].TypeOf == CellTypeWall || ((*distances)[index]!=0 && (*distances)[index]<d) {
		return
	}
	(*distances)[index] = d
	b.GetDistancesFromRecursive(x-1,y,d+1,distances)
	b.GetDistancesFromRecursive(x+1,y,d+1,distances)
	b.GetDistancesFromRecursive(x,y-1,d+1,distances)
	b.GetDistancesFromRecursive(x,y+1,d+1,distances)
}

func (b *Board) GetDistancesFrom(x,y int) []int {
	distances := make([]int,b.Height*b.Width)

	b.GetDistancesFromRecursive(x,y,1,&distances)	

	return distances
}

func (b *Board) GetMinBestGoalBox(goal int,inactive []bool) int {
	min := 999
	current := -1
	for i:=0;i<len(b.Boxes);i++ {
		if !inactive[i] {
			distance := b.Distances[goal][b.Boxes[i]]
			if distance < min {
				min = distance
				current = i
			}
		}
	}
	return current
}

// Get the best box for each Goal
func (b *Board) GetBestBoxFromDistance() []int {
	bestBoxes := make([]int,len(b.Goals))
	
	// it's fastest to use inactive because of the bool default value
	inactive := make([]bool,len(b.Boxes))
	
	for goal:=0;goal<len(b.Goals);goal++ {
		bestBox := b.GetMinBestGoalBox(goal,inactive)
		inactive[bestBox] = true
		bestBoxes[goal] = bestBox
	}
	return bestBoxes
}

// need b.BestBoxes := b.GetBestBoxFromDistance()
func (b *Board) GetSumOfBestBoxDistances() int {
	sum := 0
	
	for i:=0;i<len(b.Goals);i++ {
		distance := b.Distances[i][b.Boxes[b.BestBoxes[i]]]-1
		sum = sum + distance
	}
	
	return sum
}

func (b *Board) ResetFreeCells() {
	if len(b.FreeCells)>0 {
		for i:=0;i<len(b.FreeCells);i++ {
			b.Cells[b.FreeCells[i]].IsFree = false
		}
		b.FreeCells = []int{}
	}
}

func (b *Board) ResetCanMove() {
	for i:=0;i<len(b.Boxes);i++ {
		b.Cells[b.Boxes[i]].CanMove[direction.U] = false
		b.Cells[b.Boxes[i]].CanMove[direction.D] = false
		b.Cells[b.Boxes[i]].CanMove[direction.L] = false
		b.Cells[b.Boxes[i]].CanMove[direction.R] = false
	}
}

func (b *Board) GetDirections(dir direction.Direction) (int,int) {
	switch(dir) {
		case direction.U : return 0,-1
		case direction.D : return 0,1
		case direction.L : return -1,0
		case direction.R : return 1,0
	}
	return 0,0
}

func (b *Board) CheckIfBoxInDirectionCanMove(x,y int, dir direction.Direction) {
	dirx,diry := b.GetDirections(dir)
	box := b.Get(x+dirx,y+diry)
	nextthebox := b.Get(x+dirx+dirx,y+diry+diry)
	if nextthebox.TypeOf != CellTypeWall && !nextthebox.HasBox { 
		box.CanMove[dir] = true
	}
}

// true if it finds a box
func (b *Board) FindFreeCellsFrom(x,y int) bool {
	index := y*b.Width + x
	c := &b.Cells[index]
	if c.IsFree || c.TypeOf == CellTypeWall { return false }
	if c.HasBox { return true }

	c.IsFree = true
	b.FreeCells = append(b.FreeCells,index)
	if b.FindFreeCellsFrom(x+1,y) { b.CheckIfBoxInDirectionCanMove(x,y,direction.R) }
	if b.FindFreeCellsFrom(x-1,y) { b.CheckIfBoxInDirectionCanMove(x,y,direction.L) }
	if b.FindFreeCellsFrom(x,y+1) { b.CheckIfBoxInDirectionCanMove(x,y,direction.D) }
	if b.FindFreeCellsFrom(x,y-1) { b.CheckIfBoxInDirectionCanMove(x,y,direction.U) }
	return false
}

func (b *Board) FindFreeCells() {
	b.FindFreeCellsFrom(b.Player.X,b.Player.Y)
}

type MoveType int

const (
	PlayerBlockedByWall MoveType = iota
	PlayerBlockedByBox
	PlayerMoveAndPush
	PlayerMove
)

func (b *Board) MovePlayer(dir direction.Direction, undo bool) (bool, MoveType) {
	lastX := b.Player.X
	lastY := b.Player.Y
	targetX := lastX
	targetY := lastY
	nextX := targetX
	nextY := targetY

	switch dir {
	case direction.U:
		targetY--
		nextY -= 2
	case direction.D:
		targetY++
		nextY += 2
	case direction.L:
		targetX--
		nextX -= 2
	case direction.R:
		targetX++
		nextX += 2
	}

	targetCell := b.Get(targetX, targetY)

	if targetCell.TypeOf == CellTypeWall {
		return false,PlayerBlockedByWall
	} else {
		if targetCell.HasBox {
			nextCell := b.Get(nextX, nextY)
			if nextCell.TypeOf == CellTypeWall {
				return false,PlayerBlockedByWall
			} else if nextCell.HasBox {
				return false,PlayerBlockedByBox
			} else {
				nextCell.Box = targetCell.Box // works because we can't push 2 boxes at the same time
				b.Boxes[nextCell.Box] = nextY*b.Width+nextX
				targetCell.HasBox = false
				nextCell.HasBox = true
				b.Player.X = targetX
				b.Player.Y = targetY
				if undo { b.LastMove = NewLastMove(lastX,lastY,targetCell,nextCell,b.LastMove) }
				b.Update()
				return true,PlayerMoveAndPush
			}
		} else {
			b.Player.X = targetX
			b.Player.Y = targetY
			if undo { b.LastMove = NewLastMove(lastX,lastY,nil,nil,b.LastMove) }
			b.Update()
			return true,PlayerMove
		}
	}

}

type UndoType int

const (
	PlayerUndoAndUnpush UndoType = iota
	PlayerUndoMove
	NoUndo
)

func (b *Board) UndoLastMove() (bool,UndoType) {
	if b.LastMove == nil {
		return false,NoUndo
	}
	b.Player.X = b.LastMove.LastX
	b.Player.Y = b.LastMove.LastY
	var ret UndoType
	if b.LastMove.LastTargetCell != nil {
		b.LastMove.LastTargetCell.HasBox = true
		b.LastMove.LastNextCell.HasBox = false
		b.LastMove.LastTargetCell.Box = b.LastMove.LastNextCell.Box
		b.Boxes[b.LastMove.LastTargetCell.Box] = b.LastMove.LastTargetCell.Y*b.Width+b.LastMove.LastTargetCell.X
		ret = PlayerUndoMove
	} else { ret = PlayerUndoAndUnpush }
	b.LastMove = b.LastMove.PreviousMove
	b.Update()
	return true,ret
}

func (b *Board) _CheckOneBoxIsDead(x,y int) bool {
	box := b.Get(x,y)
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
		box := b.GetBox(i)
		y := box.Y
		x := box.X

		if box.TypeOf != CellTypeGoal && b._CheckOneBoxIsDead(x,y) { count++ }
	}
	return count > 0
}

func (b *Board) _CheckOneBoxIsTrapByDirWall(x,y int, dirx, diry int) bool {
	c := b.Get(x,y)
	if !c.HasBox { return false }
	if c.IsDead { return true }
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
		c.IsDead = true
		return true
	} else { return false }
}

func (b *Board) _CheckEveryBoxIsTrapByWall() bool {
	count := 0
	for i :=0;i<len(b.Boxes);i++ {
		box := b.GetBox(i)
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
		c := b.Get(b.GetBox(i).X,b.GetBox(i).Y)
		pile.Push(c)
	}
	
	// check every box as it is not free
	for {
		nextPile := NewCellPile()
		pilecount := len(pile.Cells)

		// check every box
		for current := pile.Pop(); current!=nil; current=pile.Pop() {
			x := current.X
			y := current.Y

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
		cell := b.GetBox(i)

		if !free[cell] && cell.TypeOf != CellTypeGoal {
			cell.IsDead = true
			traped = true
		}
	}
	return traped
}

// Note : b._CheckEveryBoxIsDead() is not used because _CheckEveryBoxIsStuck includes Dead ones
func (b *Board) _CheckEveryBoxIsTrap() bool {
	traped := false
	traped = b._CheckEveryBoxIsStuck() || b._CheckEveryBoxIsTrapByWall()
	return traped
}

func (b *Board) Update() {
	b.BestBoxes = b.GetBestBoxFromDistance()
	b.MaxMoves = b.GetSumOfBestBoxDistances()

	b.ResetFreeCells()
	b.ResetCanMove()
	b.FindFreeCells()
	b._CheckEveryBoxIsTrap()
}