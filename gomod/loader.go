package gomod

import (
	"os"
	"sync"

	"github.com/golang-mods/familiar/shell"
)

type loader struct {
	mutex     sync.Mutex
	fulfilled bool

	file  string
	data  []byte
	error error
}

func (loader *loader) load() (string, []byte, error) {
	loader.mutex.Lock()
	defer loader.mutex.Unlock()

	if !loader.fulfilled {
		loader.fulfilled = true

		loader.file, loader.data, loader.error = func() (string, []byte, error) {
			path, err := shell.Output("go", "env", "GOMOD")(shell.StderrSilent())
			if err != nil {
				return "", nil, err
			}

			data, err := os.ReadFile(path)
			if err != nil {
				return "", nil, err
			}

			return path, data, nil
		}()
	}

	return loader.file, loader.data, loader.error
}
