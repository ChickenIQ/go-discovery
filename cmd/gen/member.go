package gen

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/chickeniq/go-discovery/pkg/client"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func init() {
	MemberCmd.Flags().StringVarP(&memberFlags.filePath, "file", "f", "", "Path to Master Config")
	MemberCmd.Flags().StringVarP(&memberFlags.metadata, "metadata", "m", "", "Metadata for member")
	MemberCmd.MarkFlagRequired("file")
	RootCmd.AddCommand(MemberCmd)
}

var memberFlags struct {
	metadata string
	filePath string
}

var MemberCmd = &cobra.Command{
	Use:   "member",
	Short: "Generate Member Resources",
	RunE: func(cmd *cobra.Command, args []string) error {
		content, err := os.ReadFile(memberFlags.filePath)
		if err != nil {
			return err
		}

		var cfg *MasterConfig
		if err := yaml.Unmarshal(content, &cfg); err != nil {
			return err
		}

		decoded, err := base64.StdEncoding.DecodeString(cfg.PrivateKey)
		if err != nil {
			return err
		}

		masterPriv := ed25519.PrivateKey(decoded)
		masterPub := masterPriv.Public().(ed25519.PublicKey)

		memberPub, memberPriv, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return err
		}

		encodedPub := base64.StdEncoding.EncodeToString(memberPub)
		memberSig := ed25519.Sign(masterPriv, fmt.Append([]byte(encodedPub), memberFlags.metadata))

		data, err := yaml.Marshal(client.Config{
			MasterKey: base64.StdEncoding.EncodeToString(masterPub),
			ServerURL: cfg.ServerURL,
			Member: client.Member{
				Key:       base64.StdEncoding.EncodeToString(memberPriv),
				Metadata:  memberFlags.metadata,
				Signature: base64.StdEncoding.EncodeToString(memberSig),
			},
		})
		if err != nil {
			return err
		}

		return output(data)
	},
}
