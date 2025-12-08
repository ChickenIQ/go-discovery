package cmd

import (
	"github.com/chickeniq/go-discovery/cmd/gen"
	"github.com/chickeniq/go-discovery/cmd/sync"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(sync.RootCmd)
	RootCmd.AddCommand(gen.RootCmd)
}

var RootCmd = &cobra.Command{
	Use:               "discovery",
	Short:             "A client for the discovery service",
	SilenceErrors:     true,
	SilenceUsage:      true,
	DisableAutoGenTag: true,
}

func Execute() error {
	return RootCmd.Execute()
}
