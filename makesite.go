package main

import (
	"flag"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"fmt"
	"github.com/bregydoc/gtranslate"
)

type Post struct {
	Content string
}

func main() {
	file := flag.String("file", "", "The name of the .txt file")
	dir := flag.String("dir", "", "The directory to find .txt files")
	flag.Parse()


	if *file != "" {
		processFile(*file)
	} else if *dir != "" {
		processDirectory(*dir)
	} else {
		println("Please provide either a --file or a --dir flag.")
	}
}

func processFile(filename string) {
	postContentBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	postContent := string(postContentBytes)

	translatedContent, err := gtranslate.TranslateWithParams(
		postContent,
		gtranslate.TranslationParams{
			From: "en",
			To:   "es",
		},
	)
	if err != nil {
		panic(err)
	}

	post := &Post{
		Content: translatedContent,
	}

	tmpl, err := template.ParseFiles("template.tmpl")
	if err != nil {
		panic(err)
	}

	newFilename := strings.TrimSuffix(filename, ".txt") + ".html"
	newFile, err := os.Create(newFilename)
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
