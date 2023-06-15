package main

import (
	"flag"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"github.com/bregydoc/gtranslate"
	"time"
	"github.com/fatih/color"
	"fmt"
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

	start := time.Now()

	files, err := ioutil.ReadDir(".")
	if err != nil {
		panic(err)
	}

	count := 0
	totalSize := int64(0) // in bytes

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".txt") {
			processFile(file.Name())
			count++

			// Get size of the generated HTML file
			htmlFile, err := os.Stat(strings.TrimSuffix(file.Name(), ".txt") + ".html")
			if err != nil {
				panic(err)
			}

			totalSize += htmlFile.Size()
		}
	}

	elapsed := time.Since(start)
	kbSize := float64(totalSize) / 1024.0

	// Print success message
	green := color.New(color.FgGreen, color.Bold).SprintFunc()
	bold := color.New(color.Bold).SprintFunc()
	fmt.Printf("%s Generated %s pages (%.1fkB total) in %.2f seconds.\n",
		green("Success!"), bold(count), kbSize, elapsed.Seconds())
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
