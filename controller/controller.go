package controller

import (
	"fmt"

	pixelgl "github.com/gopxl/pixel/v2"
	"github.com/TheInvader360/sokoban-go/direction"
	"github.com/TheInvader360/sokoban-go/model"
)

type Controller struct {
	m *model.Model
	ShowFreeSpace bool
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
		case pixelgl.KeyF:
			c.toggleShowFreeSpace()
		case pixelgl.KeyZ:
			c.tryUndoLastMove()
		case pixelgl.KeyR:
			c.restartLevel()
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

// toggle show/hide Free Space
func (c *Controller) toggleShowFreeSpace() {
	if (c.ShowFreeSpace) {
		c.ShowFreeSpace = false
		c.m.Board.ResetFreeSpace()
	} else {
		c.ShowFreeSpace = true
		c.m.Board.CheckEveryFreeSpaceFromPlayer(c.m.Boards)
	}
}

// tryMovePlayer - Move player (and an adjacent box where appropriate) in the specified direction if possible. Check for board completion (and handle appropriately) if a box is moved
func (c *Controller) tryMovePlayer(dir direction.Direction) {
	lastX := c.m.Board.Player.X
	lastY := c.m.Board.Player.Y
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

	if targetCell.TypeOf == model.CellTypeWall {
		fmt.Printf("%v: Player blocked (wall)\n", dir)
	} else {
		if targetCell.HasBox {
			nextCell := c.m.Board.Get(nextX, nextY)
			if nextCell.TypeOf == model.CellTypeWall {
				fmt.Printf("%v: Box blocked (wall)\n", dir)
			} else if nextCell.HasBox {
				fmt.Printf("%v: Box blocked (box)\n", dir)
			} else {
				c.m.Board.LastMove = model.NewLastMove(lastX,lastY,targetX,targetY,nextX,nextY,c.m.Board.LastMove)
				targetCell.HasBox = false
				nextCell.HasBox = true
				c.m.Board.Player.X = targetX
				c.m.Board.Player.Y = targetY
				fmt.Printf("%v: Player moved (push)\n", dir)
				if c.m.Board.IsComplete() {
					c.m.State = model.StateLevelComplete
					fmt.Print("*** Level complete! ***\n(space key to continue)\n")
				}
				if (c.ShowFreeSpace) {
					go func() {
						c.m.Board.CheckEveryFreeSpaceFromPlayer(c.m.Boards)
					}()
				}
			}
		} else {
			c.m.Board.LastMove = model.NewLastMove(lastX,lastY,-1,-1,-1,-1,c.m.Board.LastMove)
			c.m.Board.Player.X = targetX
			c.m.Board.Player.Y = targetY
			fmt.Printf("%v: Player moved (clear)\n", dir)
		}
	}
}

func (c *Controller) tryUndoLastMove() {
	if c.m.Board.LastMove == nil {
		return
	}
	c.m.Board.Player.X = c.m.Board.LastMove.LastX
	c.m.Board.Player.Y = c.m.Board.LastMove.LastY
	if c.m.Board.LastMove.LastTargetX != -1 {
		c.m.Board.Get(c.m.Board.LastMove.LastTargetX,c.m.Board.LastMove.LastTargetY).HasBox = true
		c.m.Board.Get(c.m.Board.LastMove.LastNextX,c.m.Board.LastMove.LastNextY).HasBox = false
	}
	c.m.Board.LastMove = c.m.Board.LastMove.PreviousMove
	fmt.Printf("Player undo last moved\n")

	if (c.ShowFreeSpace) {
			c.m.Board.CheckEveryFreeSpaceFromPlayer(c.m.Boards)
	}
}

// tryStartNextLevel - Starts the next level if the current one isn't the last, else sets game state to game complete
func (c *Controller) tryStartNextLevel() {
	if c.m.LM.HasNextLevel() {
		c.m.LM.ProgressToNextLevel()
		l := c.m.LM.GetCurrentLevel()
		c.m.Board = model.NewBoard(l.MapData, l.Width, l.Height)
		c.m.Boards = make(map[string]*model.Board)
		c.m.State = model.StatePlaying
		if (c.ShowFreeSpace) {
			c.m.Board.CheckEveryFreeSpaceFromPlayer(c.m.Boards)
		}
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
	c.m.Boards = make(map[string]*model.Board)
	c.m.State = model.StatePlaying
	if (c.ShowFreeSpace) {
			c.m.Board.CheckEveryFreeSpaceFromPlayer(c.m.Boards)
	}
	fmt.Printf("Restart level %d\n", c.m.LM.GetCurrentLevelNumber())
}
