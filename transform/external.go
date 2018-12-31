package transform

type External struct {
	Program   string

	cmd         *exec.Cmd
	halted      bool
	stdin       io.Writer
	stdout      io.Reader
}

func (t *External) Eval(code string) (interface{}, error) {
	if t.halted {
		return nil, errors.New("halted")
	}

	t.send(code)
	result, err := t.read()
	return result, err
}

func (t *External) Set(name string, val interface{}) error {
	bytes, _ := json.Marshal(val)
	setCode := fmt.Sprintf("%s = %s; null", name, bytes)
	t.send(setCode)
	t.read()
	return nil
}

func (t *External) Destroy() {
	if t.cmd != nil {
		log.Println(t.logPrefix(), "Killing process")
		t.cmd.Process.Kill()
	}
}

func (t *External) Init() {
	cmd := exec.Command(t.Program, "")
	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		log.Println(t.logPrefix(), "exec error", err)
		t.halted = true
	} else {
		t.halted = false
		t.cmd = cmd
		t.stdin = stdin
		t.stdout = stdout
		go t.readAndLogStderr(stderr)
	}
}

func (t *External) send(code string) {
	t.stdin.Write([]byte(code))
}

func (t *External) read() (interface{}, error) {
	bytes := make([]byte, PIPE_SIZE)
	if n, err := t.stdout.Read(bytes); err != nil {
		log.Println(t.logPrefix(), "stdout read error", err)
		t.halted = true
		return nil, errors.New("read error")
	} else {
		var result interface{}
		err := json.Unmarshal(bytes[:n], &result)
		return result, err
	}
}

func (t *External) readAndLogStderr(pipe io.Reader) {
	reader := bufio.NewReader(pipe)
	for {
		if msg, err := reader.ReadString('\n'); err == nil {
			log.Println(t.logPrefix(), ">>", msg[:len(msg)-1], "<<")
		} else {
			log.Println(t.logPrefix(), "stderr read error", err)
			break
		}
	}
}

func (t *External) logPrefix() string {
	return "transform.External{" + t.Program + "}"
}
