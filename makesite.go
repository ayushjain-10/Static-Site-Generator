package main

import (
	"flag"
	"html/template"
	"io/ioutil"
	"os"
	"strings"
)

type Post struct {
	Content string
}


func main() {
	filename := flag.String("file", "first-post.txt", "The name of the .txt file")
	flag.Parse()

	postContentBytes, err := ioutil.ReadFile(*filename)
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

	newFilename := strings.TrimSuffix(*filename, ".txt") + ".html"

	newFile, err := os.Create(newFilename)
	if err != nil {
		panic(err)
	}
	defer newFile.Close()

	if err := tmpl.Execute(newFile, post); err != nil {
		panic(err)
	}
}
