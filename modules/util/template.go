package util

import (
	"fmt"
	"github.com/markbates/pkger"
	"html/template"
	"io/ioutil"
)

// LoadTemplate ...
func LoadTemplate(path string) (*template.Template, error) {

	templateLayout, err := pkger.Open(fmt.Sprintf("/%s", path))
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(templateLayout)
	if err != nil {
		return nil, err
	}

	t, err := template.New("").Parse(string(data))
	if err != nil {
		return nil, err
	}
	return t, nil
}
