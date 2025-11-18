package internal

import (
	"bytes"
	"io"
	"os"
	"regexp"
	"strings"
	"text/template"
)

func SubstitutePlaceholders(reader io.Reader) (io.Reader, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	// any environment variable like {{PLACE_HOLDER}}
	re := regexp.MustCompile(`\{\{([A-Z0-9_]+)\}\}`)
	result := re.ReplaceAllStringFunc(string(data), func(match string) string {
		// Extract the placeholder inside {{PLACE_HOLDER}}
		placeholder := re.FindStringSubmatch(match)[1]
		value := os.Getenv(placeholder)
		if value == "" {
			// if the placeholder can't be replaced, it remains unmodified
			return match
		}
		return value
	})
	return bytes.NewReader([]byte(result)), nil
}

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
