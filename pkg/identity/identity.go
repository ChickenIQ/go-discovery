package identity

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/chickeniq/go-discovery/pkg/discovery"
)

type Identity struct {
	privateKey ed25519.PrivateKey
	metadata   string
	signature  []byte

	member *discovery.Member
}

func NewIdentity(privateKey ed25519.PrivateKey, metadata string, signature []byte) *Identity {
	return &Identity{
		privateKey: privateKey,
		metadata:   metadata,
		signature:  signature,
	}
}

func (id *Identity) PrivateKey() ed25519.PrivateKey {
	return id.privateKey
}

// PublicKey returns the public key corresponding to the identity's private key.
func (id *Identity) PublicKey() ed25519.PublicKey {
	return id.PrivateKey().Public().(ed25519.PublicKey)
}

// Metadata returns the metadata associated with the identity.
func (id *Identity) Metadata() string {
	return id.metadata
}

// Signature returns the signature associated with the identity.
func (id *Identity) Signature() []byte {
	return id.signature
}

// Sign creates a digital signature for the provided data using the identity's private key.
func (id *Identity) Sign(data []byte) []byte {
	return ed25519.Sign(id.PrivateKey(), data)
}

// GetMember converts the Identity into a discovery.Member struct.
func (id *Identity) Member() *discovery.Member {
	if id.member != nil {
		return id.member
	}

	id.member = &discovery.Member{
		Key:       base64.StdEncoding.EncodeToString(id.PublicKey()),
		Metadata:  id.Metadata(),
		Signature: base64.StdEncoding.EncodeToString(id.Signature()),
	}

	return id.member
}

// NewBody constructs a Body struct containing the provided data, current timestamp, and a signature.
func (id *Identity) NewBody(data []byte) *discovery.Body {
	timestamp := time.Now().UnixMilli()
	mergedData := fmt.Append([]byte(id.Metadata()), data, timestamp)

	return &discovery.Body{
		Data:      base64.StdEncoding.EncodeToString(data),
		Timestamp: timestamp,
		Signature: base64.StdEncoding.EncodeToString(id.Sign(mergedData)),
	}
}
