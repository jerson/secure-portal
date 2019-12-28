package util

import (
	"github.com/markbates/pkger"
	"html/template"
	"io/ioutil"
)

func LoadTemplate(path string) (*template.Template, error) {

	templateLayout, err := pkger.Open(path)
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
