package pkg

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Package struct {
	Name            string            `json:"name"`
	Version         string            `json:"version"`
	Scripts         map[string]string `json:"scripts"`
	RawDependencies map[string]string `json:"dependencies"`
	Dependencies    []Dependency      `json:"-"`
}

func ReadPackageJSON() (*Package, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	file := filepath.Join(cwd, "package.json")
	p, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var pkg Package
	err = json.Unmarshal(p, &pkg)
	if err != nil {
		return nil, err
	}

	pkg.Dependencies, err = parseDependencies(pkg.RawDependencies)

	return &pkg, err
}
