package controllers

import (
	"bytes"
	"html/template"
)

type TemplateData struct {
	BaseImage string
	// GitRepo   string
	// GitBranch string
	// BuildCmd  string
}

func GenerateScript(templatePaths []string, data TemplateData) (string, error) {
	tmpl, err := template.ParseFiles(templatePaths...)
	tmpl = template.Must(tmpl, err)
	if err != nil {
		return "test", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "test-2", err
	}

	return buf.String(), nil
}
