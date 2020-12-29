package cli

import (
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/skmatz/fs"
)

const (
	headSize = 10
	itemSize = 10
)

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func head(content string) string {
	lines := strings.Split(content, "\n")
	return strings.Join(lines[:min(len(lines), headSize)], "\n") + "\n"
}

func (c *CLI) SelectSnippet(snippets []fs.Snippet) (*fs.Snippet, error) {
	funcMap := promptui.FuncMap
	funcMap["head"] = head

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   promptui.IconSelect + " {{ .ID | cyan }}",
		Inactive: "  {{ .ID | faint }}",
		Selected: promptui.IconGood + " {{ .Name }}",
		Details: `
{{ "ID:" | faint }}	{{ .ID }}
{{ "Name:" | faint }}	{{ .Name }}
{{ "Content:" | faint }}
{{ .Content | head }}
`,
		FuncMap: funcMap,
	}

	searcher := func(input string, index int) bool {
		id := strings.ToLower(snippets[index].ID)
		name := strings.ToLower(snippets[index].Name)
		return strings.Contains(id, input) || strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:             "Select snippet",
		Items:             snippets,
		Size:              itemSize,
		Templates:         templates,
		Searcher:          searcher,
		StartInSearchMode: true,
	}

	idx, _, err := prompt.Run()
	if err != nil {
		return nil, err
	}
	return &snippets[idx], nil
}
