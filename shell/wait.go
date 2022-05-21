package shell

import (
	"bufio"
	"io"
	"os"
	"syscall"
	"time"

	"github.com/magefile/mage/mg"
)

func WaitIf(conditions func(line string) bool, name string, arguments ...string) func(...Option) (func() error, error) {
	return func(options ...Option) (func() error, error) {
		var stdout io.Writer
		if mg.Verbose() {
			stdout = os.Stdout
		}

		kill, _, err := ExecWaitIf(conditions, stdout, name, arguments...)(options...)
		if err != nil {
			return nil, err
		}

		return func() error { _, err := kill(); return err }, nil
	}
}

func ExecWaitIf(
	conditions func(line string) bool,
	stdout io.Writer, name string, arguments ...string,
) func(...Option) (func() (bool, error), bool, error) {
	return func(options ...Option) (func() (bool, error), bool, error) {
		opts := buildOptions(options)
		if stdout == nil && opts.verbose {
			stdout = os.Stdout
		}

		reader, writer := io.Pipe()
		if stdout == nil {
			stdout = writer
		} else {
			stdout = io.MultiWriter(stdout, writer)
		}

		kill, ok, err := start(stdout, opts, name, arguments...)
		if err != nil {
			return nil, ok, err
		}

		scanner := bufio.NewScanner(reader)
		for scanner.Scan() && !conditions(scanner.Text()) {
		}
		if err := scanner.Err(); err != nil {
			return nil, false, err
		}

		go func() {
			defer reader.Close()
			defer writer.Close()

			// drop data
			for scanner.Scan() {
			}
		}()

		return kill, ok, err
	}
}

func start(stdout io.Writer, options *commandOptions, name string, arguments ...string) (func() (bool, error), bool, error) {
	const waitTime = 5 * time.Second

	command, name, arguments := buildCommand(stdout, name, arguments, options)
	ok, err := executeError(command.Start(), name, arguments)
	if err != nil {
		return nil, ok, err
	}

	errorChannel := make(chan error)
	go func() {
		defer close(errorChannel)
		errorChannel <- command.Wait()
	}()
	return func() (bool, error) {
		select {
		case err := <-errorChannel:
			return executeError(err, name, arguments)
		default:
			if err := killChildren(command, syscall.SIGINT); err != nil {
				return executeError(err, name, arguments)
			}
		}

		select {
		case err := <-errorChannel:
			return executeError(err, name, arguments)
		case <-time.NewTimer(waitTime).C:
			if err := killChildren(command, syscall.SIGKILL); err != nil {
				return executeError(err, name, arguments)
			}
			return executeError(<-errorChannel, name, arguments)
		}
	}, ok, err
}
