package handlers

import (
	"encoding/json"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestCountFilesFromEntryPoint(t *testing.T) {
	e := echo.New()
	q := make(url.Values)
	q.Set("entryPoint", "../data")
	q.Set("queryTarget", "0")
	req := httptest.NewRequest(http.MethodGet, "/file?"+q.Encode(), nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, GetFolderStatsHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var response struct {
			Message string			`json:"message"`
			Result  map[string]int	`json:"result"`
		}

		err := json.Unmarshal([]byte(strings.TrimSpace(rec.Body.String())), &response)
		if err != nil {
			log.Fatalf("Failed to parse as json, error: %v", err)
		}
		assert.Equal(t, 2, response.Result["fileCount"])
	}
}

func TestCountAverageNumberOfAlphaCharsPerTextFile(t *testing.T) {
	e := echo.New()
	q := make(url.Values)
	q.Set("entryPoint", "../data")
	q.Set("queryTarget", "1")
	req := httptest.NewRequest(http.MethodGet, "/file?"+q.Encode(), nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, GetFolderStatsHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		//log.Info(rec.Body.String())

		type fileAverageWordLengthResult struct {
			FileStats 		  map[string]float32  `json:"fileStats"`
			StandardDeviation float32			  `json:"standardDeviation"`
		}

		var response struct {
			Result  fileAverageWordLengthResult  `json:"result"`
			Message string						 `json:"message"`
		}

		err := json.Unmarshal([]byte(strings.TrimSpace(rec.Body.String())), &response)
		if err != nil {
			log.Fatalf("Failed to parse as json, error: %v", err)
		}
		assert.Equal(t, float32(574), response.Result.StandardDeviation)
	}
}

func TestCountAverageWordLengthPerTextFile(t *testing.T) {
	e := echo.New()
	q := make(url.Values)
	q.Set("entryPoint", "../data")
	q.Set("queryTarget", "2")
	req := httptest.NewRequest(http.MethodGet, "/file?"+q.Encode(), nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, GetFolderStatsHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		//log.Info(rec.Body.String())

		type fileAverageWordLengthResult struct {
			FileStats 		  map[string]float32  `json:"fileStats"`
			StandardDeviation float32			  `json:"standardDeviation"`
		}

		var response struct {
			Result  fileAverageWordLengthResult  `json:"result"`
			Message string						 `json:"message"`
		}
		err := json.Unmarshal([]byte(strings.TrimSpace(rec.Body.String())), &response)
		if err != nil {
			log.Fatalf("Failed to parse as json, error: %v", err)
		}
		assert.Equal(t, float32(5.312649), response.Result.StandardDeviation)
	}
}

func TestCountTotalNumberOfBytes(t *testing.T) {
	e := echo.New()
	q := make(url.Values)
	q.Set("entryPoint", "../data")
	q.Set("queryTarget", "3")
	req := httptest.NewRequest(http.MethodGet, "/file?"+q.Encode(), nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, GetFolderStatsHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		//log.Info(rec.Body.String())

		var response struct {
			Result  map[string]int64  `json:"result"`
			Message string            `json:"message"`
		}

		err := json.Unmarshal([]byte(strings.TrimSpace(rec.Body.String())), &response)
		if err != nil {
			log.Fatalf("Failed to parse as json, error: %v", err)
		}
		assert.Equal(t, int64(1407), response.Result["totalBytesCount"])
	}
}