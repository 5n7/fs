package fs

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/atotto/clipboard"
)

type Snippet struct {
	ID      string `toml:"id"`
	Name    string `toml:"name"`
	Content string `toml:"content"`
	Path    string `toml:"path"`
}

type mode int

const (
	modeContent mode = iota
	modePath
)

func (s *Snippet) getMode() (mode, error) {
	if s.Content != "" {
		return modeContent, nil
	}
	if s.Path != "" {
		return modePath, nil
	}
	return 0, fmt.Errorf("invalid snippet mode")
}

func (s *Snippet) getContentFromPath() (string, error) {
	bytes, err := ioutil.ReadFile(s.Path)
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
	}

	s.Content = content
	return nil
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
