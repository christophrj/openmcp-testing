package internal

import (
	"io"
	"os"
	"strings"
	"text/template"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/e2e-framework/klient/k8s"
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

func ToUnstructured(obj k8s.Object) (*unstructured.Unstructured, error) {
	u, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return nil, err
	}
	return &unstructured.Unstructured{
		Object: u,
	}, nil
}

func IgnoreNotFound(err error) error {
	if errors.IsNotFound(err) {
		return nil
	}
	return err
}
