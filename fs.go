package fs

import (
	"os"

	"github.com/naoina/toml"
	"golang.org/x/sync/errgroup"
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

	var eg errgroup.Group
	for i := range fs.Snippets {
		i := i
		eg.Go(func() error {
			if err := fs.Snippets[i].setContent(); err != nil {
				return err
			}

			fs.Snippets[i].setColoredContent()
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return &fs, nil
}
