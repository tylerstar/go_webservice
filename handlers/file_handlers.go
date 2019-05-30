package handlers

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"io/ioutil"
	"net/http"
	"os"
)

func CreateNewFileHandler(c echo.Context) error {
	// Get parameters
	filePath := c.FormValue("filePath")
	content := c.FormValue("content")

	// Ensure the parameters are not null
	if filePath == "" || content == "" {
		return echo.NewHTTPError(http.StatusBadRequest,
			"Parameter 'filePath' or 'content' cannot be null.")
	}

	// Check if file already exists
	if _, err := os.Stat(filePath); err == nil {
		message := fmt.Sprintf("File '%s' already exists.", filePath)
		return echo.NewHTTPError(http.StatusBadRequest, message)
	}

	// Write unified error Message to client
	// Write different internal error messages for debug
	errorMessage := fmt.Sprintf("Failed to create file: %s.", filePath)

	// Assume file's filePath was totally determined by the parameter
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open/create file %s in write only mode, error: %v", filePath, err)
		return echo.NewHTTPError(http.StatusInternalServerError, errorMessage)
	}

	// Write content into the file
	_, err = file.WriteString(content)
	if err != nil {
		log.Fatalf("Error occurred while writing file: %v", err)
		// Remove the file if failed to close the file
		if _, err := os.Stat(filePath); err == nil {
			err = os.Remove(filePath)
			log.Fatalf("Unable to remove file after failed to close it, file: %s, error: %v", filePath, err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, errorMessage)
	}

	// Ensure the file was closed properly or otherwise the content might be flushed into the file
	err = file.Close()
	if err != nil {
		log.Fatalf("Failed to close the file after writing new content, file: %s, error: %v", filePath, err)
		// Remove the file if failed to close the file
		if _, err := os.Stat(filePath); err == nil {
			err = os.Remove(filePath)
			log.Fatalf("Unable to remove file after failed to close it, file: %s, error: %v", filePath, err)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, errorMessage)
	}

	// Response
	var response struct {
		Message string `json:"Message"`
	}
	response.Message = fmt.Sprintf("File '%s' has been created.", filePath)
	return c.JSON(http.StatusCreated, &response)
}


func GetFileContentHandler(c echo.Context) error {
	// Get parameters
	filePath := c.QueryParam("filePath")

	// Ensure parameter is not null
	if filePath == "" {
		return echo.NewHTTPError(http.StatusBadRequest,
			"Parameter 'filePath' cannot be null.")
	}

	// Ensure the existence of file
	if _, err := os.Stat(filePath); err != nil {
		message := fmt.Sprintf("File '%s' doesn't exist.", filePath)
		return echo.NewHTTPError(http.StatusBadRequest, message)
	}

	// Read file's content
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read content from file: %s, error: %v",filePath, err)
		message := fmt.Sprintf("Failed to get content from file: %s", filePath)
		return echo.NewHTTPError(http.StatusInternalServerError, message)
	}

	// Response
	content := string(b)

	type fileContentResult struct {
		Content string `json:"content"`
	}

	var response struct {
		Message string				`json:"Message"`
		Result  fileContentResult	`json:"Result"`
	}
	response.Message = "Retrieved successfully."
	response.Result = fileContentResult{Content: content}
	return c.JSON(http.StatusOK, &response)
}

func ReplaceFileContentHandler(c echo.Context) error {
	// Get parameters
	filePath := c.FormValue("filePath")
	content := c.FormValue("content")

	// Ensure parameter is not null
	if filePath == "" {
		return echo.NewHTTPError(http.StatusBadRequest,
			"Parameter 'filePath' or 'content' cannot be null.")
	}

	// Ensure the existence of file
	if _, err := os.Stat(filePath); err != nil {
		message := fmt.Sprintf("File '%s' doesn't exist.", filePath)
		return echo.NewHTTPError(http.StatusBadRequest, message)
	}

	// Write to tmp file before remove the old file
	// Rename the tmp file if only everything running well
	tmpPath := filePath + ".tmp"
	file, err := os.OpenFile(tmpPath, os.O_WRONLY|os.O_CREATE, 0644)

	// Write unified error errorMessage to client
	errorMessage := fmt.Sprintf("Unable to replace content of file: %s", filePath)

	if err != nil {
		log.Fatalf("Unable to open file: %s, error: %v", tmpPath, err)
		return echo.NewHTTPError(http.StatusInternalServerError, errorMessage)
	}
	_, err = file.WriteString(content)
	if err != nil {
		log.Fatalf("Unable to write content into file: %s, error: %v", tmpPath, err)
		return echo.NewHTTPError(http.StatusInternalServerError, errorMessage)
	}
	if err = file.Close(); err != nil {
		log.Fatalf("Unable to close file: %s, error: %v", tmpPath, err)
		return echo.NewHTTPError(http.StatusInternalServerError, errorMessage)
	}

	// Rename the tmp file
	if err := os.Rename(tmpPath, filePath); err != nil {
		log.Fatalf("Failed to rename file from %s to %s, error: %v", tmpPath, filePath, err)
		return echo.NewHTTPError(http.StatusInternalServerError, errorMessage)
	}

	// Response
	var response struct {
		Message string	`json:"Message"`
	}
	response.Message = fmt.Sprintf("File '%s' content has been replaced.", filePath)
	return c.JSON(http.StatusCreated, &response)
}

func RemoveFileHandler(c echo.Context) error {
	filePath := c.FormValue("filePath")

	// Ensure parameter is not null
	if filePath == "" {
		return echo.NewHTTPError(http.StatusBadRequest,
			"Parameter 'filePath' or 'content' cannot be null.")
	}

	// Ensure the existence of file
	if _, err := os.Stat(filePath); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest,
			"File '%s' doesn't exist.", filePath)
	}

	// Remove the file
	if err := os.Remove(filePath); err != nil {
		log.Fatalf("Error occurred while removing the file '%s', error: %v", filePath, err)
		return echo.NewHTTPError(http.StatusInternalServerError,
			"Failed to remove file: %s", filePath)
	}

	// Response
	var response struct{
		Message string	`json:"Message"`
	}
	response.Message = fmt.Sprintf("File '%s' content has been removed.", filePath)
	return c.JSON(http.StatusOK, &response)
}