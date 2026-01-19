package model

type LastMove struct {
	LastX, LastY int
	LastTargetX, LastTargetY int
	LastNextX, LastNextY int

	PreviousMove *LastMove
}

// LastMove - Memorize the last move and effect on board
func NewLastMove(x, y int, lasttargetX, lasttargetY int, lastnextX, lastnextY int, lastMove *LastMove) *LastMove {
	return &LastMove{LastX: x, LastY: y, LastTargetX: lasttargetX, LastTargetY: lasttargetY, LastNextX: lastnextX, LastNextY: lastnextY, PreviousMove: lastMove}
}
