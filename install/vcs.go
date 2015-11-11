package install

import (
	"bytes"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ChrisMcKenzie/caddy/pkg"
)

type vcsCmd struct {
	cmd         string
	createCmd   []string
	downloadCmd []string
}

func NewVcsRetreiver(rawUrl *url.URL) (*vcsCmd, error) {
	return &GitVcs, nil
}

var GitVcs = vcsCmd{
	cmd: "git",
	createCmd: []string{
		"clone {repo} {dir}",
		"--git-dir={dir}/.git submodule update --init --recursive",
		"--git-dir={dir}/.git checkout {version}",
	},
	downloadCmd: []string{
		"pull --ff-only",
		"submodule update --init --recursive",
		"checkout {version}",
	},
}

// expand rewrites s to replace {k} with match[k] for each key k in match.
func expand(match map[string]string, s string) string {
	for k, v := range match {
		s = strings.Replace(s, "{"+k+"}", v, -1)
	}
	return s
}

func (v *vcsCmd) Download(dir string, dep *pkg.Dependency) error {
	var err error
	fullPath := filepath.Join(dir, dep.Name)
	host := url.URL{
		Host:   dep.Hosted.Host,
		Scheme: dep.Hosted.Scheme,
		Path:   dep.Hosted.Path,
	}
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		err = os.MkdirAll(fullPath, os.ModePerm)
		if err != nil {
			return err
		}
		_, err = v.download(dep.Name, host.String(), dep.Hosted.Fragment)
	} else {
		_, err = v.update(dep.Name, host.String(), dep.Hosted.Fragment)
	}

	return err
}

func (v *vcsCmd) update(dir, repo, version string) (string, error) {
	if out, err := v.fixDetachedHead(dir); err != nil {
		return out, err
	}
	for _, cmd := range v.downloadCmd {
		out, err := v.run(dir, dir, cmd, "dir", dir, "repo", repo, "version", version)
		if err != nil {
			return out, err
		}
	}

	return "", nil
}

func (v *vcsCmd) download(dir, repo, version string) (string, error) {
	for _, cmd := range v.createCmd {
		out, err := v.run("vendor", dir, cmd, "dir", dir, "repo", repo, "version", version)
		if err != nil {
			return out, err
		}
	}

	return "", nil
}

func (v *vcsCmd) run(cwd, dir, cmdline string, keyval ...string) (string, error) {
	m := make(map[string]string)
	for i := 0; i < len(keyval); i += 2 {
		m[keyval[i]] = keyval[i+1]
	}
	args := strings.Fields(cmdline)
	for i, arg := range args {
		args[i] = expand(m, arg)
	}

	cmd := exec.Command(v.cmd, args...)
	cmd.Dir = cwd
	var buf bytes.Buffer
	cmd.Stderr = &buf
	cmd.Stdout = &buf
	err := cmd.Run()
	out := buf.String()
	return out, err
}

func (v *vcsCmd) fixDetachedHead(dir string) (string, error) {
	return v.run(dir, dir, "checkout master")
}
