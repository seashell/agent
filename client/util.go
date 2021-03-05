package client

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/seashell/agent/pkg/uuid"
)

// Read file contents if they exist, or persist and return a default value otherwise.
func (c *Client) readFileLazy(path string, s string) (string, error) {

	var out string

	buf, err := ioutil.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return "", err
	}

	if len(buf) != 0 {
		out = string(buf)
	} else {
		out = s
		if err := ioutil.WriteFile(path, []byte(s), 0700); err != nil {
			return "", err
		}
	}

	return out, nil
}

func renderTemplateToFile(tmplStr string, out string, data interface{}) error {

	tmpl, err := template.New(strings.Split(uuid.Generate(), "-")[0]).Parse(tmplStr)
	if err != nil {
		return fmt.Errorf("error parsing template : %v", err)

	}

	f, err := os.Create(out)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}

	defer f.Close()

	err = tmpl.Execute(f, data)
	if err != nil {
		return fmt.Errorf("error rendering template: %v", err)
	}

	return nil
}
