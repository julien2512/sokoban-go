package controller

import (
	"fmt"
	"time"

	pixelgl "github.com/gopxl/pixel/v2"
	"github.com/TheInvader360/sokoban-go/direction"
	"github.com/TheInvader360/sokoban-go/model"
)

type Controller struct {
	m *model.Model
	rewind *time.Ticker
}

// NewController - Creates a controller
func NewController(m *model.Model) *Controller {
	c := Controller{
		m: m,
	}

	return &c
}

// StartNewGame - Starts a new game at level 1
func (c *Controller) StartNewGame() {
	c.m.LM.Reset()
	c.tryStartNextLevel()
}

// HandleInput - Handles user input as appropriate (game state dependent behaviour)
func (c *Controller) HandleInput(key pixelgl.Button) {
	switch c.m.State {
	case model.StateRewind:
		switch key {
		case pixelgl.KeyR:
			c.toggleRewind()
		}
	case model.StatePlaying:
		switch key {
		case pixelgl.KeyUp:
			c.tryMovePlayer(direction.U)
		case pixelgl.KeyDown:
			c.tryMovePlayer(direction.D)
		case pixelgl.KeyLeft:
			c.tryMovePlayer(direction.L)
		case pixelgl.KeyRight:
			c.tryMovePlayer(direction.R)
		case pixelgl.KeyTab:
			c.trySwitchPlayer()
		case pixelgl.KeyZ:
			c.tryUndoLastMove()
		case pixelgl.KeyR:
			c.toggleRewind()
		}
	case model.StateLevelComplete:
		if key == pixelgl.KeySpace {
			c.tryStartNextLevel()
		}
	case model.StateGameComplete:
		if key == pixelgl.KeySpace {
			c.StartNewGame()
		}
	}
}


func (c *Controller) Rewind() {
	if (c.m.State == model.StateRewind) {
		if c.m.Board.LastMove == nil {
			c.toggleRewind()
		} else {
			c.tryUndoLastMove()
		}
	}
}

func (c *Controller) toggleRewind() {
	if (c.m.State == model.StateRewind) { 
		c.m.State = model.StatePlaying
		c.rewind.Stop()
	} else {
		c.m.State = model.StateRewind
		
		c.rewind = time.NewTicker(350 * time.Millisecond)

		go func() {
			for {
				if (c.m.State != model.StateRewind) { return }
				select {
					case <-c.rewind.C:
						c.Rewind()
				}
			}
		}()
	}
}


// tryMovePlayer - Move player (and an adjacent box where appropriate) in the specified direction if possible. Check for board completion (and handle appropriately) if a box is moved
func (c *Controller) tryMovePlayer(dir direction.Direction) {
	lastX := c.m.Board.Players[c.m.Board.ActivePlayer].X
	lastY := c.m.Board.Players[c.m.Board.ActivePlayer].Y
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

	targetCell := c.m.Board.Get(targetX, targetY)
	TargetHasPlayer,_ := c.m.Board.CheckPlayer(targetX,targetY)

	if targetCell.TypeOf == model.CellTypeWall {
		fmt.Printf("%v: Player blocked (wall)\n", dir)
	} else if TargetHasPlayer {
		fmt.Printf("%v: Player blocked (player)\n", dir)
	} else if targetCell.HasBox {
		nextCell := c.m.Board.Get(nextX, nextY)
		NextHasPlayer,_ := c.m.Board.CheckPlayer(nextX,nextY)
		if nextCell.TypeOf == model.CellTypeWall {
			fmt.Printf("%v: Box blocked (wall)\n", dir)
		} else if nextCell.HasBox {
			fmt.Printf("%v: Box blocked (box)\n", dir)
		} else if NextHasPlayer {
			fmt.Printf("%v: Box blocked (player)\n", dir)
		} else {
			c.m.Board.LastMove = model.NewLastMove(lastX,lastY,targetCell,nextCell,c.m.Board.ActivePlayer,c.m.Board.LastMove)
			targetCell.HasBox = false
			nextCell.HasBox = true
			c.m.Board.Players[c.m.Board.ActivePlayer].X = targetX
			c.m.Board.Players[c.m.Board.ActivePlayer].Y = targetY
			fmt.Printf("%v: Player moved (push)\n", dir)
			if c.m.Board.IsComplete() {
				c.m.State = model.StateLevelComplete
				fmt.Print("*** Level complete! ***\n(space key to continue)\n")
			}
		}
	} else {
			c.m.Board.LastMove = model.NewLastMove(lastX,lastY,nil,nil,c.m.Board.ActivePlayer,c.m.Board.LastMove)
			c.m.Board.Players[c.m.Board.ActivePlayer].X = targetX
			c.m.Board.Players[c.m.Board.ActivePlayer].Y = targetY
			fmt.Printf("%v: Player moved (clear)\n", dir)
	}
}

func (c *Controller) trySwitchPlayer() {
	if c.m.Board.ActivePlayer+1 == len(c.m.Board.Players) {
		c.m.Board.ActivePlayer = 0
	} else {
		c.m.Board.ActivePlayer++
	}
}

func (c *Controller) tryUndoLastMove() {
	if c.m.Board.LastMove == nil {
		return
	}
	c.m.Board.Players[c.m.Board.LastMove.LastActivePlayer].X = c.m.Board.LastMove.LastX
	c.m.Board.Players[c.m.Board.LastMove.LastActivePlayer].Y = c.m.Board.LastMove.LastY
	c.m.Board.ActivePlayer = c.m.Board.LastMove.LastActivePlayer
	if c.m.Board.LastMove.LastTargetCell != nil {
		c.m.Board.LastMove.LastTargetCell.HasBox = true
		c.m.Board.LastMove.LastNextCell.HasBox = false
	}
	c.m.Board.LastMove = c.m.Board.LastMove.PreviousMove
	fmt.Printf("Player undo last moved\n")
}

// tryStartNextLevel - Starts the next level if the current one isn't the last, else sets game state to game complete
func (c *Controller) tryStartNextLevel() {
	if c.m.LM.HasNextLevel() {
		c.m.LM.ProgressToNextLevel()
		l := c.m.LM.GetCurrentLevel()
		c.m.Board = model.NewBoard(l.MapData, l.Width, l.Height)
		c.m.State = model.StatePlaying
		fmt.Printf("Start level %d\n", c.m.LM.GetCurrentLevelNumber())
	} else {
		c.m.State = model.StateGameComplete
		fmt.Print("*** GAME COMPLETE! ***\n(space key to restart)\n")
	}
}

// restartLevel - Resets the game board to the current level's starting state
func (c *Controller) restartLevel() {
	l := c.m.LM.GetCurrentLevel()
	c.m.Board = model.NewBoard(l.MapData, l.Width, l.Height)
	c.m.State = model.StatePlaying
	fmt.Printf("Restart level %d\n", c.m.LM.GetCurrentLevelNumber())
}
