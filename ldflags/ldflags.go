package ldflags

import (
	"fmt"
	"strings"

	"github.com/golang-mods/familiar/gomod"
)

type Definition struct {
	Name  string
	Value string
}

type LDFlags struct {
	Definitions            []Definition
	MainDefinitions        []Definition
	ModuleDefinitions      []Definition
	DisableSymbolTable     bool
	DisableDWARFGeneration bool
}

const ldFlagsFieldCount = 2

func (ldflags *LDFlags) Flags() (string, error) {
	size := ldFlagsFieldCount +
		len(ldflags.Definitions) +
		len(ldflags.MainDefinitions) +
		len(ldflags.ModuleDefinitions)
	flags := make([]string, 0, size)

	for _, definition := range ldflags.Definitions {
		flags = append(flags, fmt.Sprintf("-X \"%s=%s\"", definition.Name, definition.Value))
	}

	for _, definition := range ldflags.MainDefinitions {
		flags = append(flags, fmt.Sprintf("-X \"main.%s=%s\"", definition.Name, definition.Value))
	}

	if len(ldflags.ModuleDefinitions) > 0 {
		module, err := gomod.ModulePath()
		if err != nil {
			return "", err
		}

		for _, definition := range ldflags.ModuleDefinitions {
			name := definition.Name
			if strings.ContainsRune(name, '.') {
				name = "/" + name
			} else {
				name = "." + name
			}

			flags = append(flags, fmt.Sprintf("-X \"%s%s=%s\"", module, name, definition.Value))
		}
	}

	if ldflags.DisableSymbolTable {
		flags = append(flags, "-s")
	}

	if ldflags.DisableDWARFGeneration {
		flags = append(flags, "-w")
	}

	return strings.Join(flags, " "), nil
}
