package install

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/ChrisMcKenzie/caddy/pkg"
	"github.com/ChrisMcKenzie/dropship/installer"
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
	resp, err := http.Get("http://localhost:5984/caddy/" + dep.Name + "/" + dep.Spec)
	if err != nil {
		return err
	}

	msg, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var pkg map[string]interface{}
	err = json.Unmarshal(msg, &pkg)

	dlUrl := pkg["dist"].(map[string]interface{})["url"].(string)
	resp, err = http.Get(dlUrl)

	var install installer.TarInstaller

	install.Install(filepath.Join(dir, dep.Name), resp.Body)

	return err
}
