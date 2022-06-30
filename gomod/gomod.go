package gomod

import (
	"os"
	"strings"

	"github.com/golang-mods/familiar/shell"
	"github.com/samber/lo"
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

func PathToID(path string) string {
	parts := strings.Split(path, "/")

	return strings.Join(append(lo.Reverse(strings.Split(parts[0], ".")), parts[1:]...), ".")
}
