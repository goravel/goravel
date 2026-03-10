package console

import (
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
)

type EnvDecryptCommand struct {
}

func NewEnvDecryptCommand() *EnvDecryptCommand {
	return &EnvDecryptCommand{}
}

// Signature The name and signature of the console command.
func (r *EnvDecryptCommand) Signature() string {
	return "env:decrypt"
}

// Description The console command description.
func (r *EnvDecryptCommand) Description() string {
	return "Decrypt an environment file"
}

// Extend The console command extend.
func (r *EnvDecryptCommand) Extend() command.Extend {
	return command.Extend{
		Category: "env",
		Flags: []command.Flag{
			&command.StringFlag{
				Name:    "key",
				Aliases: []string{"k"},
				Value:   "",
				Usage:   "Decryption key",
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
func (r *EnvDecryptCommand) Handle(ctx console.Context) error {
	key := convert.Default(ctx.Option("key"), os.Getenv("GORAVEL_ENV_ENCRYPTION_KEY"))
	name := convert.Default(ctx.Option("name"), support.EnvFileEncryptPath)
	if key == "" {
		ctx.Error("A decryption key is required.")
		return nil
	}

	ciphertext, err := os.ReadFile(name)
	if err != nil {
		ctx.Error("Encrypted environment file not found.")
		return nil
	}

	if file.Exists(support.EnvFilePath) && !ctx.Confirm("Environment file already exists, are you sure to overwrite?") {
		return nil
	}

	plaintext, err := r.decrypt(ciphertext, []byte(key))
	if err != nil {
		ctx.Error(fmt.Sprintf("Decrypt error: %v", err))
		return nil
	}

	err = os.WriteFile(support.EnvFilePath, plaintext, 0644)
	if err != nil {
		ctx.Error(fmt.Sprintf("Writer error: %v", err))
		return nil
	}

	ctx.Success("Encrypted environment successfully decrypted.")
	return nil
}

func (r *EnvDecryptCommand) decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(string(ciphertext))
	if err != nil {
		return nil, err
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	return r.pkcs7Unpad(plaintext)
}

func (r *EnvDecryptCommand) pkcs7Unpad(data []byte) ([]byte, error) {
	length := len(data)
	padding := int(data[length-1])
	return data[:length-padding], nil
}
