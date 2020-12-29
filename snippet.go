package fs

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/atotto/clipboard"
)

type Snippet struct {
	ID      string `toml:"id"`
	Name    string `toml:"name"`
	Content string `toml:"content"`
	Path    string `toml:"path"`
	URL     string `toml:"url"`

	ColoredContent string
}

type mode int

const (
	modeContent mode = iota
	modePath
	modeURL
)

func (s *Snippet) getMode() (mode, error) {
	if s.Content != "" {
		return modeContent, nil
	}
	if s.Path != "" {
		return modePath, nil
	}
	if s.URL != "" {
		return modeURL, nil
	}
	return 0, fmt.Errorf("invalid snippet mode")
}

func expandTilde(path string) (string, error) {
	if !strings.HasPrefix(path, "~") {
		return path, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	if path == "~" {
		return home, err
	}
	return filepath.Join(home, path[2:]), err
}

func (s *Snippet) getContentFromPath() (string, error) {
	path, err := expandTilde(s.Path)
	if err != nil {
		return "", err
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (s *Snippet) getContentFromURL() (string, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, s.URL, nil)
	if err != nil {
		return "", err
	}

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("http response not OK: %s", resp.Status)
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (s *Snippet) setContent() error {
	mode, err := s.getMode()
	if err != nil {
		return err
	}

	var content string
	switch mode {
	case modeContent:
		content = s.Content
	case modePath:
		content, err = s.getContentFromPath()
		if err != nil {
			return err
		}
	case modeURL:
		content, err = s.getContentFromURL()
		if err != nil {
			return err
		}
	}

	s.Content = content
	return nil
}

func (s *Snippet) coloredContent(theme string) string {
	lexer := lexers.Match(s.Name)
	if lexer == nil {
		lexer = lexers.Fallback
	}

	style := styles.Get(theme)
	if style == nil {
		style = styles.Fallback
	}

	formatter := formatters.Get("terminal256")
	if formatter == nil {
		formatter = formatters.Fallback
	}

	iterator, err := lexer.Tokenise(nil, s.Content)
	if err != nil {
		return s.Content
	}

	var buf bytes.Buffer
	if err := formatter.Format(&buf, style, iterator); err != nil {
		return s.Content
	}
	return buf.String()
}

var theme = os.Getenv("FS_THEME")

func (s *Snippet) setColoredContent() {
	s.ColoredContent = s.coloredContent(theme)
}

func (s *Snippet) ToClipboard() error {
	if err := clipboard.WriteAll(s.Content); err != nil {
		return err
	}
	return nil
}

func (s *Snippet) ToFile() error {
	if err := ioutil.WriteFile(s.Name, []byte(s.Content), os.ModePerm); err != nil {
		return err
	}
	return nil
}
