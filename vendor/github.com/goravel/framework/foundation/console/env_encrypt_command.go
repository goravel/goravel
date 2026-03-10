package console

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/support"
	"github.com/goravel/framework/support/convert"
	"github.com/goravel/framework/support/file"
	"github.com/goravel/framework/support/str"
)

type EnvEncryptCommand struct {
}

func NewEnvEncryptCommand() *EnvEncryptCommand {
	return &EnvEncryptCommand{}
}

// Signature The name and signature of the console command.
func (r *EnvEncryptCommand) Signature() string {
	return "env:encrypt"
}

// Description The console command description.
func (r *EnvEncryptCommand) Description() string {
	return "Encrypt an environment file"
}

// Extend The console command extend.
func (r *EnvEncryptCommand) Extend() command.Extend {
	return command.Extend{
		Category: "env",
		Flags: []command.Flag{
			&command.StringFlag{
				Name:    "key",
				Aliases: []string{"k"},
				Value:   "",
				Usage:   "Encryption key",
			},
			&command.StringFlag{
				Name:    "name",
				Aliases: []string{"n"},
				Value:   "",
				Usage:   "Encrypted environment file name",
			},
		},
	}
}

// Handle Execute the console command.
func (r *EnvEncryptCommand) Handle(ctx console.Context) error {
	key := convert.Default(ctx.Option("key"), str.Random(32))
	name := convert.Default(ctx.Option("name"), support.EnvFileEncryptPath)
	plaintext, err := os.ReadFile(support.EnvFilePath)
	if err != nil {
		ctx.Error("Environment file not found.")
		return nil
	}
	if file.Exists(name) && !ctx.Confirm("Encrypted environment file already exists, are you sure to overwrite?") {
		return nil
	}

	ciphertext, err := r.encrypt(plaintext, []byte(key))
	if err != nil {
		ctx.Error(fmt.Sprintf("Encrypt error: %v", err))
		return nil
	}

	base64Data := base64.StdEncoding.EncodeToString(ciphertext)
	err = os.WriteFile(name, []byte(base64Data), 0644)
	if err != nil {
		ctx.Error(fmt.Sprintf("Writer error: %v", err))
		return nil
	}

	ctx.Success("Environment successfully encrypted.")
	ctx.TwoColumnDetail("Key", key)
	ctx.TwoColumnDetail("Cipher", support.EnvFileEncryptCipher)
	ctx.TwoColumnDetail("Encrypted file", name)

	return nil
}

func (r *EnvEncryptCommand) encrypt(plaintext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	iv := key[:aes.BlockSize]
	plaintext = r.pkcs7Pad(plaintext)
	mode := cipher.NewCBCEncrypter(block, iv)
	ciphertext := make([]byte, len(plaintext))
	mode.CryptBlocks(ciphertext, plaintext)

	return append(iv, ciphertext...), nil
}

func (r *EnvEncryptCommand) pkcs7Pad(data []byte) []byte {
	padding := aes.BlockSize - len(data)%aes.BlockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}
