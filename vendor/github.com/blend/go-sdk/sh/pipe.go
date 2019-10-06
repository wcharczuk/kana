package sh

import (
	"io"
	"os"
	"os/exec"
	"sync"
)

// Pipe runs commands in sequence, piping combined output to the standard in of the next command.
func Pipe(commands ...*exec.Cmd) error {
	wg := sync.WaitGroup{}
	wg.Add(len(commands))

	errors := make(chan error, len(commands))
	readers := make([]io.Reader, len(commands))
	writers := make([]io.Writer, len(commands))
	for index := 0; index < len(commands); index++ {
		// set up pipes
		readers[index], writers[index] = io.Pipe()

		// wire up pipes
		if index == 0 { // the first command
			commands[index].Stdin = os.Stdin
			commands[index].Stdout = writers[index]
			commands[index].Stderr = writers[index]
		} else if index == len(commands)-1 { // the last command
			commands[index].Stdin = readers[index-1]
			commands[index].Stdout = os.Stdout
			commands[index].Stderr = os.Stderr
		} else { // intermediate commands
			commands[index].Stdin = readers[index-1]
			commands[index].Stdout = writers[index]
			commands[index].Stderr = writers[index]
		}
		go func(index int, cmd *exec.Cmd) {
			defer wg.Done()
			if index > 0 {
				defer func() {
					if typed, ok := cmd.Stdout.(*io.PipeWriter); ok {
						typed.Close()
					}
					if typed, ok := cmd.Stderr.(*io.PipeWriter); ok {
						typed.Close()
					}
					if typed, ok := cmd.Stdin.(*io.PipeReader); ok {
						typed.Close()
					}
				}()
			}
			if err := cmd.Run(); err != nil {
				if !IsEPIPE(err) {
					errors <- err
				}
			}
		}(index, commands[index])
	}
	wg.Wait()
	if len(errors) > 0 {
		return <-errors
	}
	return nil
}
