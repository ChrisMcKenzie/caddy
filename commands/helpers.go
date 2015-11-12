package commands

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"text/template"
)

var funcMap template.FuncMap

// Check if a file or directory exists.
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func writeTemplateToFile(path string, file string, template string, data interface{}) error {
	filename := filepath.Join(path, file)

	r, err := templateToReader(template, data)

	if err != nil {
		return err
	}

	err = safeWriteToDisk(filename, r)

	if err != nil {
		return err
	}
	return nil
}

func templateToReader(tpl string, data interface{}) (io.Reader, error) {
	tmpl := template.New("")
	tmpl.Funcs(funcMap)
	tmpl, err := tmpl.Parse(tpl)

	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, data)

	return buf, err
}

// Same as WriteToDisk but checks to see if file/directory already exists.
func safeWriteToDisk(inpath string, r io.Reader) (err error) {
	dir, _ := filepath.Split(inpath)
	ospath := filepath.FromSlash(dir)

	if ospath != "" {
		err = os.MkdirAll(ospath, 0777) // rwx, rw, r
		if err != nil {
			return
		}
	}

	ex, err := exists(inpath)
	if err != nil {
		return
	}
	if ex {
		return fmt.Errorf("%v already exists", inpath)
	}

	file, err := os.Create(inpath)
	if err != nil {
		return
	}
	defer file.Close()

	_, err = io.Copy(file, r)
	return
}
