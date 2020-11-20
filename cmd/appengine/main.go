package main

import (
	"log"
	"net/http"

	"github.com/bakurits/fshare-server/db"
	"github.com/bakurits/fshare-server/server"

	"github.com/bakurits/fileshare/pkg/auth"
	"github.com/sethvargo/go-envconfig"
	"google.golang.org/appengine"
)

type Config struct {
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
	log.Println("starting")
	ctx := appengine.BackgroundContext()

	var conf Config
	if err := envconfig.Process(ctx, &conf); err != nil {
		log.Fatal(err)
	}

	repository, err := db.NewRepository(conf.DBDialect, conf.ConnectionString)
	if err != nil {
		log.Fatal(err)
		return
	}

	s := &server.Server{
		AuthConfig: auth.GetConfig(conf.ClientID, conf.ClientSecret, conf.Server+"/auth"),
		Repository: repository,
	}
	s.Init()

	err = http.ListenAndServe(":"+conf.Port, s)
	log.Fatal(err)
}
