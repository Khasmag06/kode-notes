package decoder

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

type Decoder struct {
	gcm cipher.AEAD
}

func New(secretKey string) (*Decoder, error) {
	block, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	return &Decoder{
		gcm: gcm,
	}, nil
}

func (d *Decoder) Encrypt(data []byte) (string, error) {

	nonce := make([]byte, d.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := d.gcm.Seal(nil, nonce, data, nil)
	encrypted := append(nonce, ciphertext...)
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

func (d *Decoder) Decrypt(encrypted string) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return nil, fmt.Errorf("failed to decode ciphertext: %w", err)
	}

	nonceSize := d.gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := d.gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt ciphertext: %w", err)
	}

	return plaintext, nil
}
