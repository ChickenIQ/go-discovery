package client

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"

	"github.com/chickeniq/go-discovery/pkg/discovery"
)

type Member struct {
	Key       ed25519.PublicKey
	Metadata  string
	Signature []byte
}

type Body struct {
	Data      string
	Timestamp int64
	Signature []byte
}

type Entry struct {
	Member Member
	Body   Body
}

// Validate checks the integrity and authenticity of the entry using the provided master key.
func (e *Entry) Validate(masterKey ed25519.PublicKey) error {
	memberData := fmt.Append(e.Member.Key, e.Member.Metadata)
	if !ed25519.Verify(masterKey, memberData, e.Body.Signature) {
		return fmt.Errorf("invalid member signature")
	}

	bodyData := fmt.Append([]byte(e.Member.Metadata), e.Body.Data, e.Body.Timestamp)
	if !ed25519.Verify(e.Member.Key, bodyData, e.Body.Signature) {
		return fmt.Errorf("invalid body signature")
	}

	return nil
}

// ParseEntry parses an entry from the discovery server.
func ParseEntry(e discovery.Entry) (*Entry, error) {
	memberPub, err := base64.StdEncoding.DecodeString(e.Member.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to decode member key: %w", err)
	}

	memberSignature, err := base64.StdEncoding.DecodeString(e.Member.Signature)
	if err != nil {
		return nil, fmt.Errorf("failed to decode member signature: %w", err)
	}

	bodySignature, err := base64.StdEncoding.DecodeString(e.Body.Signature)
	if err != nil {
		return nil, fmt.Errorf("failed to decode body signature: %w", err)
	}

	return &Entry{
		Member: Member{
			Key:       memberPub,
			Metadata:  e.Member.Metadata,
			Signature: memberSignature,
		},
		Body: Body{
			Data:      e.Body.Data,
			Timestamp: e.Body.Timestamp,
			Signature: bodySignature,
		},
	}, nil
}
