package main

import (
	"encoding/base64"
)

// func userDataHashSum(user_data string) string {
// 	v, err := base64.StdEncoding.DecodeString(user_data)
// 	if err != nil {
// 		v = []byte(user_data)
// 	}

// 	hash := sha1.Sum(v)
// 	return hex.EncodeToString(hash[:])
// }

func base64encode(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}
