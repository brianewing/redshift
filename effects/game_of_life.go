package effects

import (
	"bytes"
	"math/rand"

	"github.com/brianewing/redshift/strip"
)

type GameOfLife struct {
	Width  int
	Height int
	Color  strip.LED

	StartingCells int

	NumPrevStates int
	Restart       bool

	life *life

	randomPrevState []bool // used to detect loops, eventually
	prevStates      [][]bool

	i int
	N int
}

func NewGameOfLife() *GameOfLife {
	return &GameOfLife{
		Color: strip.LED{255, 255, 255},
		N:     5,
		Restart: true,
	}
}

func (e *GameOfLife) Init() {
	if e.StartingCells == 0 {
		e.StartingCells = e.Width * e.Height / 4
	}
	e.life = newLife(e.Width, e.Height, e.StartingCells)
}

func (e *GameOfLife) Render(s *strip.LEDStrip) {
	if e.N == 0 {
		e.N = 1
	}

	if e.Width == 0 && e.Height == 0 {
		e.Width = s.Width
		e.Height = s.Height
	}

	state := e.life.State()

	for i, alive := range state {
		if alive {
			s.SetPixel(i, e.Color)
		}
	}

	if e.i++; e.i%e.N == 0 {
		e.life.Step()

		if e.Restart {
			for _, prevState := range e.prevStates {
				if compareStates(state, prevState) {
					e.Init()
					return
				}
			}

			if len(e.prevStates) > 500 {
				e.prevStates = e.prevStates[len(e.prevStates)-500:]
			}

			e.prevStates = append(e.prevStates, state)

			// if compareStates(s, e.randomPrevState) {
			// 	e.Init() // start over
			// 	return
			// }

			// if rand.Intn(1200) == 0 {
			// 	e.randomPrevState = s
			// }
		}
	}
}

// compareStates checks if s1 and s2 are equal (i.e. have same length and values)
func compareStates(s1, s2 []bool) bool {
	// If one is nil, the other must also be nil.
	if (s1 == nil) != (s2 == nil) {
		return false
	}
	if len(s1) != len(s2) {
		return false
	}
	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}

/*
 * Game of Life implementation
 * Copied and modified from https://golang.org/doc/play/life.go
 */

// Field represents a two-dimensional field of cells.
type lifeField struct {
	s    [][]bool
	w, h int
}

// newLifeField returns an empty field of the specified width and height.
func newLifeField(w, h int) *lifeField {
	s := make([][]bool, h)
	for i := range s {
		s[i] = make([]bool, w)
	}
	return &lifeField{s: s, w: w, h: h}
}

// Set sets the state of the specified cell to the given value.
func (f *lifeField) Set(x, y int, b bool) {
	f.s[y][x] = b
}

// Alive reports whether the specified cell is alive.
// If the x or y coordinates are outside the field boundaries they are wrapped
// toroidally. For instance, an x value of -1 is treated as width-1.
func (f *lifeField) Alive(x, y int) bool {
	x += f.w
	x %= f.w
	y += f.h
	y %= f.h
	return f.s[y][x]
}

// Next returns the state of the specified cell at the next time step.
func (f *lifeField) Next(x, y int) bool {
	// Count the adjacent cells that are alive.
	alive := 0
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if (j != 0 || i != 0) && f.Alive(x+i, y+j) {
				alive++
			}
		}
	}
	// Return next state according to the game rules:
	//   exactly 3 neighbors: on,
	//   exactly 2 neighbors: maintain current state,
	//   otherwise: off.
	return alive == 3 || alive == 2 && f.Alive(x, y)
}

// Life stores the state of a round of Conway's Game of Life.
type life struct {
	a, b *lifeField
	w, h int
}

// newLife returns a new Life game state with a random initial state.
func newLife(w, h, startingCells int) *life {
	a := newLifeField(w, h)
	if startingCells > (w * h) {
		startingCells = (w * h / 4)
	}
	for i := 0; i < startingCells; i++ {
		a.Set(rand.Intn(w), rand.Intn(h), true)
	}
	return &life{
		a: a, b: newLifeField(w, h),
		w: w, h: h,
	}
}

// Step advances the game by one instant, recomputing and updating all cells.
func (l *life) Step() {
	// Update the state of the next field (b) from the current field (a).
	for y := 0; y < l.h; y++ {
		for x := 0; x < l.w; x++ {
			l.b.Set(x, y, l.a.Next(x, y))
		}
	}
	// Swap fields a and b.
	l.a, l.b = l.b, l.a
}

func (l *life) State() []bool {
	var states []bool
	for y := 0; y < l.h; y++ {
		for x := 0; x < l.w; x++ {
			states = append(states, l.a.Alive(x, y))
		}
	}
	return states
}

func (l *life) Hash() int {
	h := 0
	for i, state := range l.State() {
		if state {
			x := i / l.w
			y := i % l.w

			h += x*50 + y
		}
	}
	return h
}

// String returns the game board as a string.
func (l *life) String() string {
	var buf bytes.Buffer
	for y := 0; y < l.h; y++ {
		for x := 0; x < l.w; x++ {
			b := byte(' ')
			if l.a.Alive(x, y) {
				b = '*'
			}
			buf.WriteByte(b)
		}
		buf.WriteByte('\n')
	}
	return buf.String()
}
