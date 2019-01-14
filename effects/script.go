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

	e.External = NewExternal()
	e.External.Program = filepath.Join(scriptsDir, e.Name)
	e.External.Args = e.Args
	e.WatchForChanges = true

	e.External.Init()
}
