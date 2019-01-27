package effects

import (
	"math/rand"

	"github.com/brianewing/redshift/strip"
)

type GGJ struct {
	Level int

	// game state
	playerX, playerY       float64 // player location
	playerVelX, playerVelY float64 // player velocity (amount added to x and y each frame)

	playerDiedState int // used to flash color when player dies

	exitX, exitY int // location of the exit

	// game layers
	background      *ggjBackground
	gameOfLife      *GameOfLife
	gameOfLifeLayer *Layer

	// game controls
	ButtonLeft, ButtonRight, ButtonJump, ButtonSpawn *LatchValue
}

var EXIT_COLOR = strip.LED{0, 0, 255}
var ENEMY_COLOR = strip.LED{255, 255, 255}
var FLOOR_COLOR = strip.LED{255, 0, 0}
var PLAYER_COLOR = strip.LED{255, 255, 0}
var PLAYER_DIED_COLOR = strip.LED{200, 0, 100}

var MOVE_VELOCITY = 0.2
var JUMP_VELOCITY = 0.3

var DECELERATION_FACTOR = 0.89

var LEVEL_BG_COLORS = map[int]strip.LED{
	0:  strip.LED{18, 203, 196},
	2:  strip.LED{196, 229, 56},
	3:  strip.LED{18, 203, 196},
	4:  strip.LED{253, 167, 223},
	5:  strip.LED{255, 0, 0},
	6:  strip.LED{255, 0, 0},
	7:  strip.LED{255, 0, 0},
	8:  strip.LED{255, 0, 0},
	9:  strip.LED{255, 0, 0},
	10: strip.LED{255, 0, 0},
}

func NewGGJ() *GGJ {
	return &GGJ{
		ButtonLeft:  &LatchValue{},
		ButtonRight: &LatchValue{},
		ButtonJump:  &LatchValue{},
		background:  &ggjBackground{},
		gameOfLife:  NewGameOfLife(),
	}
}

func (e *GGJ) Init(s *strip.LEDStrip) {
	e.resetLevel(s)
	e.gameOfLife.Init()
	e.gameOfLifeLayer = NewLayer()
	e.gameOfLifeLayer.Effects = EffectSet{EffectEnvelope{Effect: &Clear{}}, EffectEnvelope{Effect: e.gameOfLife}}
}

func (e *GGJ) Render(s *strip.LEDStrip) {
	if int(e.playerX) == e.exitX && int(e.playerY) == e.exitY {
		e.Level += 1
		e.resetLevel(s)
	}

	e.performMovement(s)

	e.background.Render(s)
	e.gameOfLife.Render(s)

		if e.Level > 1 {
			e.Level = 0
		}
		e.playerDiedState = 10
		e.resetLevel(s)
		return
	}

	e.drawExit(s)
	e.drawPlayer(s)
	e.drawFlash(s)

	e.handleInput(s)
}

func (e *GGJ) handleInput(s *strip.LEDStrip) {
	if e.playerDiedState > 0 {
		return
	}
	if e.ButtonJump.Read() && e.playerY > 0 {
		e.playerVelY = -JUMP_VELOCITY
	}
	if e.ButtonLeft.Read() && e.playerX > 0 {
		e.playerVelX = -MOVE_VELOCITY
	}
	if e.ButtonRight.Read() && s.Width > int(e.playerX) {
		e.playerVelX = MOVE_VELOCITY
	}
}

func (e *GGJ) performMovement(s *strip.LEDStrip) {
	if e.playerX >= 0 && int(e.playerX+e.playerVelX) <= s.Width-1 {
		e.playerX += e.playerVelX
	}
	if e.playerY >= 0 && int(e.playerY+e.playerVelY) <= s.Height-1 {
		e.playerY += e.playerVelY
	}

	e.playerVelX *= DECELERATION_FACTOR
	e.playerVelY *= DECELERATION_FACTOR

	// wrap around the grid
	// e.playerX = math.Abs(math.Mod(e.playerX, float64(s.Width)))
	// e.playerY = math.Abs(math.Mod(e.playerY, float64(s.Height)))
}

func (e *GGJ) resetLevel(s *strip.LEDStrip) {
	e.playerX = 0
	e.playerY = float64(s.Height - 1)
	e.playerVelX = 0
	e.playerVelY = 0

	exitPosition := rand.Intn(s.Size)
	e.exitX = exitPosition % s.Width
	e.exitY = int(exitPosition / s.Width)

	// e.background.Color = blendHcl(strip.LED{0, 0, 0}, LEVEL_BG_COLORS[e.Level], float64(3+(7-e.Level/7)))
	e.background.Color = LEVEL_BG_COLORS[e.Level]
	e.gameOfLife.N = 10
	// e.gameOfLife.N = 40 - e.Level*2
	e.gameOfLife.StartingCells = (s.Width * s.Height / 4) /// (10 - e.Level) * 3)
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

type ggjBackground struct {
	Color strip.LED
	layer *Layer
}

func (e *ggjBackground) Render(s *strip.LEDStrip) {
	(&Fill{Color: e.Color}).Render(s)
	return
	if e.layer == nil {
		e.layer = NewLayer()
		e.layer.Effects = EffectSet{
			EffectEnvelope{Effect: NewFill()},
			EffectEnvelope{
				Effect: NewBrightness(),
				Controls: ControlSet{
					ControlEnvelope{
						Control: &TweenControl{
							BaseControl: BaseControl{Field: "Level"},
							Function:    "sin", Min: 200, Max: 245, Speed: 0.5,
						},
					},
				},
			},
		}
	}

	fill, _ := e.layer.Effects[0].Effect.(*Fill)
	fill.Color = e.Color

	e.layer.Render(s)
}

func compareColors(a, b strip.LED) bool {
	return len(a) == 3 && len(b) == 3 && a[0] == b[0] && a[1] == b[1] && a[2] == b[2]
}

// got you, you died, lol, oh no, oh well, bye bye
// welcome home
