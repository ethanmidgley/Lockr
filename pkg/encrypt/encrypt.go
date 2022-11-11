package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
)

// BytesTool | this is an interface to structure our encryption and decryption better
type BytesTool struct {
}

// NewBytesTool | creates new object to encrypt/decrypt bytes slice
func NewBytesTool() BytesTool {
	return BytesTool{}
}

// Bytes | we will encrypt a slice of bytes with key
var Bytes = NewBytesTool()

// Encrypt will encrypt data to the data given with the key passed in to it
func (b *BytesTool) Encrypt(data, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

// Decrypt will take bytes data and transfer it into the original state of data
func (b *BytesTool) Decrypt(data, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()

	nonce, cipherText := data[:nonceSize], data[nonceSize:]

	plainText, err := aesGCM.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return nil, err
	}
	return plainText, nil

}

// PasswordToKey will take a user inputted password and turn it in to a 256 bit encryption key
func PasswordToKey(password string) []byte {
	hash := sha256.Sum256([]byte(password))
	return hash[:]
}
