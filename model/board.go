package model

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
}

type Board struct {
	Width, Height int
	Boxes []int
	Goals []int
	Distances [][]int
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
			index := (y*b.Width)+x
			code := string(mapData[index])
			cell := Cell{}
			switch code {
			case "@":
				b.Player = NewPlayer(x, y)
			case "$":
				cell.HasBox = true
				b.Boxes = append(b.Boxes,index)
			case ".":
				cell.TypeOf = CellTypeGoal
				b.Goals = append(b.Goals,index)
			case "#":
				cell.TypeOf = CellTypeWall
			case "+":
				cell.TypeOf = CellTypeGoal
				b.Player = NewPlayer(x, y)
			case "*":
				cell.TypeOf = CellTypeGoal
				cell.HasBox = true
				b.Boxes = append(b.Boxes,index)
				b.Goals = append(b.Goals,index)
			}
			b.Cells[(y*b.Width)+x] = cell
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