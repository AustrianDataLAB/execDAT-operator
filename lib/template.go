package lib

import (
	"bytes"
	"html/template"

	taskv1alpha1 "github.com/AustrianDataLAB/execDAT-operator/api/v1alpha1"
)

type InitTemplateData struct {
	BaseImage string
}

func CreateTemplate[D InitTemplateData | taskv1alpha1.BuildSpec](templatePaths []string, data D) (string, error) {
	tmpl, err := template.ParseFiles(templatePaths...)
	tmpl = template.Must(tmpl, err)
	if err != nil {
		return "ERROR", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "ERROR", err
	}

	return buf.String(), nil
}
