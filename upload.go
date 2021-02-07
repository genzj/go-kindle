package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

func upload(c echo.Context) error {
	//-----------
	// Read file
	//-----------

	// Source
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Destination
	dir, err := ioutil.TempDir("", "uploaded-")
	if err != nil {
		return err
	}

	tmpfn := filepath.Join(dir, file.Filename)

	err = func() error {
		dst, err := os.Create(tmpfn)
		if err != nil {
			return err
		}
		defer dst.Close()
		// Copy
		_, err = io.Copy(dst, src)

		return err
	}()
	if err != nil {
		return err
	}

	report := make(chan error)
	context := c.(*serverContext)
	context.ch <- &kindlePushMessage{
		from:     context.config.FromEmail,
		to:       context.config.ToEmail,
		filename: file.Filename,
		fullpath: tmpfn,
		logger:   c.Logger(),
		cleanup: func() {
			c.Logger().Debugf("clean tmp folder %s", dir)
			os.RemoveAll(dir)
		},
		report: func(err error) {
			if err != nil {
				c.Logger().Errorf(
					"sending %s from %s to %s failed: %s",
					tmpfn, context.config.FromEmail,
					context.config.ToEmail, err.Error(),
				)
			}
			c.Logger().Debugf(
				"sending %s from %s to %s successfully",
				tmpfn, context.config.FromEmail,
				context.config.ToEmail,
			)
			report <- err
		},
	}

	err = <-report
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, nil)
}
