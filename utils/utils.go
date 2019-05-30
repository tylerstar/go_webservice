package utils

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

func GetAllFilePathsFromEntryPoint(entryPoint string) ([]string, error) {
	var filePaths []string
	err := filepath.Walk(entryPoint,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Get the status of the directory
			fi, err := os.Stat(path)
			if err != nil {
				return err
			}

			// Ensure it's a directory
			if fi.Mode().IsDir() {
				return nil
			}

			filePaths = append(filePaths, path)
			return nil
		})
	if err != nil {
		return nil, err
	}
	return filePaths, nil
}

func CountFileAlphaChars(filePath string) (int, error) {
	count := 0

	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		for _, r := range scanner.Text() {
			if unicode.IsLetter(r) {
				count++
			}
		}
	}
	return count, nil
}

func CountFileAverageWordLength(filePath string) (float32, error) {
	wordCount := 0
	totalWordLength := 0

	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	// Count number of words & total length of all words
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words := strings.Fields(scanner.Text())
		for _, word := range words {
			wordLength := len(word)
			wordCount++
			totalWordLength += wordLength
		}
	}

	// Calculate the average length of each word
	wordAverageLength := float32(totalWordLength) / float32(wordCount)

	return wordAverageLength, nil
}