package effects

import (
	"math/rand"

	"github.com/brianewing/redshift/strip"
)

var EXIT_COLOR = strip.LED{0, 0, 255}
var ENEMY_COLOR = strip.LED{255, 255, 255}
var FLOOR_COLOR = strip.LED{255, 0, 0}
var PLAYER_COLOR = strip.LED{255, 255, 0}
var PLAYER_DIED_COLOR = strip.LED{200, 0, 100}

var MOVE_VELOCITY = 0.4
var JUMP_VELOCITY = 0.4

var DECELERATION_FACTOR = 0.65

var LEVEL_BG_COLORS = map[int]strip.LED{
	0:  strip.LED{30, 30, 30},
	1:  strip.LED{18, 203, 196},
	2:  strip.LED{196, 229, 56},
	3:  strip.LED{18, 203, 196},
	4:  strip.LED{253, 167, 223},
	5:  strip.LED{255, 123, 50},
	6:  strip.LED{50, 123, 255},
	7:  strip.LED{255, 40, 190},
	8:  strip.LED{255, 190, 40},
	9:  strip.LED{140, 255, 60},
	10: strip.LED{16, 36, 100},
}

// GGJ is a game created for Global Game Jam 2019 at Farset Labs
type GGJ struct {
	Level int

	// Controller buttons
	ButtonLeft, ButtonRight, ButtonJump, ButtonDown *LatchValue

	// Game state
	playerX, playerY       float64 // player location
	playerVelX, playerVelY float64 // player velocity (amount added to x and y each frame)

	playerDiedState int  // used to flash color when player dies
	playerHasMoved  bool // resets with level

	exitX, exitY int // location of the exit

	// Layers
	background *ggjBackground
	gameOfLife *GameOfLife
}

func NewGGJ() *GGJ {
	return &GGJ{
		ButtonLeft:  &LatchValue{},
		ButtonRight: &LatchValue{},
		ButtonJump:  &LatchValue{},
		ButtonDown:  &LatchValue{},
		background:  &ggjBackground{},
		gameOfLife:  NewGameOfLife(),
	}
}

func (e *GGJ) Init(s *strip.LEDStrip) {
	e.Level = 0
	e.resetLevel(s)
	e.gameOfLife.Init()
}

func (e *GGJ) Render(s *strip.LEDStrip) {
	if e.isPlayerAtExit(s) {
		e.nextLevel(s)
	}

	e.handleInput(s)
	e.performMovement(s)

	e.background.Render(s)
	e.gameOfLife.Render(s)

	if e.isCollidingWithEnemy(s) && e.playerHasMoved {
		e.playerDiedState = 10
		e.prevLevel(s)
		return
	}

	e.drawExit(s)
	e.drawPlayer(s)
	e.drawFlash(s)
}

func (e *GGJ) handleInput(s *strip.LEDStrip) {
	if e.playerDiedState > 0 {
		return
	}
	if e.ButtonJump.Read() {
		e.playerVelY = -JUMP_VELOCITY
	}
	if e.ButtonDown.Read() {
		e.playerVelY = JUMP_VELOCITY
	}
	if e.ButtonLeft.Read() {
		e.playerVelX = -MOVE_VELOCITY
	}
	if e.ButtonRight.Read() {
		e.playerVelX = MOVE_VELOCITY
	}
}

func (e *GGJ) performMovement(s *strip.LEDStrip) {
	e.playerX += e.playerVelX
	e.playerY += e.playerVelY

	if e.playerX < 0 {
		e.playerX += float64(s.Width)
	}
	if e.playerY < 0 {
		e.playerY += float64(s.Height)
	}
	if e.playerX > float64(s.Width)-0.5 {
		e.playerX -= float64(s.Width)
	}
	if e.playerY > float64(s.Height)-0.5 {
		e.playerY -= float64(s.Height)
	}

	if e.playerVelX != 0 || e.playerVelY != 0 {
		e.playerHasMoved = true
	}

	e.playerVelX *= DECELERATION_FACTOR
	e.playerVelY *= DECELERATION_FACTOR

	// wrap around the grid
	// e.playerX = math.Abs(math.Mod(e.playerX, float64(s.Width)))
	// e.playerY = math.Abs(math.Mod(e.playerY, float64(s.Height)))
}

func (e *GGJ) nextLevel(s *strip.LEDStrip) {
	e.Level += 1
	e.resetLevel(s)
}

func (e *GGJ) prevLevel(s *strip.LEDStrip) {
	if e.Level >= 1 {
		e.Level -= 1
		e.resetLevel(s)
	}
}

func (e *GGJ) resetLevel(s *strip.LEDStrip) {
	e.playerX, e.playerY = 0, float64(s.Height-1)
	e.playerVelX, e.playerVelY = 0, 0

	e.playerHasMoved = false

	e.exitX, e.exitY = rand.Intn(s.Width), rand.Intn(s.Height)

	e.background.Color = LEVEL_BG_COLORS[e.Level]

	e.gameOfLife.StartingCells = (s.Width * s.Height / 6)
	e.gameOfLife.N = 61 - (e.Level * 6)
}

func (e *GGJ) drawExit(s *strip.LEDStrip) {
	s.SetXY(e.exitX, e.exitY, EXIT_COLOR)
}

func (e *GGJ) drawPlayer(s *strip.LEDStrip) {
	s.SetXY(int(e.playerX), int(e.playerY), PLAYER_COLOR)
}

func (e *GGJ) drawFlash(s *strip.LEDStrip) {
	if e.playerDiedState%3 >= 1 {
		(&Fill{Color: PLAYER_DIED_COLOR}).Render(s)
	}
	if e.playerDiedState >= 1 {
		e.playerDiedState--
	}
}

func (e *GGJ) isCollidingWithEnemy(s *strip.LEDStrip) bool {
	enemy := s.GetXY(int(e.playerX), int(e.playerY))
	return compareColors(enemy, ENEMY_COLOR)
}

func (e *GGJ) isPlayerAtExit(s *strip.LEDStrip) bool {
	return int(e.playerX) == e.exitX && int(e.playerY) == e.exitY
}

type ggjBackground struct {
	Color strip.LED
	layer *Layer
}

func (e *ggjBackground) Render(s *strip.LEDStrip) {
	// (&Fill{Color: e.Color}).Render(s)
	// return
	if e.layer == nil {
		e.layer = NewLayer()
		e.layer.Blend.Factor = 0.05
		gol2 := NewGameOfLife()
		gol2.Color = strip.LED{0, 255, 0}
		e.layer.Effects = EffectSet{
			EffectEnvelope{Effect: NewFill()},
			// EffectEnvelope{Effect: gol2},
		}

		e.layer.Init()

		if rainbow, ok := e.layer.Effects[0].Effect.(*Rainbow); ok {
			rainbow.Size = uint16(s.Size * 3)
			rainbow.Reverse = true
			rainbow.Blend.Factor = 0.01
		}
	}

	if fill, ok := e.layer.Effects[0].Effect.(*Fill); ok {
		fill.Color = e.Color
	}

	e.layer.Render(s)
}

func compareColors(a, b strip.LED) bool {
	return len(a) == 3 && len(b) == 3 && a[0] == b[0] && a[1] == b[1] && a[2] == b[2]
}

// got you, you died, lol, oh no, oh well, bye bye
// welcome home
