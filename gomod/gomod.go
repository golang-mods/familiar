package gomod

import (
	"os"

	"github.com/golang-mods/familiar/shell"
	"golang.org/x/mod/modfile"
)

func ModulePath() (string, error) {
	path, err := shell.Output("go", "env", "GOMOD")(shell.StderrSilent())
	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return modfile.ModulePath(data), nil
}
