package model

type LastMove struct {
	LastX, LastY int
	LastTargetCell, LastNextCell *Cell
	LastActivePlayer int

	PreviousMove *LastMove
}

// LastMove - Memorize the last move and effect on board
func NewLastMove(x, y int, targetCell, nextCell *Cell, lastActivePlayer int, lastMove *LastMove) *LastMove {
	return &LastMove{LastX: x, LastY: y, LastTargetCell: targetCell, LastNextCell: nextCell, LastActivePlayer: lastActivePlayer, PreviousMove: lastMove}
}
