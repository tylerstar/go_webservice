package utils

import (
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetAllFilePathsFromEntryPoint(t *testing.T) {
	entryPoint := "../data"

	filePaths, err := GetAllFilePathsFromEntryPoint(entryPoint)
	if err != nil {
		log.Fatalf("Failed to files from entry point: %s, error: %v", entryPoint, err)
	}
	assert.Equal(t, 2, len(filePaths))
}

func TestCountFileAverageWordLength(t *testing.T) {
	filePath := "../data/text_files/text1.txt"

	averageWordLength, err := CountFileAverageWordLength(filePath)
	if err != nil {
		log.Fatalf("Failed to get average word length from the entry point: %s, error: %v", filePath, err)
	}
	assert.Equal(t, float32(5.4329896), averageWordLength)
}

func TestCountFileAlphaChars(t *testing.T) {
	filePath := "../data/text_files/text2.txt"

	fileAlphaCharsCount, err := CountFileAlphaChars(filePath)
	if err != nil {
		log.Fatalf("Failed to get the number of alphanumeric characters of file: %s, error: %v",
			filePath, err)
	}
	assert.Equal(t, 134, fileAlphaCharsCount)
}
