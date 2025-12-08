package sync

import (
	"fmt"
	"os"

	"github.com/chickeniq/go-discovery/pkg/client"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var SyncFlags struct {
	filePath string
	data     string
}

func init() {
	RootCmd.Flags().StringVarP(&SyncFlags.filePath, "file", "f", "", "Path to Member Config")
	RootCmd.Flags().StringVarP(&SyncFlags.data, "data", "d", "", "Data to update member resources with")
	RootCmd.MarkFlagRequired("data")
	RootCmd.MarkFlagRequired("file")
}

var RootCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync Member Resources",
	RunE: func(cmd *cobra.Command, args []string) error {
		content, err := os.ReadFile(SyncFlags.filePath)
		if err != nil {
			return err
		}

		var cfg client.Config
		if err := yaml.Unmarshal(content, &cfg); err != nil {
			return err
		}

		client, err := client.NewClient(&cfg)
		if err != nil {
			return err
		}

		entries, err := client.Sync(SyncFlags.data)
		if err != nil {
			return err
		}

		if len(*entries) == 0 {
			fmt.Println("No entries found")
			return nil
		}

		for _, entry := range *entries {
			entryData, err := yaml.Marshal(entry)
			if err != nil {
				return err
			}
			fmt.Println(string(entryData))
		}

		return nil
	},
}
