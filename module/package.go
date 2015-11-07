package module

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Package struct {
	Name         string            `json:"name"`
	Dependencies map[string]string `json:"dependencies"`
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

	return &pkg, err
}
