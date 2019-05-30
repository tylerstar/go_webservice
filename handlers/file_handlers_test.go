package handlers

import (
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestCreateNewFileHandler(t *testing.T) {
	e := echo.New()
	f := make(url.Values)
	f.Set("filePath", "../data/test.txt")
	f.Set("content", "Hello, World!")
	req := httptest.NewRequest(http.MethodPost, "/file", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, CreateNewFileHandler(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}

func TestGetFileContentHandler(t *testing.T) {
	e := echo.New()
	q := make(url.Values)
	q.Set("filePath", "../data/test.txt")
	req := httptest.NewRequest(http.MethodGet, "/file?"+q.Encode(), nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, GetFileContentHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestReplaceFileContentHandler(t *testing.T) {
	e := echo.New()
	f := make(url.Values)
	f.Set("filePath", "../data/test.txt")
	f.Set("content", "New Content")
	req := httptest.NewRequest(http.MethodPost, "/file", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, ReplaceFileContentHandler(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
	}
}

func TestRemoveFileHandler(t *testing.T) {
	e := echo.New()
	f := make(url.Values)
	f.Set("filePath", "../data/test.txt")
	req := httptest.NewRequest(http.MethodPost, "/file", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, RemoveFileHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}