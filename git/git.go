package git

import (
	"errors"

	"github.com/golang-mods/familiar/shell"
)

func Revision(current string) (string, error) {
	return validateEmpty(
		shell.Output("git", "-C", current, "rev-parse", "HEAD")(shell.StderrSilent()),
	)
}

func RevisionShort(current string) (string, error) {
	return validateEmpty(
		shell.Output("git", "-C", current, "rev-parse", "--short", "HEAD")(shell.StderrSilent()),
	)
}

func Tag(current string) (string, error) {
	return validateEmpty(
		shell.Output("git", "-C", current, "tag", "--points-at", "HEAD")(shell.StderrSilent()),
	)
}

func Clean(current string) error {
	return shell.Run("git", "-C", current, "clean", "-fdX")()
}

var ErrEmpty = errors.New("value is empty")

func validateEmpty(value string, err error) (string, error) {
	if err != nil {
		return value, err
	}
	if value == "" {
		return "", ErrEmpty
	}
	return value, err
}
