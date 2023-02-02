package gomod

import (
	"strings"

	"github.com/samber/lo"
	"golang.org/x/mod/modfile"
)

var internalLoader loader

func File() (*modfile.File, error) {
	path, data, err := internalLoader.load()
	if err != nil {
		return nil, err
	}

	return modfile.Parse(path, data, nil)
}

func ModulePath() (string, error) {
	_, data, err := internalLoader.load()
	if err != nil {
		return "", err
	}

	return modfile.ModulePath(data), nil
}

func PathToID(path string) string {
	parts := strings.Split(path, "/")

	return strings.Join(append(lo.Reverse(strings.Split(parts[0], ".")), parts[1:]...), ".")
}
