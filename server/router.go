package server

import (
	"net/http"

	"github.com/bakurits/fshare-server/db"
	"github.com/bakurits/fshare-server/mail"

	"github.com/bakurits/fshare-common/auth"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/schema"
)

var schemaDecoder = schema.NewDecoder()

// Server dependencies for endpoints
type Server struct {
	r *gin.Engine

	AuthConfig *auth.Config

	Repository    *db.Repository
	StaticFileDir string

	MailSender *mail.SenderMail
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.r.ServeHTTP(w, r)
}

// Init initializes server
func (s *Server) Init() {
	router := gin.Default()

	store := sessions.NewCookieStore([]byte("secret"))
	router.Use(sessions.Sessions("fileshare", store))

	router.GET("/", s.userExtractorMiddleware(s.homePageHandler()))

	router.GET("/login", s.loginPageHandler())
	router.POST("/login", s.loginHandler())

	router.GET("/change-password", s.userExtractorMiddleware(s.changePasswordPageHandler()))
	router.POST("/change-password", s.userExtractorMiddleware(s.changePasswordHandler()))

	router.POST("/logout", s.logoutHandler())

	router.GET("/password-recovery-request", s.passwordRecoveryRequestPageHandler())
	router.POST("/password-recovery-request", s.passwordRecoveryRequestHandler())

	router.GET("/password-recovery", s.emailPasswordRecoveryHandler())

	router.GET("/auth", s.authHandler())

	router.Static("static", s.StaticFileDir)

	router.GET("/api/token", s.getUserTokenHandler())

	s.r = router
}
