package update

import (
	"os"

	"github.com/chickeniq/go-discovery/pkg/client"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var UpdateFlags struct {
	filePath string
	data     string
}

func init() {
	RootCmd.Flags().StringVarP(&UpdateFlags.filePath, "file", "f", "", "Path to Member Config")
	RootCmd.Flags().StringVarP(&UpdateFlags.data, "data", "d", "", "Data to update member resources with")
	RootCmd.MarkFlagRequired("data")
	RootCmd.MarkFlagRequired("file")
}

var RootCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Member Resources",
	RunE: func(cmd *cobra.Command, args []string) error {
		content, err := os.ReadFile(UpdateFlags.filePath)
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

		return client.Update(UpdateFlags.data)
	},
}
