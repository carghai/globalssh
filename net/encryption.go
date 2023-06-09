package net

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
	"log"
)

const DecryptError = "decrypt Failed, invalid key"

func (net Net) encrypt(text string) (string, error) {
	if net.EncryptionKey == nil {
		return text, nil
	}
	plaintext := []byte(text)

	// Generate a random initialization vector (IV)
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	// Pad the plaintext to a multiple of the block size
	paddedPlaintext := pad(plaintext, aes.BlockSize)

	// Create a new CBC mode cipher using the block and IV
	mode := cipher.NewCBCEncrypter(net.EncryptionKey, iv)

	// Encrypt the padded plaintext
	ciphertext := make([]byte, len(paddedPlaintext))
	mode.CryptBlocks(ciphertext, paddedPlaintext)

	// Combine the IV and ciphertext and return the result
	return string(append(iv, ciphertext...)), nil
}

func (net Net) decrypt(text string) (string, error) {
	if net.EncryptionKey == nil {
		return text, nil
	}
	ciphertext := []byte(text)

	// Extract the IV from the ciphertext
	iv := ciphertext[:aes.BlockSize]

	// Extract the actual ciphertext
	ciphertext = ciphertext[aes.BlockSize:]

	// Create a new CBC mode cipher using the block and IV
	mode := cipher.NewCBCDecrypter(net.EncryptionKey, iv)

	// Decrypt the ciphertext
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	// Remove padding from the decrypted plaintext
	plaintext, err := unpad(plaintext)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// Pad the input to a multiple of the block size using PKCS7 padding
func pad(data []byte, blockSize int) []byte {
	padding := blockSize - (len(data) % blockSize)
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// Remove PKCS7 padding from the input
func unpad(data []byte) ([]byte, error) {
	padding := int(data[len(data)-1])
	if padding > len(data) {
		return nil, errors.New(DecryptError)
	}
	return data[:len(data)-padding], nil
}

// Pad the key to 32 bytes
func padKey(key []byte) []byte {
	paddedKey := make([]byte, 32)
	copy(paddedKey, key)
	return paddedKey
}

func NewKey(key string) cipher.Block {
	var EncryptionKey cipher.Block
	if key != "" {
		paddedKey := padKey([]byte(key))
		key, err := aes.NewCipher(paddedKey)
		if err != nil {
			log.Fatal("Failed To Create Encryption Key, Please fix it: ", err)
		}
		EncryptionKey = key
	} else {
		EncryptionKey = nil
	}
	return EncryptionKey
}
