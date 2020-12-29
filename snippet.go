package fs

import (
	"io/ioutil"
	"os"

	"github.com/atotto/clipboard"
)

type Snippet struct {
	ID      string `toml:"id"`
	Name    string `toml:"name"`
	Content string `toml:"content"`
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
