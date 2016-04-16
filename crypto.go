package linebot

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

func sign(body []byte, key string) []byte {
	hash := hmac.New(sha256.New, []byte(key))
	hash.Write(body)
	return hash.Sum(nil)
}

func verify(sig string, body []byte, key string) bool {
	sb, err := base64.StdEncoding.DecodeString(sig)
	if err != nil {
		return false
	}
	if hmac.Equal(sb, sign(body, key)) {
		return true
	} else {
		return false
	}
}
