package tools

import (
	"bufio"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/golang-mods/set"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/magefile/mage/target"
)

func Sync(list, bin string) error {
	if updated, err := target.Dir(bin, list); err != nil {
		return err
	} else if !updated {
		return nil
	}

	return SyncForce(list, bin)
}

func SyncForce(list, bin string) error {
	packages, err := readLines(list)
	if err != nil {
		return err
	}

	packageBinaries := make([]string, len(packages))
	for i, pkg := range packages {
		packageBinaries[i] = getBinaryName(pkg)
	}
	sort.Strings(packageBinaries)

	if binaries, err := os.ReadDir(bin); err == nil {
		removeBinaries := set.SortedDifference(binaries, packageBinaries, func(lhs fs.DirEntry, rhs string) int {
			name := getName(lhs.Name())
			if name < rhs {
				return -1
			}
			if name > rhs {
				return 1
			}
			return 0
		})

		for _, entry := range removeBinaries {
			if err := sh.Rm(filepath.Join(bin, entry.Name())); err != nil {
				return err
			}
		}
	}

	env := map[string]string{"GOBIN": bin}
	for _, pkg := range packages {
		if err := sh.RunWith(env, mg.GoCmd(), "install", pkg); err != nil {
			return err
		}
	}

	return nil
}

func readLines(name string) ([]string, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func getBinaryName(name string) string {
	return strings.Split(path.Base(name), "@")[0]
}

func getName(name string) string {
	return name[:len(name)-len(filepath.Ext(name))]
}
