package transform

import (
	// JS interpreter written in Go
	"github.com/robertkrimen/otto"
)

type Otto struct {
	vm *otto.Otto
}

func (t *Otto) init() {
	if t.vm == nil {
		t.vm = otto.New()
	}
}

func (t *Otto) Eval(code string) (result interface{}, err error) {
	t.init()

	result, err := c.vm.Run(c.Transform)
	newVal, _ := result.Export()

	return newVal, err
}

func (t *Otto) Set(name string, value interface{}) error {
	t.init()
	return t.vm.Set(name, value)
}
