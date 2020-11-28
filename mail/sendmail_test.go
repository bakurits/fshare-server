package mail

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestService_send(t *testing.T) {
	is := assert.New(t)
	err := SendMail(
		"test mail", "test subject", "giorgi.baghdavadze@optio.ai",
		"giorgi121", "gbagh16@freeuni.edu.ge", "giorgi.baghdavadze@optio.ai")
	is.NoError(err)
}
