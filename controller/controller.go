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
	ShowFreeSpace bool
	autoplay bool
	autoTime *time.Ticker
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
		case pixelgl.KeyA:
			c.toggleAutoplay()
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

func (c *Controller) Autoplay() {
	player := c.m.Board.Player
	board  := c.m.Board
	cell := board.Get(player.X,player.Y)
	pathDir := cell.PathDir
	
	if (c.autoplay && c.m.State == model.StatePlaying && pathDir!=direction.None) {
		c.tryMovePlayer(pathDir)	
	}
}

func (c *Controller) toggleAutoplay() {
	if (c.autoplay) { 
		c.autoplay = false
		c.autoTime.Stop()
	} else {
		c.autoplay = true
		
		c.autoTime = time.NewTicker(500 * time.Millisecond)

		go func() {
			for {
				if (!c.autoplay) { return }
				select {
					case <-c.autoTime.C:
						c.Autoplay()
				}
			}
		}()
	}
}

// toggle show/hide Free Space
func (c *Controller) toggleShowFreeSpace() {
	if (c.ShowFreeSpace) {
		c.ShowFreeSpace = false
	} else {
		c.ShowFreeSpace = true
		c.m.Board = c.m.Board.GetBoard(c.m.Boards)
		c.m.Board.CheckEveryBoxMoveFromPlayer(c.m.Boards)
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
				c.m.Moves++
				c.m.Board = c.m.Board.MoveBoxAndCheck(targetX,targetY,dir,c.m.Boards)
				c.m.LastMove = model.NewLastMove(lastX,lastY,targetX,targetY,nextX,nextY,c.m.LastMove)
				fmt.Printf("%v: Player moved (push)\n", dir)
				c.m.Board.CheckEveryBoxMoveFromPlayer(c.m.Boards)
				if c.m.Board.IsComplete() {
					c.m.State = model.StateLevelComplete
					fmt.Print("*** Level complete! ***\n(space key to continue)\n")
				}
			}
		} else {
			c.m.Moves++
			c.m.LastMove = model.NewLastMove(lastX,lastY,-1,-1,-1,-1,c.m.LastMove)
			c.m.Board.Player.X = targetX
			c.m.Board.Player.Y = targetY
			c.m.Board.CheckEveryBoxMoveFromPlayer(c.m.Boards)
			fmt.Printf("%v: Player moved (clear)\n", dir)
		}
	}
}

func (c *Controller) tryUndoLastMove() {
	if c.m.LastMove == nil {
		return
	}
	c.m.Board = c.m.Board.Duplicate()
	c.m.Board.Player.X = c.m.LastMove.LastX
	c.m.Board.Player.Y = c.m.LastMove.LastY
	c.m.Moves--
	if c.m.LastMove.LastTargetX != -1 {
		lastCell := c.m.Board.Get(c.m.LastMove.LastTargetX,c.m.LastMove.LastTargetY)
		nextCell := c.m.Board.Get(c.m.LastMove.LastNextX,c.m.LastMove.LastNextY)
		lastCell.HasBox = true
		nextCell.HasBox = false
		lastCell.Box = nextCell.Box
		c.m.Board.Boxes[lastCell.Box].X = c.m.LastMove.LastTargetX
		c.m.Board.Boxes[lastCell.Box].Y = c.m.LastMove.LastTargetY
		c.m.Board = c.m.Board.GetBoard(c.m.Boards)
	}
	c.m.LastMove = c.m.LastMove.PreviousMove
	fmt.Printf("Player undo last moved\n")

	if (c.ShowFreeSpace) {
			c.m.Board.CheckEveryBoxMoveFromPlayer(c.m.Boards)
	}
}

func (c *Controller) loadLevel() {
	l := c.m.LM.GetCurrentLevel()
	c.m.Board = model.NewBoard(l.MapData, l.Width, l.Height)
	c.m.Boards = make(map[string]*model.Board)
	c.m.LastMove = nil
	c.m.Moves = 0
	start := time.Now()		
	c.m.Board.CheckEveryBoxMoveFromPlayer(c.m.Boards)
	c.m.SolveDuration = time.Now().Sub(start)
	c.m.BestMoves = c.m.Board.GetBestPosition().BestLength
	c.m.State = model.StatePlaying
}

// tryStartNextLevel - Starts the next level if the current one isn't the last, else sets game state to game complete
func (c *Controller) tryStartNextLevel() {
	if c.m.LM.HasNextLevel() {
		c.m.LM.ProgressToNextLevel()
		c.loadLevel()
		fmt.Printf("Start level %d\n", c.m.LM.GetCurrentLevelNumber())
	} else {
		c.m.State = model.StateGameComplete
		fmt.Print("*** GAME COMPLETE! ***\n(space key to restart)\n")
	}
}

// restartLevel - Resets the game board to the current level's starting state
func (c *Controller) restartLevel() {
	c.loadLevel()
	fmt.Printf("Restart level %d\n", c.m.LM.GetCurrentLevelNumber())
}
