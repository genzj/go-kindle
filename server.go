package main

import (
	"crypto/subtle"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/bcrypt"
)

type serverContext struct {
	echo.Context
	ch     emailChannel
	config *configuration
}

func main() {
	log.SetLevel(log.DEBUG)
	config := readConfig()
	ch := startEmailHandler(config.SMTPServer, config.SMTPPort, config.SMTPUser, config.SMTPPass)
	e := echo.New()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Logger().SetLevel(log.DEBUG)
			cc := &serverContext{
				c,
				ch,
				config,
			}
			return next(cc)
		}
	})
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.BasicAuth(func(user, pass string, c echo.Context) (bool, error) {
		config := c.(*serverContext).config
		if subtle.ConstantTimeCompare([]byte(user), []byte(config.AuthUser)) == 1 &&
			bcrypt.CompareHashAndPassword([]byte(config.AuthPass), []byte(pass)) == nil {
			return true, nil
		}
		return false, nil
	}))

	csrfConfig := middleware.DefaultCSRFConfig
	csrfConfig.TokenLookup = "form:_csrf"
	e.Use(middleware.CSRFWithConfig(csrfConfig))

	e.Static("/", "public")

	e.POST("/upload", upload)
	e.OPTIONS("/upload", func(c echo.Context) error { return nil },
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{http.MethodPost},
		}))
	e.Logger.Fatal(e.Start(":1323"))
}
