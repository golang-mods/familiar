package shell

import (
	"io"
	"os"
)

type commandOptions struct {
	environment map[string]string
	directory   string
	verbose     bool
	stderr      io.Writer
	arguments   []string
}

func buildOptions(options []Option) *commandOptions {
	commandOptions := &commandOptions{
		stderr: os.Stderr,
	}
	for _, option := range options {
		option(commandOptions)
	}
	return commandOptions
}

type Option func(*commandOptions)

func Env(environment map[string]string) Option {
	return func(options *commandOptions) {
		options.environment = environment
	}
}

func Dir(directory string) Option {
	return func(options *commandOptions) {
		options.directory = directory
	}
}

func Verbose() Option {
	return func(options *commandOptions) {
		options.verbose = true
	}
}

func Stderr(stderr io.Writer) Option {
	return func(options *commandOptions) {
		options.stderr = stderr
	}
}

func StderrSilent() Option {
	return func(options *commandOptions) {
		options.stderr = nil
	}
}

func Args(arguments ...string) Option {
	return func(options *commandOptions) {
		options.arguments = arguments
	}
}
