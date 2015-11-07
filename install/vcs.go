package install

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type vcsCmd struct {
	cmd         string
	createCmd   []string
	downloadCmd []string
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

func Download(dir, repo, version string) error {
	var err error
	var out string
	fmt.Printf("\n\tInstalling %s...", repo)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		out, err = GitVcs.download(dir, repo, version)
	} else {
		out, err = GitVcs.update(dir, repo, version)
	}

	if err != nil {
		fmt.Printf(" [\033[0;31mERR\033[0m]\n")
		fmt.Printf("\tFailed to download dependency: \n\t%s\n", out)
		return err
	}

	fmt.Printf(" [\033[0;32mOK\033[0m]\n")
	return nil
}

func (v *vcsCmd) update(dir, repo, version string) (string, error) {
	if out, err := v.fixDetachedHead(dir); err != nil {
		return out, err
	}
	for _, cmd := range v.downloadCmd {
		out, err := v.run(dir, cmd, "dir", dir, "repo", repo, "version", version)
		if err != nil {
			return out, err
		}
	}

	return "", nil
}

func (v *vcsCmd) download(dir, repo, version string) (string, error) {
	for _, cmd := range v.createCmd {
		out, err := v.run(dir, cmd, "dir", dir, "repo", repo, "version", version)
		if err != nil {
			fmt.Printf(" [\033[0;31mERR\033[0m]\n")
			fmt.Printf("\tFailed to download dependency: \n\t%s\n", out)
			return out, err
		}
	}

	return "", nil
}

func (v *vcsCmd) run(dir, cmdline string, keyval ...string) (string, error) {
	m := make(map[string]string)
	for i := 0; i < len(keyval); i += 2 {
		m[keyval[i]] = keyval[i+1]
	}
	args := strings.Fields(cmdline)
	for i, arg := range args {
		args[i] = expand(m, arg)
	}

	cmd := exec.Command(v.cmd, args...)
	cmd.Dir = dir
	var buf bytes.Buffer
	cmd.Stderr = &buf
	cmd.Stdout = &buf
	err := cmd.Run()
	out := buf.String()
	return out, err
}

func (v *vcsCmd) fixDetachedHead(dir string) (string, error) {
	return v.run(dir, "checkout master")
}
