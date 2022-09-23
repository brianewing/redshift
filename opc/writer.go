package opc

import "io"
//import "log"
import "sync"

type IOWriter struct {
	io.Writer
	sync.Mutex
}

func (w *IOWriter) WriteOpc(msg Message) error {
	w.Lock()
	defer w.Unlock()

	_, err := w.Write(msg.Bytes())
	return err
}

