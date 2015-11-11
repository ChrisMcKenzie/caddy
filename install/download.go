package install

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/ChrisMcKenzie/caddy/pkg"
)

func Download(global bool, dep *pkg.Dependency) error {
	var err error
	retreiver, err := getRetreiverByType(dep)
	if err != nil {
		return err
	}

	var path string
	if global {
		path = filepath.Join(os.Getenv("GOPATH"), "/src")
	} else {
		path = "vendor"
	}

	err = retreiver.Download(path, dep)
	return err
}

func getRetreiverByType(dep *pkg.Dependency) (Retreiver, error) {
	switch dep.Type {
	case "hosted":
		return getRetreiverByUrlScheme(dep)
	case "range", "version":
		return NewRegistryRetreiver("http://localhost:8000", "", "")
	}
	return nil, errors.New("unable to determine dependency source")
}

func getRetreiverByUrlScheme(dep *pkg.Dependency) (Retreiver, error) {
	switch dep.Hosted.Scheme {
	case "https", "ssh":
		return NewVcsRetreiver(dep.Hosted)
	}
	return nil, errors.New("unable to determine dependency source")
}
