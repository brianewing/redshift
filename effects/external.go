package effects

import (
	"bufio"
	"io"
	"log"
	"os/exec"
	"sync"
	"time"

	"github.com/brianewing/redshift/strip"
	"github.com/fsnotify/fsnotify"
)

const PIPE_SIZE = 65536

// External is an Effect which spawns a process and allows it to be part of an LED effects chain by reading and writing rgb byte values over stdio.
// The protocol is very simple: just rgbrgbrgbrgbrgbrgbrgbrgbrgbrgb byte values exchanged via stdin/stdout, error messages over stderr.
type External struct {
	Program string
	Args    []string

	WatchForChanges bool
	changeDebouncer *time.Timer
	changeWatcher   *fsnotify.Watcher

	IsShellScript bool

	cmd    *exec.Cmd
	halted bool

	stdin  io.Writer
	stdout io.Reader

	reloadMutex sync.Mutex
}

// NewExternal makes a new External Effect.
func NewExternal() *External {
	changeDebouncer := time.NewTimer(0)
	<-changeDebouncer.C

	return &External{
		WatchForChanges: true,
		changeDebouncer: time.NewTimer(99999 * time.Hour),
	}
}

func (e *External) Init() {
	e.startProcess()

	if e.WatchForChanges {
		go e.watchForChanges()
	}
}

func (e *External) Render(s *strip.LEDStrip) {
	e.reloadMutex.Lock()

	// if child process has exited with an error, just passthrough / skip this effect
	if e.halted {
		e.reloadMutex.Unlock()
		return
	}

	e.sendFrame(s.Buffer) // write the buffer to the program's stdin, i.e. request a frame
	e.readFrame(s.Buffer) // wait until the program replies, then copy its response into the strip buffer

	e.reloadMutex.Unlock()
}

func (e *External) Destroy() {
	if e.cmd != nil {
		e.cmd.Process.Kill()
		e.cmd.Process.Release()
	}
	if e.changeWatcher != nil {
		e.changeWatcher.Close()
	}
}

func (e *External) startProcess() {
	cmd := exec.Command(e.Program, e.Args...)

	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		log.Println(e.logPrefix(), "exec error", err)
		e.halted = true
	} else {
		log.Println(e.logPrefix(), "process running")
		e.cmd, e.stdin, e.stdout = cmd, stdin, stdout
		e.halted = false
		go e.readAndLogStderr(stderr)
	}
}

func (e *External) watchForChanges() {
	if watcher, err := fsnotify.NewWatcher(); err != nil {
		log.Println(e.logPrefix(), "error watching for changes", err)
	} else {
		watcher.Add(e.Program)
		e.changeWatcher = watcher
		for {
			select {
			case <-e.changeDebouncer.C:
				e.reload()
			case _, ok := <-watcher.Events:
				if !ok { // channel has closed, let's exit this loop
					return
				}
				e.changeDebouncer.Reset(50 * time.Millisecond)
			}
		}
	}
}

func (e *External) reload() {
	e.reloadMutex.Lock()
	defer e.reloadMutex.Unlock()

	log.Println(e.logPrefix(), "reload()")

	if e.cmd != nil {
		e.cmd.Process.Kill()
		e.cmd.Process.Release()
	}

	e.cmd = nil
	e.halted = false

	e.startProcess()
}

func (e *External) sendFrame(buffer strip.Buffer) {
	frame := buffer.MarshalBytes()
	if e.IsShellScript {
		// Many shells can't process null bytes so we replace them with \001
		for i, byte := range frame {
			if byte == 0 {
				frame[i] = 1
			}
		}
	}
	e.stdin.Write(frame)
}

func (e *External) readFrame(buffer strip.Buffer) {
	bytes := make([]byte, PIPE_SIZE)
	if n, err := e.stdout.Read(bytes); err != nil {
		log.Println(e.logPrefix(), "stdout read error", err)
		e.halted = true
	} else {
		// log.Println(e.logPrefix(), "got", n, "bytes", bytes[:n])
		buffer.UnmarshalBytes(bytes[:n])
	}
}

func (e *External) readAndLogStderr(pipe io.Reader) {
	reader := bufio.NewReader(pipe)
	for {
		if msg, err := reader.ReadString('\n'); err == nil {
			log.Println(e.logPrefix(), ">>", msg[:len(msg)-1], "<<")
		} else {
			log.Println(e.logPrefix(), "stderr read error", err)
			break
		}
	}
}

func (e *External) logPrefix() string {
	return "effects.External{" + e.Program + "}"
}
