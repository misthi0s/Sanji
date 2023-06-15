package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"

	"golang.org/x/crypto/pbkdf2"
)

func webServer(port string, key string) {
	http.HandleFunc("/favicon.ico", favicon)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", "Sanji")
		if r.Method != http.MethodPost {
			http.Error(w, "It doesn’t matter if you spend 10 thousand berries or one million berries, you should never waste food.", http.StatusMethodNotAllowed)
			return
		}

		requestContents, err := io.ReadAll(r.Body)
		if err != nil {
			return
		}
		decryptedResults, err := decryptAES(key, string(requestContents))
		if err != nil {
			fmt.Println("Error encountered: ", err)
			return
		}
		fmt.Println("Results:\n", string(decryptedResults))
	})

	if err := http.ListenAndServe(":"+port, nil); err != http.ErrServerClosed {
		fmt.Println("Unable to start HTTP server. Error: " + err.Error())
	}
}

func favicon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/x-icon")
	w.Header().Set("Cache-Control", "public, max-age=7776000")
	http.ServeFile(w, r, "../images/Sanji.ico")
}

func decryptAES(passphrase string, ciphertext string) ([]byte, error) {
	salt, err := hex.DecodeString(ciphertext[0:24])
	if err != nil {
		return nil, err
	}
	iv, _ := hex.DecodeString(ciphertext[24:48])
	data, _ := hex.DecodeString(ciphertext[48:])
	key := pbkdf2.Key([]byte(passphrase), salt, 1000, 32, sha256.New)
	b, _ := aes.NewCipher(key)
	aesgcm, _ := cipher.NewGCM(b)
	data, _ = aesgcm.Open(nil, iv, data, nil)
	return data, nil
}

func main() {
	if len(os.Args) > 2 {
		port := os.Args[1]
		passphrase := os.Args[2]
		webServer(port, passphrase)
	} else {
		fmt.Println("Usage: server_decrypt.exe <listen_port> <decryption_passphrase>")
		os.Exit(1)
	}
}
