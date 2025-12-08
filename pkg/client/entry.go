package client

import (
	"crypto/ed25519"
	"fmt"
)

type Entry struct {
	Member Member
	Body   Body
}

// IsValid verifies the signatures of the entry using the provided master key.
func (e *Entry) IsValid(masterKey ed25519.PublicKey) bool {
	memberSignature, err := b64Decode(e.Member.Signature)
	if err != nil {
		return false
	}

	memberData := fmt.Append([]byte(e.Member.Key), e.Member.Metadata)
	if !ed25519.Verify(masterKey, memberData, memberSignature) {
		return false
	}

	memberKey, err := b64Decode(e.Member.Key)
	if err != nil {
		return false
	}

	bodySignature, err := b64Decode(e.Body.Signature)
	if err != nil {
		return false
	}

	bodyData := fmt.Append([]byte(e.Body.Data), e.Body.Timestamp)
	return ed25519.Verify(ed25519.PublicKey(memberKey), bodyData, bodySignature)
}
