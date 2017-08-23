package main

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
)

func userDataHashSum(user_data string) string {

	v, err := base64.StdEncoding.DecodeString(user_data)
	if err != nil {
		v = []byte(user_data)
	}

	hash := sha1.Sum(v)
	return hex.EncodeToString(hash[:])
}

func base64encode(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}
