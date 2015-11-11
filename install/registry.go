package install

import (
	"net/url"

	"github.com/ChrisMcKenzie/caddy/pkg"
)

type RegistryRetreiver struct {
	registry *url.URL
	username string
	password string
}

func NewRegistryRetreiver(rawUrl, user, pass string) (*RegistryRetreiver, error) {
	registry, err := url.Parse(rawUrl)
	if err != nil {
		return nil, err
	}

	ret := &RegistryRetreiver{registry, user, pass}

	return ret, nil
}

func (r RegistryRetreiver) Download(dir string, dep *pkg.Dependency) error {
	return nil
}
