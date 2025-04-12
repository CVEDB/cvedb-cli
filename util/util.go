package util

import (
	"log"
	"os"
)

type Config struct {
	User struct {
		Token         string
		TokenFilePath string
	}
	BaseUrl    string
	Dependency string
}

var (
	Cfg          Config
	SpaceName    string
	ProjectName  string
	WorkflowName string
	URL          string
)

func GetToken() string {
	if Cfg.User.Token != "" {
		return Cfg.User.Token
	}

	if Cfg.User.TokenFilePath != "" {
		token, err := os.ReadFile(Cfg.User.TokenFilePath)
		if err != nil {
			log.Fatal("Couldn't read the token file: ", err)
		}
		Cfg.User.Token = string(token)
		return Cfg.User.Token
	}

	if tokenEnv, tokenSet := os.LookupEnv("CVEDB_TOKEN"); tokenSet {
		Cfg.User.Token = tokenEnv
		return tokenEnv
	}

	log.Fatal("Cvedb authentication token not set! Use --token, --token-file, or CVEDB_TOKEN environment variable.")
	return ""
}
