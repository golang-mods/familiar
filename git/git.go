package git

import "github.com/golang-mods/familiar/shell"

func Revision(current string) (string, error) {
	return shell.Output("git", "-C", current, "rev-parse", "HEAD")(shell.StderrSilent())
}

func RevisionShort(current string) (string, error) {
	return shell.Output("git", "-C", current, "rev-parse", "--short", "HEAD")(shell.StderrSilent())
}

func Tag(current string) (string, error) {
	return shell.Output("git", "-C", current, "tag", "--points-at", "HEAD")(shell.StderrSilent())
}
