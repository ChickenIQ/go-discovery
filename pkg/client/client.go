package client

import (
	"bytes"
	"crypto/ed25519"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"time"
)

type Client struct {
	Config     *Config
	PublicKey  string
	PrivateKey ed25519.PrivateKey
}

// NewClient creates a new Client instance from the provided configuration.
func NewClient(c *Config) (*Client, error) {
	keyBytes, err := b64Decode(c.Member.Key)
	if err != nil {
		return nil, err
	}

	key := ed25519.PrivateKey(keyBytes)
	pubKey := b64Encode(key.Public().(ed25519.PublicKey))

	return &Client{
		Config:     c,
		PrivateKey: key,
		PublicKey:  pubKey,
	}, nil
}

// Sync sends the provided body to the server and retrieves the list of entries, verifying their signatures.
func (c *Client) Sync(body string) (*[]Entry, error) {
	timestamp := time.Now().UnixMilli()
	bodySignature := ed25519.Sign(c.PrivateKey, fmt.Append([]byte(c.Config.Member.Signature), body, timestamp))

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
			Signature: b64Encode(bodySignature),
		},
	})
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(c.Config.ServerURL, "application/json", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var entries []Entry
	if err := json.Unmarshal(respBody, &entries); err != nil {
		return nil, err
	}

	masterKey, err := b64Decode(c.Config.MasterKey)
	if err != nil {
		return nil, err
	}

	for i := range entries {
		if !entries[i].IsValid(ed25519.PublicKey(masterKey)) {
			return nil, fmt.Errorf("invalid entry for member: %s", entries[i].Member.Key)
		}
	}

	return &entries, nil
}
