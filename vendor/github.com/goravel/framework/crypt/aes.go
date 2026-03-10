package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"

	"github.com/goravel/framework/contracts/config"
	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/errors"
	"github.com/goravel/framework/support"
	"github.com/goravel/framework/support/color"
	"github.com/goravel/framework/support/convert"
)

type AES struct {
	json foundation.Json
	key  []byte
}

// NewAES returns a new AES hasher.
func NewAES(config config.Config, json foundation.Json) (*AES, error) {
	key := config.GetString("app.key")

	// Don't use AES in artisan when the key is empty.
	if support.RuntimeMode == support.RuntimeArtisan && len(key) == 0 {
		return nil, errors.CryptAppKeyNotSet
	}

	keyLength := len(key)
	// check key length before using it
	if keyLength != 16 && keyLength != 24 && keyLength != 32 {
		color.Errorf("[Crypt] Invalid APP_KEY length. Expected 16, 24, or 32 bytes, but got %d bytes.\n", len(key))
		color.Default().Println("Please reset it using the following command:")
		color.Default().Println("go run . artisan key:generate")
		return nil, errors.CryptInvalidAppKeyLength.Args(keyLength)
	}

	return &AES{
		key:  convert.UnsafeBytes(key),
		json: json,
	}, nil
}

// EncryptString encrypts the given string, and returns the iv and ciphertext as base64 encoded strings.
func (b *AES) EncryptString(value string) (string, error) {
	block, err := aes.NewCipher(b.key)
	if err != nil {
		return "", err
	}

	plaintext := convert.UnsafeBytes(value)

	iv := make([]byte, 12)
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	ciphertext := aesgcm.Seal(nil, iv, plaintext, nil)

	var jsonEncoded []byte
	jsonEncoded, err = b.json.Marshal(map[string][]byte{
		"iv":    iv,
		"value": ciphertext,
	})

	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(jsonEncoded), nil
}

// DecryptString decrypts the given iv and ciphertext, and returns the plaintext.
func (b *AES) DecryptString(payload string) (string, error) {
	decodePayload, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return "", err
	}

	decodeJson := make(map[string][]byte)
	err = b.json.Unmarshal(decodePayload, &decodeJson)
	if err != nil {
		return "", err
	}

	// check if the json payload has the correct keys
	if _, ok := decodeJson["iv"]; !ok {
		return "", errors.CryptMissingIVKey
	}
	if _, ok := decodeJson["value"]; !ok {
		return "", errors.CryptMissingValueKey
	}

	decodeIv := decodeJson["iv"]
	decodeCiphertext := decodeJson["value"]

	block, err := aes.NewCipher(b.key)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plaintext, err := aesgcm.Open(nil, decodeIv, decodeCiphertext, nil)
	if err != nil {
		return "", err
	}

	return convert.UnsafeString(plaintext), nil
}
