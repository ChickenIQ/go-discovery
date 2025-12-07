package gen

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var masterFlags struct {
	serverURL string
}

func init() {
	MasterCmd.Flags().StringVarP(&masterFlags.serverURL, "server", "s", "https://discovery.chickeniq.workers.dev", "Discovery Server URL")
	RootCmd.AddCommand(MasterCmd)
}

type MasterConfig struct {
	PrivateKey string `yaml:"privateKey"`
	ServerURL  string `yaml:"serverUrl"`
}

var MasterCmd = &cobra.Command{
	Use:   "master",
	Short: "Generate Master Resources",
	RunE: func(cmd *cobra.Command, args []string) error {
		_, privateKey, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return err
		}

		data, err := yaml.Marshal(MasterConfig{
			PrivateKey: base64.StdEncoding.EncodeToString(privateKey),
			ServerURL:  masterFlags.serverURL,
		})
		if err != nil {
			return err
		}

		return output(data)
	},
}
