package main

import (
	"log"
	"time"

	"gopkg.in/gomail.v2"
)

type emailChannel = chan<- *gomail.Message

func startEmailHandler(host string, port int, username, password string) emailChannel {
	ch := make(chan *gomail.Message)

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
				if !open {
					if s, err = d.Dial(); err != nil {
						panic(err)
					}
					open = true
				}
				if err := gomail.Send(s, m); err != nil {
					log.Print(err)
				}
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
	m.Attach(fullpath)
	return m
}
