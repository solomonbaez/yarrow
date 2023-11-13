package main

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
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

func FileDisplay(files []File) {
	for _, f := range files {
		fmt.Printf("%s\t%s\n", f.Name, map[bool]string{true: "DIR", false: "FILE"}[f.IsDir])
	}
}

func main() {
	files, e := ReadDir(".")
	if e != nil {
		err := fmt.Errorf("fatal: %w", e)
		log.Fatal().Err(err).Msg("")
		os.Exit(1)
	}

	FileDisplay(files)
}
