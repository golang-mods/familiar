package command

import (
	"os"
	"strings"

	"github.com/golang-mods/familiar/ldflags"
	"github.com/magefile/mage/mage"
	"github.com/magefile/mage/mg"
)

func Compile(output string, flags *ldflags.LDFlags) error {
	if flags == nil {
		flags = &ldflags.LDFlags{}
	}
	ldflags, err := flags.Flags()
	if err != nil {
		return err
	}

	arguments := []string{"-ldflags", ldflags, "-compile", output}
	code := mage.ParseAndRun(os.Stdout, os.Stdout, os.Stdin, arguments)
	if code != 0 {
		return mg.Fatalf(code, "Execute: mage %s", strings.Join(arguments, " "))
	}
	return nil
}
