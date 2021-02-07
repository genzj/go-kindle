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

	// defer os.RemoveAll(dir) // clean up

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

	context := c.(*serverContext)
	m := newBookMessage(context.config.FromEmail, context.config.ToEmail, file.Filename, tmpfn)
	context.ch <- m

	return c.JSON(http.StatusCreated, nil)
}
