package animator

import (
	"github.com/brianewing/redshift/effects"
	"github.com/brianewing/redshift/strip"
	"testing"
)

type testEffect struct {
	didInit, didDestroy bool
}

func (e *testEffect) Render(s *strip.LEDStrip) {}
func (e *testEffect) Init()                    { e.didInit = true }
func (e *testEffect) Destroy()                 { e.didDestroy = true }

func makeTestAnimation() (*Animator, *testEffect) {
	e := &testEffect{}
	return &Animator{Effects: []effects.Effect{e}}, e
}

func TestEffectInit(t *testing.T) {
	a, e := makeTestAnimation()
	a.Render()

	if !e.didInit {
		t.Fatal("effect did not init")
	}
}

func TestEffectDestroy(t *testing.T) {
	a, e := makeTestAnimation()
	e2 := &testEffect{}

	a.Render()
	a.SetEffects([]effects.Effect{e2})

	if a.Render(); !e.didDestroy {
		t.Fatal("effect did not destroy")
	}

	if e2.didDestroy {
		t.Fatal("new effect did destroy")
	}
}
