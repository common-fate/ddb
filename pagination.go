package ddb

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strings"
)

// Cursor is a type alias to make pagination easier to work with
// and avoid us forgetting to decrypt a user-supplied cursor before
// including it in the database query.
type Cursor struct {
	Pk string
	Sk string
}

// String returns the plaintext value of the cursor.
func (c *Cursor) String() string {
	return fmt.Sprintf("%s:%s", c.Pk, c.Sk)
}

// Encrypt a cursor with a provided AES key. You can create a key by calling storage.Secret()
func (c *Cursor) Encrypt(secret *PaginationSecret) (string, error) {
	if secret == nil {
		return "", errors.New("PaginationSecret was nil")
	}

	cipherBlock, err := aes.NewCipher(secret.Value)
	if err != nil {
		return "", err
	}

	aead, err := cipher.NewGCM(cipherBlock)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aead.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(aead.Seal(nonce, nonce, []byte(c.String()), nil)), nil
}

// Decrypt a cursor from a user provided 'nextToken' value.
// The 'nextToken' value must be base64 encoded and encrypted with AES.
// The 'secret' argument is the AES encryption secret.
func DecryptCursor(nextToken string, secret *PaginationSecret) (*Cursor, error) {
	if secret == nil {
		return nil, errors.New("PaginationSecret was nil")
	}

	encryptData, err := base64.URLEncoding.DecodeString(nextToken)
	if err != nil {
		return nil, err
	}

	cipherBlock, err := aes.NewCipher(secret.Value)
	if err != nil {
		return nil, err
	}

	aead, err := cipher.NewGCM(cipherBlock)
	if err != nil {
		return nil, err
	}

	nonceSize := aead.NonceSize()
	if len(encryptData) < nonceSize {
		return nil, err
	}

	nonce, cipherText := encryptData[:nonceSize], encryptData[nonceSize:]
	plainData, err := aead.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(string(plainData), ":")
	if len(parts) != 2 {
		return nil, errors.New("expected cursor to be PRIMARY_KEY:SORT_KEY format")
	}

	return &Cursor{Pk: parts[0], Sk: parts[1]}, nil
}

// PaginationSecret
type PaginationSecret struct {
	Value []byte
}

// NewPaginationSecret returns a 32 bytes AES key for encrypting the cursor with.
func NewPaginationSecret() (*PaginationSecret, error) {
	key := make([]byte, 32)

	if _, err := rand.Read(key); err != nil {
		return nil, err
	}

	s := PaginationSecret{
		Value: key,
	}

	return &s, nil
}
