package main

import (
	"bytes"
	"fmt"
	"github.com/gomarkdown/markdown"
	jsoniter "github.com/json-iterator/go"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
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

func parseArticleContent(articlePath string) string {
	log.Println("- loading article from", articlePath)
	data, err := ioutil.ReadFile(articlePath)
	if err != nil {
		log.Fatal("failed to read article contents", err)
	}
	return string(markdown.ToHTML(data, nil, nil))
}

func compileBlog(config BlogConfig) {
	files, err := template.ParseFiles(config.Template)
	if err != nil {
		log.Fatal("failed to parse template", err)
	}

	for _, article := range config.Articles {
		var output bytes.Buffer

		// create new file for article
		files.Execute(&output, map[string]interface{}{
			"title": article.Title,
			"articleContent": template.HTML(parseArticleContent(article.Source)),
			"publishDate": article.PublishedAt,
		})

		fileName := strings.TrimSuffix(article.Source, filepath.Ext(article.Source))
		ioutil.WriteFile(fmt.Sprintf("%s.html", fileName), output.Bytes(), 0644)
	}
}

func main() {
	args := os.Args
	if len(args) <= 1 {
		fmt.Println("Give me my config file!")
		return
	}

	configFile := args[1]
	log.Println("Loading config file from", configFile)

	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal("failed to load blog config", err)
	}

	var config BlogConfig
	if err := jsoniter.Unmarshal(data, &config); err != nil {
		log.Fatal("failed to read blog config", err)
	}

	compileBlog(config)

	//output := markdown.ToHTML(md, nil, nil)
	//fmt.Println(output)
}
