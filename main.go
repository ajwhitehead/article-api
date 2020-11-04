// make_http_request.go
package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net"
	"html"
)

func main() {
	ConstructArticleFromURL("https://www.quantamagazine.org/deep-neural-networks-help-to-explain-living-brains-20201028/")

	//TODO - set up source control. Work out how to have packages. Construct package to write epub file using go-epub https://pkg.go.dev/github.com/bmaupin/go-epub
}

type Article struct {
	title           string
	blurb           string
	author          string
	leadingImageUrl string
	sections        []Section
}

type Section struct {
}

type Item struct {
	itemType ItemType
	content  string
}

type ItemType uint32

const (
	paragraph ItemType = iota
	image
)

func ConstructArticleFromURL(url string) {
	document, err := goquery.NewDocument(url)

	articleRoot := document.Find("#postBody")

	title, err := ExtractDataFromSingleNode(articleRoot, ".post__title__title")
	if err != nil {
		log.Fatal(err)
	}
	
	blurb, err := ExtractDataFromSingleNode(articleRoot, ".post__title__excerpt")
	if err != nil {
		log.Fatal(err)
	}

	author, err := ExtractDataFromSingleNode(articleRoot, ".mv05")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(title)
	fmt.Println(blurb)

	imageNodes := articleRoot.Find("img").Nodes

	for _, node := range imageNodes {
		println(ExtractImageUrlFromImgNode(node))
	}

	contentRoot := articleRoot.Find(".post__content__section")

	paragraphs := contentRoot.Find("p").Nodes

	for _, paragraph := range paragraphs {
		fmt.Println(paragraph.FirstChild.Data)
	}
}

func ExtractDataFromSingleNode(s *goquery.Selection, selector string) (string, error) {

	node := s.Find(selector).Nodes

	if len(node) > 1 {
		return "", fmt.Errorf("More than one node matched selector" + selector)
	}

	return node[0].FirstChild.Data, nil
}

func ExtractImageUrlFromImgNode(a *Node) (string, error) {

	for _, attribute := range node.Attr {
		if attribute.Key == "src" {
			return attribute.Val
		}
	}
}

