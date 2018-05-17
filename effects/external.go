package effects

import (
	"bufio"
	"github.com/brianewing/redshift/strip"
	"github.com/fsnotify/fsnotify"
	"io"
	"log"
	"os/exec"
	"sync"
)

const PIPE_SIZE = 65536

// todo: write tests / benchmarks

type External struct {
	Program   string
	Args      []string
	Shellhack bool

	cmd         *exec.Cmd
	halted      bool
	stdin       io.Writer
	stdout      io.Reader
	watcher     *fsnotify.Watcher
	reloadMutex sync.Mutex
}

func (e *External) Render(s *strip.LEDStrip) {
	e.reloadMutex.Lock()

	if e.cmd == nil && !e.halted {
		//log.Println(e.logPrefix(), "Starting process")
		e.startProcess()
	}
	if e.watcher == nil {
		log.Println(e.logPrefix(), "Watching for changes")
		go e.watchForChanges()
	}
	if e.halted {
		e.reloadMutex.Unlock()
		return
	}

	// write the buffer to the program's stdin, i.e. request a frame
	e.sendFrame(s.Buffer)
	// wait until the program replies, then copy its response into the strip buffer
	e.readFrame(s.Buffer)

	e.reloadMutex.Unlock()
}

func (e *External) Destroy() {
	if e.cmd != nil {
		log.Println(e.logPrefix(), "Killing process")
		e.cmd.Process.Kill()
	}
	if e.watcher != nil {
		log.Println(e.logPrefix(), "Stopping watcher")
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
		for range watcher.Events {
			e.reload()
		}
	}
}

func (e *External) reload() {
	e.reloadMutex.Lock()
	log.Println(e.logPrefix(), "reloading")
	if e.cmd != nil {
		e.cmd.Process.Kill()
		e.cmd = nil
	}
	e.halted = false
	e.reloadMutex.Unlock()
}

func (e *External) sendFrame(buffer strip.Buffer) {
	frame := buffer.MarshalBytes()
	if e.Shellhack {
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
		//log.Println(e.logPrefix(), "got", n, "bytes", bytes[:n])
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
