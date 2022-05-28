package shell

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

func Run(name string, arguments ...string) func(...Option) error {
	return func(options ...Option) error {
		var stdout io.Writer
		if mg.Verbose() {
			stdout = os.Stdout
		}
		_, err := Exec(stdout, name, arguments...)(options...)
		return err
	}
}

func Output(name string, arguments ...string) func(...Option) (string, error) {
	return func(options ...Option) (string, error) {
		buffer := &bytes.Buffer{}
		_, err := Exec(buffer, name, arguments...)(options...)
		return strings.TrimSuffix(buffer.String(), "\n"), err
	}
}

func Exec(stdout io.Writer, name string, arguments ...string) func(...Option) (bool, error) {
	return func(options ...Option) (bool, error) {
		opts := buildOptions(options)
		if stdout == nil && opts.verbose {
			stdout = os.Stdout
		}
		command, name, arguments := buildCommand(stdout, name, arguments, opts)

		verboseCommand(name, arguments...)

		return executeError(command.Run(), name, arguments)
	}
}

func buildCommand(stdout io.Writer, name string, arguments []string, options *commandOptions) (*exec.Cmd, string, []string) {
	offset := len(arguments)
	args := make([]string, offset+len(options.arguments))
	copy(args, arguments)
	copy(args[offset:], options.arguments)

	name, args = expand(options.environment, name, args)
	command := exec.Command(name, args...)
	command.Env = getEnvironment(options.environment)
	command.Dir = options.directory
	command.Stdin = os.Stdin
	command.Stdout = stdout
	command.Stderr = options.stderr

	return command, name, args
}

func executeError(err error, name string, arguments []string) (bool, error) {
	if err == nil {
		return true, nil
	} else if sh.CmdRan(err) {
		code := sh.ExitStatus(err)
		return true, mg.Fatalf(code, `running "%s %s" failed with exit code %d`, name, strings.Join(arguments, " "), code)
	} else {
		return false, fmt.Errorf(`failed to run "%s %s: %v"`, name, strings.Join(arguments, " "), err)
	}
}

func expand(environment map[string]string, name string, arguments []string) (string, []string) {
	mapping := func(key string) string {
		if value, ok := environment[key]; ok {
			return value
		}
		return os.Getenv(key)
	}

	name = os.Expand(name, mapping)
	for i := range arguments {
		arguments[i] = os.Expand(arguments[i], mapping)
	}
	return name, arguments
}

func getEnvironment(environment map[string]string) []string {
	env := os.Environ()
	for key, value := range environment {
		env = append(env, key+"="+value)
	}
	return env
}

func verboseCommand(command string, arguments ...string) {
	if mg.Verbose() {
		var quoted []string
		for _, argument := range arguments {
			quoted = append(quoted, fmt.Sprintf("%q", argument))
		}
		log.Println("exec:", command, strings.Join(quoted, " "))
	}
}
