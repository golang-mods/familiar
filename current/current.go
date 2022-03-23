package current

import (
	"path/filepath"
	"runtime"
)

func File() string {
	_, file, _, _ := runtime.Caller(1)
	return file
}

func Directory() string {
	_, file, _, _ := runtime.Caller(1)
	return filepath.Dir(file)
}
