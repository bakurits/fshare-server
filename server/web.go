package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bakurits/fshare-server/db"

	"github.com/bakurits/ph"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (s *Server) homePageHandler() handlerWithUser {
	type Homepage struct {
		Email         string
		IsPasswordSet bool
	}

	return func(user db.User, c *gin.Context) {
		s.executeTemplate(c.Writer, Homepage{Email: user.Email, IsPasswordSet: user.Password != ""}, true, "homepage")
	}
}

func (s *Server) loginPageHandler() gin.HandlerFunc {
	type LoginResponse struct {
		AuthLink string
	}

	return func(c *gin.Context) {
		session := sessions.Default(c)
		if session.Get(EmailSessionKey) != nil {
			c.Redirect(http.StatusSeeOther, "/")
			return
		}
		state := randToken()
		session.Set(StateSessionKey, state)
		_ = session.Save()

		s.executeTemplate(c.Writer, LoginResponse{AuthLink: s.getLoginURL(state)}, true, "login")
	}
}

func (s *Server) loginHandler() gin.HandlerFunc {
	type PostForm struct {
		Email    string `schema:"email"`
		Password string `schema:"password"`
	}
	return func(c *gin.Context) {
		session := sessions.Default(c)
		if err := c.Request.ParseForm(); err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("bad request"))
			return
		}
		var req PostForm
		if err := schemaDecoder.Decode(&req, c.Request.PostForm); err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("bad request"))
			return
		}
		u, err := s.Repository.Users.Get(req.Email)
		if err != nil {
			_ = c.AbortWithError(http.StatusUnauthorized, errors.New("user doesn't exists"))
			return
		}
		if !ph.Compare(u.Password, req.Password) {
			_ = c.AbortWithError(http.StatusUnauthorized, errors.New("incorrect credentials"))
			return
		}
		session.Set(EmailSessionKey, req.Email)
		_ = session.Save()

		c.Redirect(http.StatusSeeOther, "/")
	}
}

func (s *Server) logoutHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		session.Delete("email")
		_ = session.Save()
		c.Redirect(http.StatusSeeOther, "/")
	}
}

func (s *Server) authHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		retrievedState := session.Get(StateSessionKey)
		if retrievedState != c.Query(StateSessionKey) {
			_ = c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("invalid session state: %s", retrievedState))
			return
		}

		client, err := s.AuthConfig.ClientFromCode(c.Query("code"))
		if err != nil {
			_ = c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		_ = s.Repository.Users.Add(db.User{
			Email:    client.Email,
			Password: "",
			Token:    db.TokenStore(*client.Token),
		})

		session.Set(EmailSessionKey, client.Email)
		_ = session.Save()
		c.Redirect(http.StatusSeeOther, "/")
	}
}

func (s *Server) changePasswordPageHandler() handlerWithUser {
	return func(_ db.User, c *gin.Context) {
		s.executeTemplate(c.Writer, struct{}{}, true, "change-password")
	}
}

func (s *Server) changePasswordHandler() handlerWithUser {
	type Request struct {
		Password        string `schema:"password"`
		PasswordConfirm string `schema:"passwordConfirm"`
	}
	return func(u db.User, c *gin.Context) {
		if err := c.Request.ParseForm(); err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("bad request"))
			return
		}
		var req Request
		if err := schemaDecoder.Decode(&req, c.Request.PostForm); err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("bad request"))
			return
		}
		if req.Password == "" || req.Password != req.PasswordConfirm {
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("bad request"))
			return
		}

		passHash, err := ph.HashAndSalt(req.Password)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, errors.New("internal server error"))
			return
		}
		_ = s.Repository.Users.Update(db.User{
			Email:    u.Email,
			Password: passHash,
		})

		c.Redirect(http.StatusSeeOther, "/")
	}
}

func (s *Server) emailPasswordRecoveryHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		email := s.getEmailFromPasswordRecoveryRequest(c)
		if email == "" {
			c.Redirect(http.StatusSeeOther, "/")
			return
		}

		session.Set(EmailSessionKey, email)
		_ = session.Save()
		c.Redirect(http.StatusSeeOther, "/change-password")
	}
}

func (s *Server) passwordRecoveryRequestPageHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		s.executeTemplate(c.Writer, struct{}{}, true, "password-recovery-request")
	}
}

func (s *Server) passwordRecoveryRequestHandler() gin.HandlerFunc {
	type Request struct {
		Email string `schema:"email"`
	}
	return func(c *gin.Context) {
		if err := c.Request.ParseForm(); err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("bad request"))
			return
		}
		var req Request
		if err := schemaDecoder.Decode(&req, c.Request.PostForm); err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("bad request"))
			return
		}
		token := uuid.New().String()
		err := s.Repository.PasswordRestoreRequests.Add(db.PasswordRecoveryRequest{
			Token: token,
			Email: req.Email,
		})
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, errors.New("internal server error"))
			return
		}

		err = s.MailSender.SendMail("Password Recovery", fmt.Sprintf("https://fshare.dev/password-recovery?token=%s", token), req.Email)
		if err != nil {
			log.Printf("Mail was not send, because of %v", err)
		}
	}
}
