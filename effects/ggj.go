package effects

import (
	"math/rand"

	"github.com/brianewing/redshift/strip"
)

var EXIT_COLOR = strip.LED{0, 0, 255}
var ENEMY_COLOR = strip.LED{255, 255, 255}
var PLAYER_COLOR = strip.LED{255, 230, 110}
var PLAYER_DIED_COLOR = strip.LED{210, 0, 50}
var LEVEL_INDICATOR_COLOR = strip.LED{255, 255, 255}

// var DECELERATION_FACTOR = 0.95
var DECELERATION_FACTOR = 0.0

type gglLevel struct {
	Color       strip.LED
	GameOfLifeN int
	Velocity    float64
}

var LEVELS = []gglLevel{
	{strip.LED{30, 30, 30}, 999999999999, 0.25},
	{strip.LED{18, 203, 196}, 45, 0.3},
	{strip.LED{253, 167, 223}, 30, 0.33},
	{strip.LED{255, 123, 50}, 25, 0.36},
	{strip.LED{50, 123, 255}, 20, 0.40},
	{strip.LED{255, 40, 190}, 15, 0.44},
	{strip.LED{255, 126, 40}, 13, 0.51},
	{strip.LED{140, 255, 60}, 9, 0.54},
	{strip.LED{118, 36, 100}, 7, 0.58},
	{strip.LED{70, 90, 98}, 6, 0.6},
	{strip.LED{19, 130, 141}, 5, 0.62},
	{strip.LED{8, 30, 127}, 4, 0.64},
	{strip.LED{136, 30, 80}, 3, 0.67},
	{strip.LED{250, 70, 50}, 2, 0.7},
	{strip.LED{30, 198, 41}, 1, 0.72},
}

// GGJ is a game created for Global Game Jam 2019 at Farset Labs
type GGJ struct {
	Level        int
	EnemyEffects EffectSet // set to replace default enemy

	// Controller buttons
	ButtonLeft, ButtonRight, ButtonJump, ButtonDown *LatchValue

	// Game state
	playerX, playerY       float64 // player location
	playerVelX, playerVelY float64 // player velocity (amount added to x and y each frame)

	playerDiedState int  // used to flash color when player dies
	playerHasMoved  bool // resets with level
	playerHealth    uint8

	exitX, exitY int // location of the exit

	// Layers
	background *ggjBackground
	gameOfLife *GameOfLife
	winState   Effect
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
	// e.Level = 0
	e.playerHealth = 255
	e.resetLevel(s)
	// for i, _ := range LEVELS {
	// LEVELS[i].Velocity = LEVELS[i].Velocity / 1.1
	// }
}

func (e *GGJ) Render(s *strip.LEDStrip) {
	if e.Level >= len(LEVELS) {
		e.drawWinState(s)
		return
	}

	if e.isPlayerAtExit(s) {
		e.nextLevel(s)
		return
	}

	e.handleInput(s)
	e.performMovement(s)

	e.background.Render(s)

	if len(e.EnemyEffects) == 0 {
		e.gameOfLife.Render(s)
	} else {
		e.EnemyEffects.Render(s)
	}

	if e.isCollidingWithEnemy(s) && e.playerHasMoved && e.playerDiedState == 0 {
		e.playerDiedState = 20
		e.prevLevel(s)
		return
	}

	e.drawExit(s)
	e.drawLevelIndicator(s)
	e.drawPlayer(s)
	e.drawFlash(s)
}

func (e *GGJ) handleInput(s *strip.LEDStrip) {
	if e.playerDiedState > 0 {
		return
	}

	velocity := LEVELS[e.Level].Velocity

	if e.ButtonJump.Read() {
		e.playerVelY += -velocity
	}
	if e.ButtonDown.Read() {
		e.playerVelY += velocity
	}
	if e.ButtonLeft.Read() {
		e.playerVelX += -velocity
	}
	if e.ButtonRight.Read() {
		e.playerVelX += velocity
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

	if e.playerVelX > 0.05 || e.playerVelY > 0.05 {
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
	}
	e.resetLevel(s)
}

func (e *GGJ) resetLevel(s *strip.LEDStrip) {
	if len(LEVELS) <= e.Level {
		return
	}

	e.playerX, e.playerY = 0, float64(s.Height-1)
	e.playerVelX, e.playerVelY = 0, 0

	e.playerHasMoved = false

	e.exitX, e.exitY = rand.Intn(s.Width), rand.Intn(s.Height)

	e.background.Color = LEVELS[e.Level].Color

	e.gameOfLife.StartingCells = (s.Width * s.Height / 7)
	e.gameOfLife.N = LEVELS[e.Level].GameOfLifeN
	e.gameOfLife.Init()
}

func (e *GGJ) drawExit(s *strip.LEDStrip) {
	s.SetXY(e.exitX, e.exitY, EXIT_COLOR)
}

func (e *GGJ) drawLevelIndicator(s *strip.LEDStrip) {
	for i := 0; i <= e.Level; i++ {
		s.SetPixel(i, blendHcl(s.Buffer[i], LEVEL_INDICATOR_COLOR, 0.4))
	}
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

func (e *GGJ) drawWinState(s *strip.LEDStrip) {
	if e.winState == nil {
		rainbow := NewRainbow()
		rainbow.Size = uint16(s.Size * 4)
		rainbow.Speed = 0.6
		mirror := NewMirror()
		mirror.Effects = EffectSet{EffectEnvelope{Effect: rainbow}}
		e.winState = mirror
		// mood := NewMoodEffect()
		// mood.Speed = 2
		// mood.Init()
		// e.winState = mood
	}
	e.winState.Render(s)
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
		e.layer.Blend.Factor = 0.03
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
