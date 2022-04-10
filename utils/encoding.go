package utils

import "encoding/base64"

// Encode aaab
func Encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

// Decode aaa
func Decode(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}
