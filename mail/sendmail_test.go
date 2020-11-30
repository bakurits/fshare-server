package mail

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_send(t *testing.T) {
	is := assert.New(t)
	mailSender := Sender{
		Email:    "giorgi.baghdavadze@optio.ai",
		Password: "giorgi121",
	}

	err := mailSender.SendMail(
		"test mail", "test subject",
		"gbagh16@freeuni.edu.ge", "giorgi.baghdavadze@optio.ai", "bakuricucxashvili@gmail.com")
	is.NoError(err)
}
