package cmd

import (
	"github.com/skmatz/fs/cli"
	"github.com/spf13/cobra"
)

var (
	opt     cli.Options
	version bool
)

func runRoot(cmd *cobra.Command, args []string) error {
	if version {
		return runVersion(cmd, args)
	}

	c := cli.New()
	return c.Run(opt)
}

var rootCmd = &cobra.Command{
	Use:   "fs",
	Short: "Run FS",
	Long:  "Run FS.",
	RunE:  runRoot,
}

func init() {
	rootCmd.Flags().StringVarP(&opt.Mode, "mode", "m", "file", "output mode")
	rootCmd.Flags().BoolVarP(&version, "version", "V", false, "show version")
}

func Execute() {
	rootCmd.Execute() //nolint:errcheck
}
