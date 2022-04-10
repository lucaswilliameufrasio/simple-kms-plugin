package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

// Encrypt method is to encrypt or hide any classified text
func Encrypt(text, MySecret string) (string, error) {
	secretBuffer := []byte(MySecret)
	block, err := aes.NewCipher(secretBuffer)
	if err != nil {
		return "", err
	}

	plainText := []byte(text)
	cipherText := make([]byte, len(plainText))

	iv := make([]byte, aes.BlockSize)
	fmt.Println(len(iv))
	cfb := cipher.NewCFBEncrypter(block, iv)

	cfb.XORKeyStream(cipherText, plainText)

	return Encode(cipherText), nil
}

// Decrypt method is to extract back the encrypted text
func Decrypt(text, MySecret string) (string, error) {
	block, err := aes.NewCipher([]byte(MySecret))
	if err != nil {
		return "", err
	}
	cipherText := Decode(text)
	iv := make([]byte, aes.BlockSize)
	cfb := cipher.NewCFBDecrypter(block, iv)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)
	return string(plainText), nil
}
