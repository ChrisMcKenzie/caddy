package pkg

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/wmark/semver"
)

var registryModuleRegEx = "^(?:([^/]+?)[/])?([^/]+?)$"

type Dependency struct {
	Raw     string
	Name    string
	Scope   string
	Type    string
	Spec    string
	RawSpec string
	Hosted  *url.URL
}

func parseDependencies(p map[string]string) ([]Dependency, error) {
	var deps []Dependency
	for pkg, version := range p {
		dep, err := parse(pkg, version)
		if err != nil {
			return nil, err
		}
		deps = append(deps, dep)
	}

	return deps, nil
}

func parse(depName, version string) (dep Dependency, err error) {
	matches, err := regexp.MatchString(registryModuleRegEx, depName+"@"+version)
	if err != nil {
		return
	}

	if matches {
		nameParts := strings.Split(depName, "/")
		typ := getVersionType(version)
		dep = Dependency{
			Raw:     depName + "@" + version,
			Name:    depName,
			Scope:   nameParts[0],
			Type:    typ,
			Spec:    version,
			RawSpec: version,
		}
	} else {
		var hosted *url.URL
		hosted, err = url.Parse(version)
		dep = Dependency{
			Raw:    version,
			Name:   depName,
			Type:   "hosted",
			Spec:   hosted.Fragment,
			Hosted: hosted,
		}
		if err != nil {
			return
		}
	}

	return
}

func getVersionType(version string) (typ string) {
	_, err := semver.NewVersion(version)
	if err != nil && err.Error() == "Given string does not resemble a Version." {
		semver.NewRange(version)
		return "range"
	}

	return "version"
}
