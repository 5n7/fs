package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/skmatz/fs"
)

type CLI struct{}

func New() *CLI {
	return &CLI{}
}

type Options struct {
	Mode string
}

func (c *CLI) defaultConfigPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "fs", "fs.toml"), nil
}

func (c *CLI) Run(opt Options) error {
	path, err := c.defaultConfigPath()
	if err != nil {
		return err
	}

	app, err := fs.New(path)
	if err != nil {
		return err
	}

	snippet, err := c.SelectSnippet(app.Snippets)
	if err != nil {
		return err
	}

	switch opt.Mode {
	case "clipboard":
		return snippet.ToClipboard()
	case "file":
		return snippet.ToFile()
	default:
		return fmt.Errorf("unknown mode: %s", opt.Mode)
	}
}
