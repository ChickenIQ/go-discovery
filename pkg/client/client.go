package client

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"

	"github.com/chickeniq/go-discovery/pkg/discovery"
)

type IdentityProvider interface {
	Member() *discovery.Member
	NewBody(data []byte) *discovery.Body
}

type Client struct {
	MasterKey ed25519.PublicKey
	Identity  IdentityProvider
	ServerURL string
}

// New creates a new Client instance with the provided master key, server URL, and identity.
func NewClient(masterKey ed25519.PublicKey, serverURL string, identity IdentityProvider) *Client {
	return &Client{
		MasterKey: masterKey,
		ServerURL: serverURL,
		Identity:  identity,
	}
}

// Sync synchronizes the client's data with the discovery server and returns the list of entries.
func (c *Client) Sync(data []byte) (*[]Entry, error) {
	masterKey := base64.StdEncoding.EncodeToString(c.MasterKey)

	entry := discovery.Entry{
		Member: *c.Identity.Member(),
		Body:   *c.Identity.NewBody(data),
	}

	entries, err := discovery.Sync(c.ServerURL, masterKey, entry)
	if err != nil {
		return nil, err
	}

	var parsedEntries []Entry
	for _, e := range *entries {
		entry, err := ParseEntry(e)
		if err != nil {
			return nil, fmt.Errorf("failed to parse entry: %w", err)
		}

		parsedEntries = append(parsedEntries, *entry)
	}

	return &parsedEntries, nil
}
