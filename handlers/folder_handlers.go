package handlers

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"go_webservice/utils"
	"net/http"
	"os"
	"strconv"
)

// QueryTarget
//
//·0      Total number of files in that folder.
//·1      Average number of alphanumeric characters per text file (and standard deviation) in that folder.
//·2      Average word length (and standard deviation) in that folder.
//·3      Total number of bytes stored in that folder.
const (
	fileNumber                           = 0
	AverageNumberOfAlphaCharsPerTextFile = 1
	AverageWordLengthPerTextFile         = 2
	TotalNumberOfBytes                   = 3
)


func GetFolderStatsHandler(c echo.Context) error {
	entryPoint := c.QueryParam("entryPoint")
	queryTarget := c.QueryParam("queryTarget")

	// Ensure parameter is not null
	if entryPoint == "" || queryTarget == "" {
		return echo.NewHTTPError(http.StatusBadRequest,
			"Parameter 'entryPoint' or 'queryTarget' cannot be null.")
	}

	// Ensure the existence of file
	if _, err := os.Stat(entryPoint); err != nil {
		log.Warnf("Entry point '%s' doesn't exists.", entryPoint)
		return echo.NewHTTPError(http.StatusConflict,
			"Entry point '%s' doesn't exist.", entryPoint)
	}

	// Get the status of the directory
	fi, err := os.Stat(entryPoint)
	if err != nil {
		log.Fatalf("Failed to get stats of entryPoint: %s, error: %v", entryPoint, err)
		return echo.NewHTTPError(http.StatusInternalServerError,
			"Currently unable to get stats of entry point '%s'.", entryPoint)
	}

	// Ensure it's a directory
	if !fi.Mode().IsDir() {
		log.Warnf("EntryPoint '%s' is not directory", entryPoint)
		return echo.NewHTTPError(http.StatusBadRequest,
			"Entry point '%s' is not a directory", entryPoint)
	}

	// Ensure the value of queryTarget is valid
	queryNumber, err := strconv.Atoi(queryTarget)
	if err != nil {
		log.Warnf("Failed to convert '%s' to number", queryTarget)
		return echo.NewHTTPError(http.StatusBadRequest,
			"Invalid value, parameter 'queryTarget' expect a int from 0 ~ 3, got %s", queryTarget)
	}

	// Get all the files first
	filePaths, err := utils.GetAllFilePathsFromEntryPoint(entryPoint)
	if err != nil {
		log.Fatalf("Error occurred while listing file from the entry point: %s, error: %v",
			entryPoint, err)
		return echo.NewHTTPError(http.StatusInternalServerError,
			"Failed to list the filePaths from the entry point: %s", entryPoint)
	}

	// Response might be different according to the value of queryTarget
	var response interface{}

	switch queryNumber {
	case fileNumber:
		response, err = CountFilesFromEntryPoint(entryPoint, filePaths)
		if err != nil {
			return err
		}
	case AverageNumberOfAlphaCharsPerTextFile:
		response, err = CountAverageNumberOfAlphaCharsPerTextFile(entryPoint, filePaths)
		if err != nil {
			return err
		}
	case AverageWordLengthPerTextFile:
		response, err = CountAverageWordLengthPerTextFile(entryPoint, filePaths)
		if err != nil {
			return err
		}
	case TotalNumberOfBytes:
		response, err = CountTotalNumberOfBytes(entryPoint, filePaths)
		if err != nil {
			return err
		}
	default:
		log.Warnf("Got invalid value from parameter 'queryTarget', got %s", queryTarget)
		return echo.NewHTTPError(http.StatusBadRequest,
			"Invalid value, parameter 'queryTarget' expect a int from 0 ~ 3, got %s", queryTarget)
	}

	return c.JSON(http.StatusOK, &response)
}

func CountFilesFromEntryPoint(entryPoint string, filePaths []string) (interface{}, error) {
	fileCount := len(filePaths)

	var response struct {
		Message string			`json:"message"`
		Result  map[string]int	`json:"result"`
	}

	response.Message = fmt.Sprintf("Successfully calculate the number of file from the entry point: %s.",
		entryPoint)
	response.Result = make(map[string]int)
	response.Result["fileCount"] = fileCount
	return response, nil
}

func CountAverageNumberOfAlphaCharsPerTextFile(entryPoint string, filePaths []string) (interface{}, error) {
	// Get number of alpha chars per file
	fileAlphaCharsCountMap := make(map[string]int)
	totalFileAlphaCharsCount := 0

	var alphaCharsNumber int
	var err error
	for _, filePath := range filePaths {
		alphaCharsNumber, err = utils.CountFileAlphaChars(filePath)
		if err != nil {
			return nil, echo.NewHTTPError(http.StatusInternalServerError,
				"Failed to count the alpha characters in file: %s, error: %v", filePath, err)
		}
		fileAlphaCharsCountMap[filePath] = alphaCharsNumber
		totalFileAlphaCharsCount += alphaCharsNumber
	}

	// Calculate the standard deviation
	fileCount := len(filePaths)
	charsCountPerFileStandDeviation := float32(totalFileAlphaCharsCount) / float32(fileCount)

	type fileAlphaCharsCountResult struct {
		FileStats         map[string]int  `json:"fileStats"`
		StandardDeviation float32		  `json:"standardDeviation"`
	}

	var response struct {
		Result  fileAlphaCharsCountResult  `json:"result"`
		Message string					   `json:"message"`
	}

	// Response
	response.Message = fmt.Sprintf("Calculate the number of alphanumeric characters successfully from the entry point: %s",
		entryPoint)
	response.Result = fileAlphaCharsCountResult{
		FileStats:         fileAlphaCharsCountMap,
		StandardDeviation: charsCountPerFileStandDeviation}
	return response, nil
}

func CountAverageWordLengthPerTextFile(entryPoint string, filePaths []string) (interface{}, error) {
	fileAverageWordLengthMap := make(map[string]float32)
	var totalFileAverageWordLength float32

	for _, filePath := range filePaths {
		fileAverageWordLength, err := utils.CountFileAverageWordLength(filePath)
		if err != nil {
			log.Fatalf("Error occurred while calculating the average word length in file: %s, error: %v",
				filePath, err)
			return nil, echo.NewHTTPError(http.StatusInternalServerError,
				"Failed to count the average world length per text file from the entry point: %s", entryPoint)
		}
		fileAverageWordLengthMap[filePath] = fileAverageWordLength
		totalFileAverageWordLength += fileAverageWordLength
	}

	// Calculate the standard deviation
	fileCount := len(filePaths)
	fileAverageWordLengthStandardDeviation := totalFileAverageWordLength / float32(fileCount)

	// Response
	type fileAverageWordLengthResult struct {
		FileStats 		  map[string]float32  `json:"fileStats"`
		StandardDeviation float32			  `json:"standardDeviation"`
	}

	var response struct {
		Result  fileAverageWordLengthResult  `json:"result"`
		Message string						 `json:"message"`
	}

	response.Message = fmt.Sprintf(
		"Successfully calculate the average word length per text file from the entry point: %s",
		entryPoint)
	response.Result = fileAverageWordLengthResult{
		FileStats: fileAverageWordLengthMap,
		StandardDeviation: fileAverageWordLengthStandardDeviation}
	return response, nil
}

func CountTotalNumberOfBytes(entryPoint string, filePaths []string) (interface{}, error) {
	// Get file size (number of bytes) of each file
	fileSizeMap := make(map[string]int64)
	var totalNumberOfBytes int64

	for _, filePath := range filePaths {
		file, err := os.Open(filePath)
		if err != nil {
			log.Fatalf("Failed to open the file: %s, error: %v", filePath, err)
			return nil, echo.NewHTTPError(http.StatusInternalServerError,
				"Failed to read the file from the entry point: %s", entryPoint)
		}

		fi, err := file.Stat()
		if err != nil {
			log.Fatalf("Error occurred while getting stats of file: %s, error: %v", filePath, err)
			return nil, echo.NewHTTPError(http.StatusInternalServerError,
				"Failed to count the total bytes from the entry point: %s", entryPoint)
		}
		fileSizeMap[filePath] = fi.Size()
		totalNumberOfBytes += fi.Size()
	}

	// Response
	type totalNumberOfBytesResult struct {
		TotalBytesCount int64	`json:"totalBytesCount"`
	}

	var response struct {
		Result  totalNumberOfBytesResult  `json:"result"`
		Message string                    `json:"message"`
	}

	response.Message = fmt.Sprintf("Successfully to calculate the total number of bytes from the entry point: %s",
		entryPoint)
	response.Result = totalNumberOfBytesResult{TotalBytesCount: totalNumberOfBytes}
	return response, nil
}