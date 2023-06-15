package utils

import (
	"crypto/aes"
	"crypto/cipher"
	crypto "crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"io"
	math "math/rand"
	"os"
	"time"

	"golang.org/x/crypto/pbkdf2"
)

func CopyPasswordFile(PasswordFile string, UserFolder string) string {
	DestPath := UserFolder + RandomString() + ".sanji"
	SourceFile, _ := os.Open(PasswordFile)
	DestFile, _ := os.Create(DestPath)

	io.Copy(DestFile, SourceFile)
	return DestPath
}

func RandomString() string {
	math.Seed(time.Now().UnixNano())
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	randString := make([]rune, 10)
	for i := range randString {
		randString[i] = letters[math.Intn(len(letters))]
	}
	return string(randString)
}

func DeriveKey(passphrase string) ([]byte, []byte) {
	salt := make([]byte, 12)
	crypto.Read(salt)
	return pbkdf2.Key([]byte(passphrase), salt, 1000, 32, sha256.New), salt
}

func EncryptAES(passphrase string, plaintext []byte) string {
	key, salt := DeriveKey(passphrase)
	iv := make([]byte, 12)
	crypto.Read(iv)
	b, _ := aes.NewCipher(key)
	aesgcm, _ := cipher.NewGCM(b)
	data := aesgcm.Seal(nil, iv, plaintext, nil)
	return hex.EncodeToString(salt) + hex.EncodeToString(iv) + hex.EncodeToString(data)
}
