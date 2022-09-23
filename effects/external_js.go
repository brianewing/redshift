package effects

import (
	"github.com/brianewing/redshift/strip"
)

type External struct {
	Name string
	Program string
	Args []string
	WatchForChanges bool
}

func NewExternal() *External {
	return &External{}
}

func (e *External) Init() {

}

func (e *External) Render(s *strip.LEDStrip) {

}

