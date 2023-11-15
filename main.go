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

func ReadDir(path string) (files []*File) {
	content, e := os.ReadDir(path)
	if e != nil && e.Error() != "operation not permitted" {
		err := fmt.Errorf("error reading directory: %w", e)
		log.Error().Err(err).Msg("")
	}

	for _, file := range content {
		files = append(files, &File{
			Name:  file.Name(),
			IsDir: file.IsDir(),
		})
	}

	return
}

func SearchFile(path string, fileName string) (err error) {
	if path == "" {
		var e error
		path, e = os.UserHomeDir()
		if e != nil {
			err = fmt.Errorf("no path provided and home directory unreachable: %w", e)
			return
		}
	}

	// unreachable path early exit
	info, e := os.Stat(path)
	if e != nil {
		log.Error().Err(e).Msg("")
		return
	}

	if info.Mode().Perm()&(1<<(uint(7))) == 1 {
		var subDir string
		files := ReadDir(path)
		for _, file := range files {
			if file.Name == fileName {
				fmt.Printf("Found file at: %s\n", filepath.Join(path, fileName))
				return
			}

			// recursive search
			if file.IsDir {
				// for clarity
				subDir = file.Name
				SearchFile(filepath.Join(path, subDir), fileName)
			}
		}
	}

	return
}

func main() {
	if len(os.Args) < 2 {
		log.Error().Msg("Please provide a file name")
		os.Exit(1)
	}
	fileName := os.Args[1]

	var filePath string
	if len(os.Args) > 2 {
		filePath = os.Args[2]
	} else {
		filePath = ""
	}
	if e := SearchFile(filePath, fileName); e != nil {
		err := fmt.Errorf("failed to fetch file %s, at %s: %w", fileName, filePath, e)
		log.Error().Err(err).Msg("")
	}
}
