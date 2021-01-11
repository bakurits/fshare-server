package main

import (
	"context"
	"log"
	"net/http"

	"github.com/bakurits/fshare-server/db"
	"github.com/bakurits/fshare-server/mail"
	"github.com/bakurits/fshare-server/server"

	"cloud.google.com/go/firestore"
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

	Email         string `env:"email"`
	EmailPassword string `env:"email_password"`
}

func main() {
	var conf config
	if err := envconfig.Process(context.Background(), &conf); err != nil {
		log.Fatal(err)
	}

	client, err := firestore.NewClient(context.Background(), conf.ProjectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer func() {
		err := client.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	repository, err := db.NewRepository(client)
	if err != nil {
		log.Fatal(err)
		return
	}

	s := &server.Server{
		AuthConfig:    auth.GetConfig(conf.ClientID, conf.ClientSecret, conf.Server+"/auth"),
		Repository:    repository,
		StaticFileDir: "static",
		MailSender: &mail.Sender{
			Email:    conf.Email,
			Password: conf.EmailPassword,
		},
	}
	s.Init()

	err = http.ListenAndServe(":"+conf.Port, s)
	log.Fatal(err)
}
