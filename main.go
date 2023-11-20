package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
)

type File struct {
	Name  string
	IsDir bool
}

func ReadDir(path string) (files []*File) {
	// cache previously unknown unreachable paths
	dir, e := os.Open(path)
	if e != nil {
		// attempt to cache unreachable path
		cache, e := os.OpenFile("cache.txt", os.O_APPEND|os.O_WRONLY, 0644)
		if e != nil {
			err := fmt.Errorf("failed to fetch cache: %w", e)
			log.Error().Err(err).Msg("")
		}
		defer cache.Close()

		if _, e := cache.WriteString(path + "\n"); e != nil {
			err := fmt.Errorf("failed to write to cache: %w", e)
			log.Error().Err(err).Msg("")
		}

		return
	}
	defer dir.Close()

	content, e := dir.Readdir(-1)
	if e != nil {
		return
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

	if _, ok := cache[path]; ok {
		return
	}

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

	return
}

var cache = make(map[string]struct{})

// TODO implement "near" and "far" commands in place of directory arg
// TODO implement conditional .txt generation of unreachable directories -> tabulation
func main() {
	var fileName string
	if len(os.Args) < 2 {
		log.Error().Msg("please provide a file name")
		os.Exit(1)

	} else if os.Args[1] == "init" {
		// base state, assumes that all commonly searched directories will be around this level
		fileName = "main.go"
		_, e := os.Create("cache.txt")
		if e != nil {
			err := fmt.Errorf("failed to create cache.txt: %w", e)
			log.Error().Err(err).Msg("")
		}
		log.Info().Msg("initialized empty cache.txt, scanning...")

	} else {
		fileName = os.Args[1]
		cacheBytes, e := os.ReadFile("cache.txt")
		if e != nil {
			err := fmt.Errorf("failed to fetch cache: %w", e)
			log.Error().Err(err).Msg("")
			os.Exit(1)
		}

		cachePaths := strings.Split(string(cacheBytes), "\n")
		for _, path := range cachePaths {
			cache[path] = struct{}{}
		}
	}

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
