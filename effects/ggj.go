package effects

import (
	"github.com/brianewing/redshift/strip"
)

type GGJ struct {
	// game state
	playerX, playerY int
	playerDiedState  int // used to flash color when player dies

	exitX, exitY int

	Level int

	gameOfLife *GameOfLife

	// game controls
	ButtonLeft, ButtonRight, ButtonJump *LatchValue
}

var EXIT_COLOR = strip.LED{0, 0, 255}
var ENEMY_COLOR = strip.LED{255, 255, 255}
var FLOOR_COLOR = strip.LED{255, 0, 0}
var PLAYER_COLOR = strip.LED{255, 255, 0}
var PLAYER_DIED_COLOR = strip.LED{200, 0, 200}

func NewGGJ() *GGJ {
	return &GGJ{
		ButtonLeft:  &LatchValue{},
		ButtonRight: &LatchValue{},
		ButtonJump:  &LatchValue{},
		gameOfLife:  NewGameOfLife(),
	}
}

func (e *GGJ) Init(s *strip.LEDStrip) {
	e.resetLevel(s)
	e.gameOfLife.Init()
}

func (e *GGJ) Render(s *strip.LEDStrip) {
	if e.playerX == e.exitX && e.playerY == e.exitY {
		e.Level += 1
		e.resetLevel(s)
	}

	if e.isCollidingWithEnemy(s) {
		e.playerDiedState = 60
		e.resetLevel(s)
	}

	e.drawFloor(s)
	e.drawExit(s)
	e.drawPlayer(s)

	e.handleInput(s)
}

func (e *GGJ) handleInput(s *strip.LEDStrip) {
	if e.ButtonJump.Read() && e.playerY > 0 {
		e.playerY -= 1
	}
	if e.ButtonLeft.Read() && e.playerX > 0 {
		e.playerX -= 1
	}
	if e.ButtonRight.Read() && s.Width > e.playerX {
		e.playerX += 1
	}
}

func (e *GGJ) drawFloor(s *strip.LEDStrip) {
	for x := 0; x < s.Width; x++ {
		s.SetXY(x, s.Height-2, FLOOR_COLOR)
		s.SetXY(x, s.Height-1, strip.LED{0, 0, 0})
	}
}

func (e *GGJ) drawExit(s *strip.LEDStrip) {
	s.SetXY(e.exitX, e.exitY, EXIT_COLOR)
}

func (e *GGJ) drawPlayer(s *strip.LEDStrip) {
	color := PLAYER_COLOR
	if e.playerDiedState%3 >= 1 {
		color = PLAYER_DIED_COLOR
	}
	if e.playerDiedState > 1 {
		e.playerDiedState--
	}
	s.SetXY(e.playerX, e.playerY, color)
}

func (e *GGJ) resetLevel(s *strip.LEDStrip) {
	e.playerX = 0
	e.playerY = s.Height - 3
	e.exitX = s.Width - 1
	e.exitY = 0
}

func (e *GGJ) isCollidingWithEnemy(s *strip.LEDStrip) bool {
	enemy := s.GetXY(e.playerX, e.playerY)
	return compareColors(enemy, ENEMY_COLOR)
}

func compareColors(a, b strip.LED) bool {
	return len(a) == 3 && len(b) == 3 && a[0] == b[0] && a[1] == b[1] && a[2] == b[2]
}
