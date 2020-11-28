package mail

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestService_send(t *testing.T) {
	is := assert.New(t)
	mailSender := SenderMail{
		Email:    "giorgi.baghdavadze@optio.ai",
		Password: "giorgi121",
	}

	err := mailSender.SendMail(
		"test mail", "test subject",
		"gbagh16@freeuni.edu.ge", "giorgi.baghdavadze@optio.ai", "bakuricucxashvili@gmail.com")
	is.NoError(err)
}
