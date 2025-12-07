package client

import (
	"bytes"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"time"
)

func NewClient(c *Config) (*Client, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(c.Member.Key)
	if err != nil {
		return nil, err
	}

	key := ed25519.PrivateKey(keyBytes)
	pubKey := base64.StdEncoding.EncodeToString(key.Public().(ed25519.PublicKey))

	return &Client{
		Config:     c,
		PrivateKey: key,
		PublicKey:  pubKey,
	}, nil
}

func (c *Client) Update(body string) error {
	timestamp := time.Now().UnixMilli()
	bodySignature := ed25519.Sign(c.PrivateKey, fmt.Append([]byte(body), timestamp))

	data, err := json.Marshal(RequestBody{
		MasterKey: c.Config.MasterKey,

		Member: Member{
			Key:       c.PublicKey,
			Metadata:  c.Config.Member.Metadata,
			Signature: c.Config.Member.Signature,
		},

		Body: Body{
			Data:      body,
			Timestamp: timestamp,
			Signature: base64.StdEncoding.EncodeToString(bodySignature),
		},
	})
	if err != nil {
		return err
	}

	resp, err := http.Post(c.Config.ServerURL, "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(respBody))

	return nil
}
