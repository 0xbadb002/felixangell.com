package main

import (
	"fmt"
	"log"
	"os"
)

// BlogArticle ...
type BlogArticle struct {
	Source      string `json:"source"`
	Title       string `json:"title"`
	PublishedAt string `json:"publishedAt"`
}

// BlogConfig ...
type BlogConfig struct {
	Template string        `json:"template"`
	Articles []BlogArticle `json:"articles"`
}

func main() {
	args := os.Args
	if len(args) <= 1 {
		fmt.Println("Give me my config file!")
		return
	}

	configFile := args[1]
	log.Println("Loading config file from", configFile)

	//output := markdown.ToHTML(md, nil, nil)
	//fmt.Println(output)
}
