package main

import (
	"context"
	"log"
	"net/http"

	"github.com/bakurits/fshare-server/db"
	"github.com/bakurits/fshare-server/server"

	"github.com/bakurits/fshare-common/auth"
	"github.com/sethvargo/go-envconfig"
)

type config struct {
	CredentialsDir string `env:"credentials_dir"`

	ConnectionString string `env:"connection_string"`
	DBDialect        string `env:"db_dialect"`

	Server string `env:"server"`
	Port   string `env:"PORT"`

	ClientID     string `env:"client_id"`
	ClientSecret string `env:"client_secret"`
	ProjectID    string `env:"project_id"`
}

func main() {
	var conf config
	if err := envconfig.Process(context.Background(), &conf); err != nil {
		log.Fatal(err)
	}

	conf.DBDialect = "mysql"
	conf.CredentialsDir = "C:\\Users\\Giorgi\\GolandProjects\\fileshare\\credentials"
	conf.ConnectionString = "giorgi:giorgi121@(localhost)/test"
	conf.ClientID = "362043341673-ofd7r9v6dtjej1u3b3kg73nd65e7b6n9.apps.googleusercontent.com"
	conf.ClientSecret = "14SbmTEQZgoVeaqG-enP2jjP"
	conf.Port = "8080"
	conf.ProjectID = "fileshare-286313"
	conf.Server = "http://localhost:8080"

	repository, err := db.NewRepository(conf.DBDialect, conf.ConnectionString)
	if err != nil {
		log.Fatal(err)
		return
	}

	s := &server.Server{
		AuthConfig:    auth.GetConfig(conf.ClientID, conf.ClientSecret, conf.Server+"/auth"),
		Repository:    repository,
		StaticFileDir: "static",
	}
	s.Init()

	err = http.ListenAndServe(":"+conf.Port, s)
	log.Fatal(err)
}
