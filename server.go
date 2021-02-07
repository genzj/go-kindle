package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type serverContext struct {
	echo.Context
	ch     emailChannel
	config *configuration
}

func main() {
	config := readConfig()
	ch := startEmailHandler(config.SMTPServer, config.SMTPPort, config.SMTPUser, config.SMTPPass)
	e := echo.New()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
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
