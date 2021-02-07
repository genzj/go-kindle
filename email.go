package main

import (
	"fmt"
	"mime"
	"path/filepath"
	"time"

	"github.com/labstack/echo/v4"
	"gopkg.in/gomail.v2"
)

var errDial = fmt.Errorf("server cannot connect Email service, get help from admin")

type kindlePushMessage struct {
	from, to           string
	filename, fullpath string
	logger             echo.Logger
	cleanup            func()
	report             func(error)
}

type emailChannel = chan<- *kindlePushMessage

func sendFile(d *gomail.Dialer, s gomail.SendCloser, logger echo.Logger, m *gomail.Message, cleanup func()) (gomail.SendCloser, error) {
	defer cleanup()

	if s != nil {
		return s, gomail.Send(s, m)
	}

	s, err := d.Dial()
	if err != nil {
		logger.Errorf("cannot dial to SMTP server %s:%d: %s", d.Host, d.Port, err)
		return nil, errDial
	}
	return s, gomail.Send(s, m)
}

func startEmailHandler(host string, port int, username, password string) emailChannel {
	ch := make(chan *kindlePushMessage)

	go func() {
		d := gomail.NewDialer(host, port, username, password)

		var s gomail.SendCloser
		var err error
		open := false

		for {
			select {
			case m, ok := <-ch:
				if !ok {
					return
				}
				mail := newBookMessage(m.from, m.to, m.filename, m.fullpath)

				if open {
					s, err = sendFile(d, s, m.logger, mail, m.cleanup)
				} else {
					s, err = sendFile(d, nil, m.logger, mail, m.cleanup)
					open = s != nil
				}
				m.report(err)

			// Close the connection to the SMTP server if no email was sent in
			// the last 30 seconds.
			case <-time.After(30 * time.Second):
				if open {
					if err := s.Close(); err != nil {
						panic(err)
					}
					open = false
				}
			}
		}
	}()

	return ch
}

func newBookMessage(from, to, filename, fullpath string) *gomail.Message {
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Upload file "+filename)
	m.SetBody("text/html", "Send file "+filename+" to kindle")
	m.Attach(
		fullpath,
		gomail.Rename(mime.QEncoding.Encode("utf-8", filepath.Base(filename))),
	)
	return m
}
