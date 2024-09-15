package utilities

import (
	"bytes"
	"encoding/json"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"regexp"
	"strings"
	"text/template"
)

func RenderInlineTemplate(tpl string, data any) (string, error) {
	if data == nil {
		data = make(map[string]any)
	}

	parsed, err := template.New("inline").Option().Parse(tpl)

	if err != nil {
		return "", err
	}

	placeholders := findAllPlaceholdersInTemplate(tpl)

	if len(placeholders) > 0 {
		encoded, _ := json.Marshal(data)

		for _, placeholder := range placeholders {
			cleaned := strings.TrimPrefix(strings.TrimSpace(placeholder), ".")

			if !gjson.ParseBytes(encoded).Get(cleaned).Exists() {
				encoded, _ = sjson.SetBytes(encoded, cleaned, "{{"+placeholder+"}}")
			}
		}

		_ = json.Unmarshal(encoded, &data)
	}

	buf := new(bytes.Buffer)

	if err := parsed.Execute(buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

var templatePlaceholderRegexp = regexp.MustCompile(`\{\{(\s*[a-zA-Z0-9_.]+\s*)\}\}`)

func findAllPlaceholdersInTemplate(template string) []string {
	matches := templatePlaceholderRegexp.FindAllStringSubmatch(template, -1)

	var placeholders []string

	for _, match := range matches {
		placeholders = append(placeholders, match[1])
	}

	return placeholders
}
