package tools

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/golang-mods/sorted"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/magefile/mage/target"
	"github.com/samber/lo"
)

// Install the packages enumlated in the list file.
// If there is already a binary in bin, it will not be installed.
// Remove binaries packages not enumlated in the list file.
func Sync(list, bin string) error {
	if updated, err := target.Dir(bin, list); err != nil {
		return err
	} else if !updated {
		return nil
	}

	return sync(list, bin, false)
}

// Install the packages enumlated in the list file.
// Even if bin already has binaries, install them again.
// Remove binaries packages not enumlated in the list file.
func Install(list, bin string) error {
	return sync(list, bin, true)
}

func sync(list, bin string, force bool) error {
	packages, err := packagesFromList(list)
	if err != nil {
		return err
	}

	binaries, err := binariesFromDirectory(bin)
	if err != nil {
		return err
	}

	installPackages := packages
	if !force {
		installPackages = sorted.Difference(packages, binaries, func(pkg, bin entry) int {
			return sorted.Compare(pkg.name, bin.name)
		})
	}

	removeBinaries := sorted.Difference(binaries, packages, func(bin, pkg entry) int {
		return sorted.Compare(bin.name, pkg.name)
	})

	if err := errors.Join(lo.Map(removeBinaries, func(entry entry, _ int) error {
		if err := sh.Rm(filepath.Join(bin, entry.raw)); err != nil {
			return err
		}

		if mg.Verbose() {
			fmt.Printf("Remove: %s\n", entry.raw)
		}
		return nil
	})...); err != nil {
		return err
	}

	env := map[string]string{"GOBIN": bin}
	return errors.Join(lo.Map(installPackages, func(entry entry, _ int) error {
		if err := sh.RunWith(env, mg.GoCmd(), "install", entry.raw); err != nil {
			return err
		}

		if mg.Verbose() {
			fmt.Printf("Install: %s\n", entry.raw)
		}
		return nil
	})...)
}

type entry struct {
	raw  string
	name string
}

func packagesFromList(name string) ([]entry, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	entries := []entry{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		raw := scanner.Text()
		base := path.Base(raw)
		if i := strings.LastIndex(base, "@"); i != -1 {
			entries = append(entries, entry{raw: raw, name: base[:i]})
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	sort.SliceStable(entries, func(i, j int) bool { return entries[i].name < entries[j].name })

	return entries, nil
}

func binariesFromDirectory(name string) ([]entry, error) {
	directoryEntries, err := os.ReadDir(name)
	if os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	entries := []entry{}
	for _, directoryEntry := range directoryEntries {
		raw := directoryEntry.Name()
		ext := filepath.Ext(raw)
		name := raw[:len(raw)-len(ext)]
		if !directoryEntry.IsDir() && (strings.ToLower(ext) == ".exe" || ext == "") {
			entries = append(entries, entry{raw: raw, name: name})
		}
	}

	sort.SliceStable(entries, func(i, j int) bool { return entries[i].name < entries[j].name })

	return entries, nil
}
