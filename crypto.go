package linebot

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

func sign(body []byte, key string) string {
	hash := hmac.New(sha256.New, []byte(key))
	hash.Write(body)
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

func verify(sig string, body []byte, key string) bool {
	if sig == sign(body, key) {
		return true
	} else {
		return false
	}
}
