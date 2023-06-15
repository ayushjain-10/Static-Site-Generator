package main

import (
	"flag"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Post struct {
	Content string
}

func main() {
	// Add new flag for file and directory
	file := flag.String("file", "", "The name of the .txt file")
	dir := flag.String("dir", "", "The directory to find .txt files")
	flag.Parse()

	// Determine whether to process a single file or all files in a directory
	if *file != "" {
		processFile(*file)
	} else if *dir != "" {
		processDirectory(*dir)
	} else {
		println("Please provide either a --file or a --dir flag.")
	}
}

func processFile(file string) {
	postContentBytes, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	postContent := string(postContentBytes)
	post := &Post{
		Content: postContent,
	}

	tmpl, err := template.ParseFiles("template.tmpl")
	if err != nil {
		panic(err)
	}

	newfile := strings.TrimSuffix(file, ".txt") + ".html"
	newFile, err := os.Create(newfile)
	if err != nil {
		panic(err)
	}
	defer newFile.Close()

	if err := tmpl.Execute(newFile, post); err != nil {
		panic(err)
	}
}

func processDirectory(dir string) {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) == ".txt" {
			processFile(path)
		}

		return nil
	})

	if err != nil {
		panic(err)
	}
}
