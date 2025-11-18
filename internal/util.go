package internal

import (
	"io"
	"os"
	"strings"
	"text/template"
)

func ExecTemplate(textTemplate string, data interface{}) (string, error) {
	tmpl, err := template.New("t").Parse(textTemplate)
	if err != nil {
		return "", err
	}
	result := strings.Builder{}
	if err := tmpl.Execute(&result, data); err != nil {
		return "", err
	}
	return result.String(), nil
}

func ExecTemplateFile(filePath string, data interface{}) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	bytes, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}
	return ExecTemplate(string(bytes), data)
}
