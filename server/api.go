package server

import (
	"log"
	"net/http"

	"github.com/bakurits/ph"
	"github.com/gin-gonic/gin"
)

func (s *Server) getUserTokenHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		email, password, ok := c.Request.BasicAuth()
		if !ok {
			log.Println("Bad request without credentials")
			c.JSON(http.StatusUnauthorized, gin.H{})
			return
		}

		user, err := s.Repository.Users.Get(c.Request.Context(), email)
		if err != nil {
			log.Printf("Can't get user with email : %s : %v\n", email, err)
			c.JSON(http.StatusUnauthorized, gin.H{})
			return
		}

		if !ph.Compare(user.Password, password) {
			log.Printf("Password is incorrect : %s : %s\n", user.Password, password)
			c.JSON(http.StatusUnauthorized, gin.H{})
			return
		}

		c.JSON(http.StatusOK, user.Token)
	}
}
