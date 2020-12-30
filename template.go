package fs

import (
	"bytes"
	"fmt"
	"regexp"
	"text/template"
	"text/template/parse"
)

var fieldRegexp = regexp.MustCompile(`\{\{\s*\.([^\s]*)\s*\}\}`)

func fieldString(field string) string {
	return fieldRegexp.FindStringSubmatch(field)[1]
}

func listTemplateFields(t *template.Template) []string {
	fields := make([]string, 0)
	for _, node := range t.Root.Nodes {
		if node.Type() == parse.NodeAction {
			fields = append(fields, fieldString(node.String()))
		}
	}
	return fields
}

func uniqueStringSlice(s []string) (u []string) {
	m := map[string]bool{}
	for _, v := range s {
		if !m[v] {
			m[v] = true
			u = append(u, v)
		}
	}
	return u
}

func (s *Snippet) ListTemplateFields() []string {
	fields := make([]string, 0)
	fields = append(fields, listTemplateFields(template.Must(template.New("snippet").Parse(s.Name)))...)
	fields = append(fields, listTemplateFields(template.Must(template.New("snippet").Parse(s.Content)))...)
	return uniqueStringSlice(fields)
}

func executeTemplate(t *template.Template, str string, data map[string]string) (string, error) {
	t, err := t.Parse(str)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (s *Snippet) ExecuteTemplate(data map[string]string) error {
	t := template.New("snippet")

	name, err := executeTemplate(t, s.Name, data)
	if err != nil {
		return err
	}
	s.Name = name

	content, err := executeTemplate(t, s.Content, data)
	if err != nil {
		return err
	}
	s.Content = content
	return nil
}
