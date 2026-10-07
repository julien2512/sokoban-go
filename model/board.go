package model

import (
	"fmt"
	"slices"
	"strings"
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
	MoveCost []int
}

type SubBoard struct {
	Boxes []int
	Cells         []Cell
	Player        *Player

	BestBoxes []int
	MaxMoves int

	BestSubBoardDir direction.Direction
	BestSubBoardBox int
	ParentSubBoard *SubBoard
	ParentSubBoardDir direction.Direction
	ParentSubBoardBox int

	str string

	FreeCells []int
}

func (S *SubBoard) GetString() string {
	if S.str != "" { return S.str }
	index := 0
	str := ""
	for i:=0;i<len(S.Cells)/2;i++ {
		num:=0
		val:=0
		for j:=0;j<2 && index<len(S.Cells);j++ {
			cell := S.Cells[index]
			if cell.HasBox {
				val = 1
			} else if cell.IsFree {
				val = 2
			} else if cell.TypeOf==CellTypeWall {
				val = 3
			} else { val = 0 }
			num = num*4+val
			index++
		}
		str += string(65+num)
	}
	S.str = str
	return str
}

func (c *Cell) CloneCell() *Cell {
	clone := &Cell{X:c.X,Y:c.Y,TypeOf:c.TypeOf,HasBox:c.HasBox,Box:c.Box,IsFree:c.IsFree,IsDead:c.IsDead}
	clone.CanMove = make([]bool,4)
	clone.MoveCost = make([]int,4)
	for i:=0;i<4;i++ {
		clone.CanMove[i] = c.CanMove[i]
		clone.MoveCost[i] = c.MoveCost[i]
	}
	return clone
}

func (s *SubBoard) CloneSubBoard() *SubBoard {
	clone := &SubBoard{}
	clone.Boxes = make([]int,len(s.Boxes))
	clone.Cells = make([]Cell,len(s.Cells))
	clone.Player = &Player{X:s.Player.X,Y:s.Player.Y}
	clone.FreeCells = make([]int,len(s.FreeCells))
	for i:=0;i<len(s.Boxes);i++ {
		clone.Boxes[i] = s.Boxes[i]
	}
	for i:=0;i<len(s.Cells);i++ {
		clone.Cells[i].X = s.Cells[i].X
		clone.Cells[i].Y = s.Cells[i].Y
		clone.Cells[i].TypeOf = s.Cells[i].TypeOf
		clone.Cells[i].HasBox = s.Cells[i].HasBox
		clone.Cells[i].Box = s.Cells[i].Box
		clone.Cells[i].IsFree = s.Cells[i].IsFree
		clone.Cells[i].IsDead = s.Cells[i].IsDead
		clone.Cells[i].CanMove = make([]bool,4)
		clone.Cells[i].MoveCost = make([]int,4)
	}
	for i:=0;i<len(s.FreeCells);i++ {
		clone.FreeCells[i] = s.FreeCells[i]
	}
	return clone
}

type Board struct {
	Width, Height int
	Goals []int
	Distances [][]int

	S *SubBoard

	LastMove      *LastMove
}

// NewBoard - Creates a board (map data encoding: Player "@", Box "$", Goal ".", Wall "#", Goal+Player "+", Goal+Box "*")
func NewBoard(mapData string, boardWidth, boardHeight int) *Board {
	b := Board{S:&SubBoard{}}

	b.Width = boardWidth
	b.Height = boardHeight

	b.S.Cells = make([]Cell, b.Width*b.Height)
	boxes  := 0

	for y := 0; y < b.Height; y++ {
		for x := 0; x < b.Width; x++ {
			index := (y*b.Width)+x
			code := string(mapData[index])
			cell := Cell{CanMove:make([]bool,4),MoveCost:make([]int,4),X:x,Y:y}
			switch code {
			case "@":
				b.S.Player = NewPlayer(x, y)
			case "$":
				cell.HasBox = true
				cell.Box    = boxes
				boxes++
				b.S.Boxes = append(b.S.Boxes,index)
			case ".":
				cell.TypeOf = CellTypeGoal
				b.Goals = append(b.Goals,index)
			case "#":
				cell.TypeOf = CellTypeWall
			case "+":
				cell.TypeOf = CellTypeGoal
				b.S.Player = NewPlayer(x, y)
				b.Goals = append(b.Goals,index)
			case "*":
				cell.TypeOf = CellTypeGoal
				cell.HasBox = true
				cell.Box    = boxes
				boxes++
				b.S.Boxes = append(b.S.Boxes,index)
				b.Goals = append(b.Goals,index)
			}
			b.S.Cells[index] = cell
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

func (b *Board) Print() {
	for y :=0;y<b.Height;y++ {
		for x :=0;x<b.Width;x++ {
			c := b.Get(x,y)
			
			if c.HasBox && c.TypeOf == CellTypeGoal && c.X==x && c.Y==y { fmt.Print("*")
			} else if x == b.S.Player.X && y == b.S.Player.Y && c.TypeOf == CellTypeGoal { fmt.Print("+")
			} else if c.TypeOf == CellTypeWall { fmt.Print("#")
			} else if c.TypeOf == CellTypeGoal { fmt.Print(".")
			} else if c.HasBox && c.TypeOf != CellTypeGoal && c.X==x && c.Y==y { fmt.Print("$")
			} else if x == b.S.Player.X && y == b.S.Player.Y && c.TypeOf != CellTypeGoal { fmt.Print("@")
			} else { fmt.Print(" ") }
		}
		fmt.Println()
	}
}

func (S *SubBoard) PrintSubBoard(W,H int) {
	for y :=0;y<H;y++ {
		for x :=0;x<W;x++ {
			c := S.GetFromSubBoard(x,y,W)
			
			if c.HasBox && c.TypeOf == CellTypeGoal && c.X==x && c.Y==y { fmt.Print("*")
			} else if x == S.Player.X && y == S.Player.Y && c.TypeOf == CellTypeGoal { fmt.Print("+")
			} else if c.TypeOf == CellTypeWall { fmt.Print("#")
			} else if c.TypeOf == CellTypeGoal { fmt.Print(".")
			} else if c.HasBox && c.TypeOf != CellTypeGoal && c.X==x && c.Y==y { fmt.Print("$")
			} else if x == S.Player.X && y == S.Player.Y && c.TypeOf != CellTypeGoal { fmt.Print("@")
			} else { fmt.Print(" ") }
		}
		fmt.Println()
	}
}

// Get - Returns the cell at the given location
func (b *Board) Get(x, y int) *Cell {
	return &b.S.Cells[(y*b.Width)+x]
}

// Get - Returns the cell at the given location
func (S *SubBoard) GetFromSubBoard(x, y int, W int) *Cell {
	return &S.Cells[(y*W)+x]
}

// Get - Returns the cell at the given location
func (b *Board) GetBox(i int) *Cell {
	return &b.S.Cells[b.S.Boxes[i]]
}

// IsComplete - Returns true if every goal cell on the board has a box
func (b *Board) IsComplete() bool {
	for _, cell := range b.S.Cells {
		if cell.TypeOf == CellTypeGoal && !cell.HasBox {
			return false
		}
	}
	return true
}

func (b *Board) GetDistancesFromRecursive(x,y,d int, distances *[]int) {
	index := (y*b.Width)+x
	if b.S.Cells[index].TypeOf == CellTypeWall || ((*distances)[index]!=0 && (*distances)[index]<d) {
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
	for i:=0;i<len(b.S.Boxes);i++ {
		if !inactive[i] {
			distance := b.Distances[goal][b.S.Boxes[i]]
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
	inactive := make([]bool,len(b.S.Boxes))
	
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
		distance := b.Distances[i][b.S.Boxes[b.S.BestBoxes[i]]]-1
		sum = sum + distance
	}
	
	return sum
}

func (b *Board) ResetFreeCells() {
	if len(b.S.FreeCells)>0 {
		for i:=0;i<len(b.S.FreeCells);i++ {
			b.S.Cells[b.S.FreeCells[i]].IsFree = false
		}
		b.S.FreeCells = []int{}
	}
}

func (b *Board) ResetCanMove() {
	for i:=0;i<len(b.S.Boxes);i++ {
		b.S.Cells[b.S.Boxes[i]].CanMove[direction.U] = false
		b.S.Cells[b.S.Boxes[i]].CanMove[direction.D] = false
		b.S.Cells[b.S.Boxes[i]].CanMove[direction.L] = false
		b.S.Cells[b.S.Boxes[i]].CanMove[direction.R] = false
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
	c := &b.S.Cells[index]
	if c.IsFree || c.TypeOf == CellTypeWall { return false }
	if c.HasBox { return true }

	c.IsFree = true
	b.S.FreeCells = append(b.S.FreeCells,index)
	if b.FindFreeCellsFrom(x+1,y) { b.CheckIfBoxInDirectionCanMove(x,y,direction.R) }
	if b.FindFreeCellsFrom(x-1,y) { b.CheckIfBoxInDirectionCanMove(x,y,direction.L) }
	if b.FindFreeCellsFrom(x,y+1) { b.CheckIfBoxInDirectionCanMove(x,y,direction.D) }
	if b.FindFreeCellsFrom(x,y-1) { b.CheckIfBoxInDirectionCanMove(x,y,direction.U) }
	return false
}

func (b *Board) FindFreeCells() {
	b.FindFreeCellsFrom(b.S.Player.X,b.S.Player.Y)
}

type MoveType int

const (
	PlayerBlockedByWall MoveType = iota
	PlayerBlockedByBox
	PlayerMoveAndPush
	PlayerMove
)

func (b *Board) MovePlayer(dir direction.Direction, undo bool) (bool, MoveType) {
	lastX := b.S.Player.X
	lastY := b.S.Player.Y
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
				b.S.Boxes[nextCell.Box] = nextY*b.Width+nextX
				targetCell.HasBox = false
				nextCell.HasBox = true
				b.S.Player.X = targetX
				b.S.Player.Y = targetY
				if undo { b.LastMove = NewLastMove(lastX,lastY,targetCell,nextCell,b.LastMove) }
				return true,PlayerMoveAndPush
			}
		} else {
			b.S.Player.X = targetX
			b.S.Player.Y = targetY
			if undo { b.LastMove = NewLastMove(lastX,lastY,nil,nil,b.LastMove) }
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
	b.S.Player.X = b.LastMove.LastX
	b.S.Player.Y = b.LastMove.LastY
	var ret UndoType
	if b.LastMove.LastTargetCell != nil {
		b.LastMove.LastTargetCell.HasBox = true
		b.LastMove.LastNextCell.HasBox = false
		b.LastMove.LastTargetCell.Box = b.LastMove.LastNextCell.Box
		b.S.Boxes[b.LastMove.LastTargetCell.Box] = b.LastMove.LastTargetCell.Y*b.Width+b.LastMove.LastTargetCell.X
		ret = PlayerUndoMove
	} else { ret = PlayerUndoAndUnpush }
	b.LastMove = b.LastMove.PreviousMove
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
	for i :=0;i<len(b.S.Boxes);i++ {
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
	for i :=0;i<len(b.S.Boxes);i++ {
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
	for i :=0;i<len(b.S.Boxes);i++ {
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
	for i :=0;i<len(b.S.Boxes);i++ {
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

func (b *Board) CheckMoveCostsForOneBox(box int,dir direction.Direction) *SubBoard {
	cell := &(b.S.Cells[b.S.Boxes[box]])
	if cell.CanMove[dir] {
		lastsubboard := b.S
		b.S = b.S.CloneSubBoard()
		b.S.ParentSubBoard = lastsubboard
		b.S.ParentSubBoardDir = dir
		b.S.ParentSubBoardBox = box
		
		playerX := cell.X
		playerY := cell.Y
		
		switch dir {
		case direction.U:
			playerY++
		case direction.D:
			playerY--
		case direction.L:
			playerX++
		case direction.R:
			playerX--
		}

		b.S.Player.X = playerX
		b.S.Player.Y = playerY
		
		b.MovePlayer(dir,false)
		b.Update()

		newsubboard := b.S
		b.S = lastsubboard
		return newsubboard
	} else {
		cell.MoveCost[dir] = 1000
		return nil
	}
}

func (b *Board) Update() {
	b.ResetFreeCells()
	b.ResetCanMove()
	b.FindFreeCells()
	if b._CheckEveryBoxIsTrap() {
		b.S.BestBoxes = []int{}
		b.S.MaxMoves = 1000
	} else {
		b.S.BestBoxes = b.GetBestBoxFromDistance()
		b.S.MaxMoves = b.GetSumOfBestBoxDistances()
	}
}


type NextTry struct {
	S *SubBoard
	Deep int
}

func SortNextTry(a, b *NextTry) int {
	if a.S.MaxMoves < b.S.MaxMoves { return -1 }
	if a.S.MaxMoves == b.S.MaxMoves { return strings.Compare(a.S.GetString(),b.S.GetString()) }
	return 1
}

type Guesser struct {
	Checked map[string]*NextTry
	Next    []*NextTry
}

func (g *Guesser) PrintGuesser() {
	fmt.Println("Print Guesser")
	for i:=0;i<len(g.Next);i++ {
		fmt.Printf("Next %0d deep %0d Moves %0d: %s\n",i,g.Next[i].Deep,g.Next[i].S.MaxMoves,g.Next[i].S.GetString())
	}
}

func NewGuesser() *Guesser {
	g := &Guesser{Checked:make(map[string]*NextTry)}
	return g
}

func (g *Guesser) Pop() *NextTry {
	if len(g.Next)==0 { return nil }
	pop := g.Next[0]
	g.Next = g.Next[1:]
	return pop
}

func (g *Guesser) Push(S *SubBoard,Deep int) bool {
	next := &NextTry{S:S,Deep:Deep}
	str := next.S.GetString()
	checked,ok := g.Checked[str]
	if !ok {
		index,_ := slices.BinarySearchFunc(g.Next,next,SortNextTry)
		g.Next = slices.Insert(g.Next,index,next)
		g.Checked[str] = next
		return true
	} else {
		if checked.Deep > Deep {
			g.Checked[str] = next
			index,_ := slices.BinarySearchFunc(g.Next,next,SortNextTry)
			for index < len(g.Next) && g.Next[index].S.GetString()!=str && g.Next[index].S.MaxMoves == next.S.MaxMoves {
				index++
			}
			if index < len(g.Next) && g.Next[index].S.GetString()==str {
				g.Next[index] = next
			}
		}
		return false
	}
}

func (S *SubBoard) CountMoves() int {
	count := 0
	for i:=0;i<len(S.Boxes);i++ {
		for dir:=0;dir<4;dir++ {
			if S.Cells[S.Boxes[i]].CanMove[dir] { count++ }
		}
	}
	return count
}

func (b *Board) Guess() {
//	maxdeep := b.S.MaxMoves+7 // Allow 5 more steps
	goal := false
	var GoalSubBoard *SubBoard
	
	g := NewGuesser()
	g.Push(b.S,0)

	saveSubBoard := b.S

	for len(g.Next)>0 && !goal {
		current := g.Pop()
		currentDeep := current.Deep
		b.S = current.S

//		fmt.Printf("Pop Board %s at %02d moves and deep %02d---------\n",b.S.GetString(),b.S.MaxMoves,currentDeep)
//		fmt.Printf("There is %02d subboards\n",len(g.Next))
//		g.PrintGuesser()
//		current.S.PrintSubBoard(b.Width,b.Height)
		for i:=0;i<len(b.S.Boxes)&&!goal;i++ {
			for dir:=0;dir<4&&!goal;dir++ {
				S := b.CheckMoveCostsForOneBox(i,direction.U+direction.Direction(dir))
				if S!=nil && S.MaxMoves<1000 { //&& currentDeep+S.MaxMoves<maxdeep  {
					if S.MaxMoves == 0 { 
						goal=true
						GoalSubBoard=S;
						break;
					}
//					fmt.Printf("Push board %s at %02d moves and %02d deep for %02d\n",S.GetString(),S.MaxMoves,currentDeep+1,maxdeep)
					g.Push(S,currentDeep+1);
//					S.PrintSubBoard(b.Width,b.Height) 
//					fmt.Printf("There is %02d subboards\n",len(g.Next))
				} // else if S!=nil { fmt.Printf("Skip board %s at %02d moves and %02d deep for %02d\n",S.GetString(),S.MaxMoves,currentDeep+1,maxdeep) 
//				}
			}
		}
	}
	b.S = saveSubBoard

	if goal {
//		fmt.Println("Goal !")
		var parent *SubBoard
//		GoalSubBoard.PrintSubBoard(b.Width,b.Height)
		for parent != saveSubBoard {
			parent = g.Checked[GoalSubBoard.ParentSubBoard.GetString()].S
			parent.BestSubBoardDir = GoalSubBoard.ParentSubBoardDir
			parent.BestSubBoardBox = GoalSubBoard.ParentSubBoardBox
//			parent.PrintSubBoard(b.Width,b.Height)
			GoalSubBoard = parent
		}
	}// else { fmt.Println("Fail !") }
}