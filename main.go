package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

type File struct {
	Name  string
	IsDir bool
}

func ReadDir(path string) (files []File, err error) {
	content, e := os.ReadDir(path)
	if e != nil {
		err = fmt.Errorf("error reading directory: %w", e)
		log.Error().Err(err).Msg("")
		return nil, err
	}

	for _, file := range content {
		files = append(files, File{
			Name:  file.Name(),
			IsDir: file.IsDir(),
		})
	}

	return
}

func SearchFile(path string, fileName string) (err error) {
	files, e := ReadDir(path)
	if e != nil {
		err = fmt.Errorf("failed to read directory: %w", e)
		log.Error().Err(err).Msg("")
		return
	}

	for _, file := range files {
		if file.Name == fileName {
			fmt.Printf("Found file at: %s\n", filepath.Join(path, fileName))
			return
		}

		// recursive search
		if file.IsDir {
			if e := SearchFile(filepath.Join(path, file.Name), fileName); e != nil {
				err = e
				return
			}
		}
	}

	return
}

func main() {
	if len(os.Args) < 3 {
		log.Error().Msg("Please provide a path and file name")
		os.Exit(1)
	}

	filePath := os.Args[1]
	fileName := os.Args[2]
	if e := SearchFile(filePath, fileName); e != nil {
		err := fmt.Errorf("failed to fetch file %s, at %s: %w", fileName, filePath, e)
		log.Error().Err(err).Msg("")
	}
}
