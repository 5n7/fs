package fs

import (
	"os"

	"github.com/naoina/toml"
)

type FS struct {
	Snippets []Snippet `toml:"snippet"`
}

func New(path string) (*FS, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var fs FS
	if err := toml.NewDecoder(f).Decode(&fs); err != nil {
		return nil, err
	}
	return &fs, nil
}
