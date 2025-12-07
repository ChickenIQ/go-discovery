package gen

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var flags struct {
	outputPath string
	force      bool
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&flags.outputPath, "output", "o", "-", "Output file path | '-' for stdout")
	RootCmd.PersistentFlags().BoolVarP(&flags.force, "force", "", false, "Force overwrite of output file if it exists")
}

var RootCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate Resources",
}

func output(data []byte) error {
	if flags.outputPath == "-" {
		_, err := os.Stdout.Write(data)
		return err
	}

	if _, err := os.Stat(flags.outputPath); err == nil && !flags.force {
		return os.ErrExist
	}

	err := os.WriteFile(flags.outputPath, data, 0600)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "Created %s\n", flags.outputPath)

	return nil
}
