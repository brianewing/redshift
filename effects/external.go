package effects

import (
	"os/exec"
	"redshift/strip"
	"log"
	"bufio"
	"io"
)

const PIPE_SIZE = 65536

// todo: write tests / benchmarks
// todo: use inotify to reload scripts when changed
// todo: .. then just save the files through the web server
// todo: define script server as a plain HTTP server, PUT/GET resources

type External struct {
	Program string
	Args []string
	Shellhack bool

	cmd *exec.Cmd
	halted bool
	stdin io.Writer
	stdout io.Reader
}

func (e *External) Render(s *strip.LEDStrip) {
	if e.cmd == nil {
		e.startProcess()
	}
	if e.halted {
		return
	}
	// write the buffer to the program's stdin, i.e. request a frame
	e.sendFrame(s.Buffer)
	// wait until the program replies, then copy its response into the strip buffer
	e.readFrame(s.Buffer)
}

func (e *External) startProcess() {
	cmd := exec.Command(e.Program, e.Args...)
	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	//log.Println(e.logPrefix(), "Starting process")
	if err := cmd.Start(); err != nil {
		log.Println(e.logPrefix(), "exec error", err)
		e.halted = true
	} else {
		e.cmd = cmd
		e.stdin = stdin
		e.stdout = stdout
		go e.readAndLogStderr(stderr)
	}
}

func (e *External) sendFrame(buffer [][]uint8) {
	frame := strip.SerializeBufferBytes(buffer)
	if e.Shellhack {
		for i, byte := range frame {
			if byte == 0 {
				frame[i] = 1
			}
		}
	}
	e.stdin.Write(frame)
}

func (e *External) readFrame(buffer [][]uint8) {
	bytes := make([]byte, PIPE_SIZE)
	if n, err := e.stdout.Read(bytes); err != nil {
		log.Println(e.logPrefix(), "stdout read error", err)
		e.halted = true
	} else {
		//log.Println(e.logPrefix(), "got", n, "bytes", bytes[:n])
		strip.UnserializeBufferBytes(buffer, bytes[:n])
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