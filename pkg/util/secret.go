package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"hash"
)

// VerifyGithubSignature verifies if the given signature is correct or not for a webhook
func VerifyGithubSignature(secret, body []byte, expected string) bool {
	expectedBytes, _ := hex.DecodeString(expected)
	h := hmac.New(func() hash.Hash {
		return sha256.New()
	}, secret)
	h.Write(body)
	return subtle.ConstantTimeCompare(h.Sum(nil), expectedBytes) == 1
}
