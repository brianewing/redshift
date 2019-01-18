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

// todo: write tests / benchmarks

type External struct {
	Program string
	Args    []string

	WatchForChanges bool
	IsShellScript   bool

	watcher *fsnotify.Watcher

	cmd    *exec.Cmd
	stdin  io.Writer
	stdout io.Reader
	halted bool

	reloadMutex sync.Mutex
	debouncer   *time.Timer
}

func NewExternal() *External {
	return &External{
		WatchForChanges: true,
		debouncer:       time.NewTimer(999999 * time.Hour),
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
	if e.watcher != nil {
		e.watcher.Close()
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
		e.halted = false
		e.cmd = cmd
		e.stdin = stdin
		e.stdout = stdout
		go e.readAndLogStderr(stderr)
	}
}

func (e *External) watchForChanges() {
	if watcher, err := fsnotify.NewWatcher(); err != nil {
		log.Println(e.logPrefix(), "error watching for changes", err)
	} else {
		e.watcher = watcher
		watcher.Add(e.Program)
		for {
			select {
			case _, ok := <-watcher.Events:
				if !ok { // channel has closed
					return
				}
				log.Println("reset debouncer")
				e.debouncer.Reset(1500 * time.Millisecond)
			case <-e.debouncer.C:
				log.Println("debouncer fired")
				e.reload()
			}
		}
	}
}

func (e *External) reload() {
	e.reloadMutex.Lock()
	log.Println(e.logPrefix(), "reloading")
	if e.cmd != nil {
		e.cmd.Process.Kill()
		e.cmd.Process.Release()
		e.cmd = nil
	}
	e.halted = false
	e.startProcess()
	log.Println("restarted process, e.halted is", e.halted)
	e.reloadMutex.Unlock()
}

func (e *External) sendFrame(buffer strip.Buffer) {
	frame := buffer.MarshalBytes()
	if e.IsShellScript {
		// shell scripts can't process null bytes so we change them to \001
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
