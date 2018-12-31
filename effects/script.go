package effects

import (
	"flag"
	"path/filepath"
)

type Script struct {
	Name string
	Args []string

	*External `json:"-"`
}

func (e *Script) Init() {
	scriptsDir := flag.Lookup("scriptsDir").Value.String()

	e.External = &External{
		Program:         filepath.Join(scriptsDir, e.Name),
		Args:            e.Args,
		WatchForChanges: true,
	}

	e.External.Init()
}
