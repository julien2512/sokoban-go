package view

import (
	"fmt"
	"image"
	_ "image/png"
	"io/ioutil"
	"os"

	pixel "github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/ext/text"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/TheInvader360/sokoban-go/model"
	"github.com/TheInvader360/sokoban-go/direction"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/colornames"
)

type spriteIndex int

const (
	SpritePlayer spriteIndex = iota
	SpriteBox
	SpriteGoal
	SpriteWall
	SpriteGoalAndPlayer
	SpriteGoalAndBox
	SpriteBoxRedCross
	SpriteLogo
	SpritePlayerInFreeSpace
	SpriteGoalInFreeSpace
	SpriteGoalAndPlayerInFreeSpace
	SpriteFreeSpace
	SpriteFreeSpaceBestPath
	SpriteGoalInFreeSpaceBestPath
	SpriteFree
	SpriteBoxGoUp
	SpriteBoxGoDown
	SpriteBoxGoLeft
	SpriteBoxGoRight
	SpriteBoxShallNotGoUp
	SpriteBoxShallNotGoDown
	SpriteBoxShallNotGoLeft
	SpriteBoxShallNotGoRight
	SpriteBoxShallGoUp
	SpriteBoxShallGoDown
	SpriteBoxShallGoLeft
	SpriteBoxShallGoRight
)

type View struct {
	m           *model.Model
	win         *opengl.Window
	scaleFactor float64
	text        *text.Text
	sprites     []*pixel.Sprite
}

// NewView - Creates a view
func NewView(m *model.Model, win *opengl.Window, scaleFactor float64) *View {
	fontFile, err := os.Open("assets/HackJack.ttf")
	if err != nil {
		panic(err)
	}
	defer fontFile.Close()
	fontBytes, err := ioutil.ReadAll(fontFile)
	if err != nil {
		panic(err)
	}
	font, err := truetype.Parse(fontBytes)
	if err != nil {
		panic(err)
	}
	face := truetype.NewFace(font, &truetype.Options{
		Size:              10 * scaleFactor,
		GlyphCacheEntries: 1,
	})
	atlas := text.NewAtlas(face, text.ASCII)
	text := text.New(pixel.V(0, 0), atlas)
	text.LineHeight = 11 * scaleFactor
	text.Color = colornames.White

	spritesheetFile, err := os.Open("assets/spritesheet.png")
	if err != nil {
		panic(err)
	}
	defer spritesheetFile.Close()
	image, _, err := image.Decode(spritesheetFile)
	if err != nil {
		panic(err)
	}
	pictureData := pixel.PictureDataFromImage(image)

	v := View{
		m:           m,
		win:         win,
		scaleFactor: scaleFactor,
		text:        text,
		sprites: []*pixel.Sprite{
			pixel.NewSprite(pictureData, pixel.R(float64(0), float64(64), float64(16), float64(80))),   // player
			pixel.NewSprite(pictureData, pixel.R(float64(16), float64(64), float64(32), float64(80))),  // box
			pixel.NewSprite(pictureData, pixel.R(float64(32), float64(64), float64(48), float64(80))),  // goal
			pixel.NewSprite(pictureData, pixel.R(float64(48), float64(64), float64(64), float64(80))),  // wall
			pixel.NewSprite(pictureData, pixel.R(float64(64), float64(64), float64(80), float64(80))),  // goal+player
			pixel.NewSprite(pictureData, pixel.R(float64(80), float64(64), float64(96), float64(80))),  // goal+box
			pixel.NewSprite(pictureData, pixel.R(float64(96), float64(64), float64(112), float64(80))),  // box red cross
			pixel.NewSprite(pictureData, pixel.R(float64(0), float64(80), float64(112), float64(128))), // logo
			pixel.NewSprite(pictureData, pixel.R(float64(0), float64(48), float64(16), float64(64))), // player in freespace
			pixel.NewSprite(pictureData, pixel.R(float64(32), float64(48), float64(48), float64(64))), // goal in freespace
			pixel.NewSprite(pictureData, pixel.R(float64(64), float64(48), float64(80), float64(64))), // goal+player in freespace
			pixel.NewSprite(pictureData, pixel.R(float64(16), float64(48), float64(32), float64(64))), // freespace
			pixel.NewSprite(pictureData, pixel.R(float64(80), float64(48), float64(96), float64(64))), // freespace best path
			pixel.NewSprite(pictureData, pixel.R(float64(96), float64(48), float64(112), float64(64))), // goal in freespace best path
			pixel.NewSprite(pictureData, pixel.R(float64(48), float64(48), float64(64), float64(64))),  // free
			pixel.NewSprite(pictureData, pixel.R(float64(32), float64(32), float64(48), float64(48))),  // box go up
			pixel.NewSprite(pictureData, pixel.R(float64(16), float64(32), float64(32), float64(48))),  // box go down
			pixel.NewSprite(pictureData, pixel.R(float64(48), float64(32), float64(64), float64(48))),  // box go left		
			pixel.NewSprite(pictureData, pixel.R(float64(0), float64(32), float64(16), float64(48))),   // box go right
			pixel.NewSprite(pictureData, pixel.R(float64(32), float64(16), float64(48), float64(32))),  // box shall not go up
			pixel.NewSprite(pictureData, pixel.R(float64(16), float64(16), float64(32), float64(32))),  // box shall not go down
			pixel.NewSprite(pictureData, pixel.R(float64(48), float64(16), float64(64), float64(32))),  // box shall not go left		
			pixel.NewSprite(pictureData, pixel.R(float64(0), float64(16), float64(16), float64(32))),   // box shall not go right
			pixel.NewSprite(pictureData, pixel.R(float64(32), float64(0), float64(48), float64(16))),  // box shall go up
			pixel.NewSprite(pictureData, pixel.R(float64(16), float64(0), float64(32), float64(16))),  // box shall go down
			pixel.NewSprite(pictureData, pixel.R(float64(48), float64(0), float64(64), float64(16))),  // box shall go left		
			pixel.NewSprite(pictureData, pixel.R(float64(0), float64(0), float64(16), float64(16))),   // box shall go right
		},
	}

	return &v
}

// Draw - Draws a graphical representation of the model's current state (called once per main game loop iteration)
func (v *View) Draw() {
	v.win.Clear(colornames.Black)

	v.drawLogoSprite()

	switch v.m.State {
	case model.StatePlaying:
		v.drawBoardLayer()
		v.printString(fmt.Sprintf("Level %02d of %02d", v.m.LM.GetCurrentLevelNumber(), v.m.LM.GetFinalLevelNumber()), 48, 7)
		v.printString(fmt.Sprintf("Boxes : %02d",len(v.m.Board.S.Boxes)),48,8)
		v.printString(fmt.Sprintf("MeanMoves : %02d",v.m.Board.S.MaxMoves),48,9)
		v.printString("---Controls---\n\nCursors:  Move\nZ:        Undo\nR:       Reset\nTab: SwitchLayer\nEscape:   Quit", 48, 14)
	case model.StateLevelComplete:
		v.drawBoardLayer()
		v.printString(fmt.Sprintf("Level %02d of %02d", v.m.LM.GetCurrentLevelNumber(), v.m.LM.GetFinalLevelNumber()), 48, 7)
		if v.m.TickAccumulator < 10 {
			v.printString("LEVEL COMPLETE", 48, 12)
		}
		v.printString("---Controls---\n\nSpace:    Next\n              \nEscape:   Quit", 48, 15)
	case model.StateGameComplete:
		v.printString("GAME COMPLETE!", 16, 10)
		if v.m.TickAccumulator < 10 {
			v.printString("CONGRATULATIONS!", 15, 12)
		} else {
			for y := 0; y <= 13; y++ {
				for x := 0; x <= 21; x++ {
					if x == 0 || x == 21 || y == 0 || y == 13 {
						v.drawBoardSprite(SpritePlayer, float64(x), float64(y), float64(1), float64(1))
					}
				}
			}
		}
		v.printString("---Controls---\n\nSpace: Restart\n              \nEscape:   Quit", 48, 15)
	}

	v.win.Update()
}

func (v *View) drawBoardDistances(goal int) {
		boardOffsetX := ((22 - v.m.Board.Width) / 2) + 1 
		boardOffsetY := ((14 - v.m.Board.Height) / 2) + 1 

		var bestBox int
		if len(v.m.Board.S.BestBoxes)>0 {
			bestBox = v.m.Board.S.BestBoxes[goal]
		} else { bestBox = -1 }

		for y := 0; y < v.m.Board.Height; y++ {
			for x := 0; x < v.m.Board.Width; x++ {
				cell := v.m.Board.Get(x, y)
				distance := v.m.Board.Distances[goal][y*v.m.Board.Width+x]

				switch cell.TypeOf {
				case model.CellTypeWall:
					v.drawBoardSprite(SpriteWall, float64(x), float64(y), float64(boardOffsetX), float64(boardOffsetY))
				default:
					if distance>0 {
						v.drawTextSprite(distance-1,float64(x),float64(y),float64(boardOffsetX),float64(boardOffsetY),cell.HasBox, cell.HasBox && cell.Box == bestBox) }
				}
			}
		}

}

func (v *View) drawArrowsDir(cell *model.Cell,x, y, boardOffsetX, boardOffsetY int, dir direction.Direction) {
	if cell.CanMove[dir] { 
		v.drawBoardSprite(SpriteBoxGoUp+spriteIndex(dir), float64(x), float64(y), float64(boardOffsetX), float64(boardOffsetY))
	}
}

func (v *View) drawArrows(cell *model.Cell,x, y, boardOffsetX, boardOffsetY int) {
	if !cell.HasBox { return }
	v.drawArrowsDir(cell,x,y,boardOffsetX,boardOffsetY,direction.U)
	v.drawArrowsDir(cell,x,y,boardOffsetX,boardOffsetY,direction.D)
	v.drawArrowsDir(cell,x,y,boardOffsetX,boardOffsetY,direction.L)
	v.drawArrowsDir(cell,x,y,boardOffsetX,boardOffsetY,direction.R)
}

func (v *View) drawBoardPath() {
	if v.m.State != model.StateGameComplete {
		boardOffsetX := ((22 - v.m.Board.Width) / 2) + 1
		boardOffsetY := ((14 - v.m.Board.Height) / 2) + 1
		for y := 0; y < v.m.Board.Height; y++ {
			for x := 0; x < v.m.Board.Width; x++ {
				cell := v.m.Board.Get(x, y)
				switch cell.TypeOf {
				case model.CellTypeNone:
					if cell.HasBox {
						if cell.IsDead { v.drawBoardSprite(SpriteBoxRedCross, float64(x), float64(y), float64(boardOffsetX), float64(boardOffsetY))
						} else { v.drawBoardSprite(SpriteBox, float64(x), float64(y), float64(boardOffsetX), float64(boardOffsetY))
						       v.drawArrows(cell,x,y,boardOffsetX,boardOffsetY)
						}
					} else if cell.IsFree {
						v.drawBoardSprite(SpriteFreeSpace, float64(x), float64(y), float64(boardOffsetX), float64(boardOffsetY))
					}
				case model.CellTypeGoal:
					if cell.HasBox {
						if cell.IsDead { v.drawBoardSprite(SpriteBoxRedCross, float64(x), float64(y), float64(boardOffsetX), float64(boardOffsetY))
						} else { v.drawBoardSprite(SpriteGoalAndBox, float64(x), float64(y), float64(boardOffsetX), float64(boardOffsetY))
							 v.drawArrows(cell,x,y,boardOffsetX,boardOffsetY)
						}
					} else if v.m.Board.S.Player.X == x && v.m.Board.S.Player.Y == y {
						if cell.IsFree {
							v.drawBoardSprite(SpriteGoalAndPlayerInFreeSpace, float64(x), float64(y), float64(boardOffsetX), float64(boardOffsetY))
						} else {
							v.drawBoardSprite(SpriteGoalAndPlayer, float64(x), float64(y), float64(boardOffsetX), float64(boardOffsetY))
						}
					} else {
						if cell.IsFree {
							v.drawBoardSprite(SpriteGoalInFreeSpace, float64(x), float64(y), float64(boardOffsetX), float64(boardOffsetY))
						} else {
							v.drawBoardSprite(SpriteGoal, float64(x), float64(y), float64(boardOffsetX), float64(boardOffsetY))
						}
					}
				case model.CellTypeWall:
					v.drawBoardSprite(SpriteWall, float64(x), float64(y), float64(boardOffsetX), float64(boardOffsetY))
				}
			}
		}
		v.drawBoardSprite(SpritePlayer, float64(v.m.Board.S.Player.X), float64(v.m.Board.S.Player.Y), float64(boardOffsetX), float64(boardOffsetY))
	}
}

func (v *View) drawBoard() {
	if v.m.State != model.StateGameComplete {
		boardOffsetX := ((22 - v.m.Board.Width) / 2) + 1
		boardOffsetY := ((14 - v.m.Board.Height) / 2) + 1
		for y := 0; y < v.m.Board.Height; y++ {
			for x := 0; x < v.m.Board.Width; x++ {
				cell := v.m.Board.Get(x, y)
				switch cell.TypeOf {
				case model.CellTypeNone:
					if cell.HasBox && v.m.Board.S.Cells[v.m.Board.S.Boxes[cell.Box]].HasBox {
						v.drawBoardSprite(SpriteBox, float64(x), float64(y), float64(boardOffsetX), float64(boardOffsetY))
					}
				case model.CellTypeGoal:
					if cell.HasBox && v.m.Board.S.Cells[v.m.Board.S.Boxes[cell.Box]].HasBox {
						v.drawBoardSprite(SpriteGoalAndBox, float64(x), float64(y), float64(boardOffsetX), float64(boardOffsetY))
					} else if v.m.Board.S.Player.X == x && v.m.Board.S.Player.Y == y {
						v.drawBoardSprite(SpriteGoalAndPlayer, float64(x), float64(y), float64(boardOffsetX), float64(boardOffsetY))
					} else {
						v.drawBoardSprite(SpriteGoal, float64(x), float64(y), float64(boardOffsetX), float64(boardOffsetY))
					}
				case model.CellTypeWall:
					v.drawBoardSprite(SpriteWall, float64(x), float64(y), float64(boardOffsetX), float64(boardOffsetY))
				}
			}
		}
		v.drawBoardSprite(SpritePlayer, float64(v.m.Board.S.Player.X), float64(v.m.Board.S.Player.Y), float64(boardOffsetX), float64(boardOffsetY))
	}
}

func (v *View) drawBoardLayer() {
	if v.m.Layer == model.LayerPlay {
		v.drawBoard()
	} else if v.m.Layer == model.LayerBox {
		v.drawBoardDistances(v.m.BoxLayer)
	} else if v.m.Layer == model.LayerPath {
		v.drawBoardPath()
	}
}

func (v *View) drawLogoSprite() {
	r := pixel.R(float64(384)*v.scaleFactor, v.win.Bounds().H()-float64(48)*v.scaleFactor, float64(496)*v.scaleFactor, v.win.Bounds().H()-float64(0)*v.scaleFactor)
	v.sprites[SpriteLogo].Draw(v.win, pixel.IM.ScaledXY(pixel.ZV, pixel.V(r.W()/v.sprites[SpriteLogo].Frame().W(), r.H()/v.sprites[SpriteLogo].Frame().H())).Moved(r.Center()))
}

func (v *View) drawBoardSprite(s spriteIndex, x, y, offsetX, offsetY float64) {
	r := pixel.R((offsetX+x)*16*v.scaleFactor, v.win.Bounds().H()-(offsetY+y+1)*16*v.scaleFactor, (offsetX+x+1)*16*v.scaleFactor, v.win.Bounds().H()-(offsetY+y)*16*v.scaleFactor)
	v.sprites[s].Draw(v.win, pixel.IM.ScaledXY(pixel.ZV, pixel.V(r.W()/v.sprites[s].Frame().W(), r.H()/v.sprites[s].Frame().H())).Moved(r.Center()))
}

func (v *View) drawTextSprite(distance int, x, y, offsetX, offsetY float64, box bool, bestbox bool) {
	r := pixel.R((offsetX+x)*16*v.scaleFactor-25, -19+v.win.Bounds().H()-(offsetY+y+1)*16*v.scaleFactor, (offsetX+x+1)*16*v.scaleFactor-25, -19+v.win.Bounds().H()-(offsetY+y)*16*v.scaleFactor)
	v.text.Clear()
	if bestbox { v.text.Color = colornames.Red
	} else if box { v.text.Color = colornames.Yellow
	} else { v.text.Color = colornames.White }
	v.text.WriteString(fmt.Sprintf("%d",distance))
	v.text.Draw(v.win, pixel.IM.ScaledXY(pixel.ZV, pixel.V(r.W()/v.text.Bounds().W(), r.H()/v.text.Bounds().H())).Moved(r.Center()))
}

// printString - prints the given string at screen position x,y (i.e. 0-63,0-22)
func (v *View) printString(s string, x, y int) {
	v.text.Clear()
	v.text.Color = colornames.White
	v.text.WriteString(s)
	v.text.Draw(v.win, pixel.IM.Moved(pixel.V(float64(x*8)*v.scaleFactor, (245-float64(y*11))*v.scaleFactor)))
}
