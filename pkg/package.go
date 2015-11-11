package pkg

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Package struct {
	Name            string            `json:"name"`
	Version         string            `json:"version,omitempty"`
	Scripts         map[string]string `json:"scripts,omitempty"`
	RawDependencies map[string]string `json:"dependencies"`
	Dependencies    []Dependency      `json:"-"`
}

func ReadCaddyJSON() (*Package, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	file := filepath.Join(cwd, "caddy.json")
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

func WriteCaddyJSON(p *Package) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	file := filepath.Join(cwd, "caddy.json")

	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(file, data, os.ModePerm)
}
