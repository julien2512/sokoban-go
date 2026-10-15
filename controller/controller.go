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
		case pixelgl.KeyZ:
			c.tryUndoLastMove()
		case pixelgl.KeyR:
			c.restartLevel()
		case pixelgl.KeyTab:
			c.switchLayer()
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

	if (c.autoplay && c.m.State == model.StatePlaying) {
		pathDir := c.m.Board.GetNextDir()
		c.tryMovePlayer(pathDir)	
	}
}

func (c *Controller) toggleAutoplay() {
	if (c.autoplay) { 
		c.autoplay = false
		c.autoTime.Stop()
	} else {
		c.autoplay = true
		
		c.autoTime = time.NewTicker(250 * time.Millisecond)

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

func (c *Controller) switchLayer() {
	if c.m.Layer == model.LayerPlay {
		c.m.Layer = model.LayerBox
		c.m.BoxLayer = 0
	} else if c.m.Layer == model.LayerBox {
		if c.m.BoxLayer < len(c.m.Board.S.Boxes)-1 {
			c.m.BoxLayer++ 
		} else {
			c.m.Layer = model.LayerPath
		}
	} else {
		c.m.Layer = model.LayerPlay
	}
}

// tryMovePlayer - Move player (and an adjacent box where appropriate) in the specified direction if possible. Check for board completion (and handle appropriately) if a box is moved
func (c *Controller) tryMovePlayer(dir direction.Direction) {
	_,typ := c.m.Board.MovePlayer(dir,true)
	c.Guess()
	switch(typ) {
		case model.PlayerBlockedByWall :
			fmt.Printf("%v: Player blocked (wall)\n", dir)
		case model.PlayerBlockedByBox :
			fmt.Printf("%v: Box blocked (box)\n", dir)
		case model.PlayerMoveAndPush : 
			fmt.Printf("%v: Player moved to %02d %02d (push)\n", dir,c.m.Board.LastMove.LastTargetCell.X,c.m.Board.LastMove.LastTargetCell.Y)
			if c.m.Board.IsComplete() {
				c.m.State = model.StateLevelComplete
				fmt.Print("*** Level complete! ***\n(space key to continue)\n")
			}
		case model.PlayerMove :
			fmt.Printf("%v: Player moved (clear)\n", dir)
	}
}

func (c *Controller) Guess() {
	autoplay := c.autoplay
	if c.autoplay { c.toggleAutoplay() }
	c.m.Guess()
	if autoplay { c.toggleAutoplay() }	
}

func (c *Controller) tryUndoLastMove() {
	_,ret := c.m.Board.UndoLastMove()
	c.Guess()

	switch(ret) {
		case model.NoUndo :
			return
		case model.PlayerUndoMove : 
			fmt.Printf("Move back to %02d %02d\n",c.m.Board.S.Player.X,c.m.Board.S.Player.Y)
	}
	fmt.Printf("Player undo last move\n")
}

func (c *Controller) loadLevel() {
	autoplay := c.autoplay
	if c.autoplay { c.toggleAutoplay() }
	l := c.m.LM.GetCurrentLevel()
	c.m.Board = model.NewBoard(l.MapData, l.Width, l.Height)
	c.m.State = model.StatePlaying
	if autoplay { c.toggleAutoplay() }
}

// tryStartNextLevel - Starts the next level if the current one isn't the last, else sets game state to game complete
func (c *Controller) tryStartNextLevel() {
	if c.m.LM.HasNextLevel() {
		c.m.LM.ProgressToNextLevel()
		c.loadLevel()
		fmt.Printf("Start level %d\n", c.m.LM.GetCurrentLevelNumber())
		c.Guess()
	} else {
		c.m.State = model.StateGameComplete
		fmt.Print("*** GAME COMPLETE! ***\n(space key to restart)\n")
	}
}

// restartLevel - Resets the game board to the current level's starting state
func (c *Controller) restartLevel() {
	c.loadLevel()
	fmt.Printf("Restart level %d\n", c.m.LM.GetCurrentLevelNumber())
	c.Guess()
}
