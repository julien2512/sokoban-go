package model

type LastMove struct {
	LastX, LastY int
	LastTargetCell, LastNextCell *Cell

	PreviousMove *LastMove
}

// LastMove - Memorize the last move and effect on board
func NewLastMove(x, y int, targetCell, nextCell *Cell, lastMove *LastMove) *LastMove {
	return &LastMove{LastX: x, LastY: y, LastTargetCell: targetCell, LastNextCell: nextCell, PreviousMove: lastMove}
}
