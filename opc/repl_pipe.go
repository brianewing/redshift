package opc

import "io"

// replPipe acts as a bridge between a repl and an opc client.
type replPipe struct {
	inputReader *io.PipeReader
	inputWriter *io.PipeWriter

	channel uint8
	client  Writer

	io.ReadWriteCloser
}

func newReplPipe(channel uint8, client Writer) replPipe {
	inputR, inputW := io.Pipe()
	return replPipe{inputR, inputW, channel, client, nil}
}

func (p replPipe) Read(buf []byte) (int, error) {
	return p.inputReader.Read(buf)
}

func (p replPipe) Write(data []byte) (int, error) {
	msg := Message{
		Channel: p.channel,
		Command: 255,
		SystemExclusive: SystemExclusive{
			Command: CmdRepl,
			Data:    data,
		},
	}
	return len(msg.Bytes()), p.client.WriteOpc(msg)
}

func (p replPipe) Close() error {
	if err := p.inputReader.Close(); err != nil {
		return err
	}
	if err := p.inputWriter.Close(); err != nil {
		return err
	}
	return nil
}
