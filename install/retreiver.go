package install

import "github.com/ChrisMcKenzie/caddy/pkg"

type Retreiver interface {
	Download(dir string, dep *pkg.Dependency) error
}
