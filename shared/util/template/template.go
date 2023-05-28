package template

import (
	"embed"
	"strings"
)

type Template string

func FromFile(fs embed.FS, filePath string) Template {
	byteArr, _ := fs.ReadFile(filePath)
	return Template(string(byteArr))
}

func (t Template) Apply(variables map[string]string) string {
	result := string(t)

	for key, value := range variables {
		result = strings.ReplaceAll(result, "{{"+key+"}}", value)
	}

	return result
}
