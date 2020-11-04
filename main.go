// make_http_request.go
package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"golang.org/x/net/html"
	"strings"
	"github.com/bmaupin/go-epub"
	"bytes"
)

func main() {
	article, err := ConstructArticleFromURL("https://www.quantamagazine.org/deep-neural-networks-help-to-explain-living-brains-20201028/")
	if err != nil {
		fmt.Print(err.Error())
		log.Fatal(err)
	}
	
	epubDocument  := epub.NewEpub(article.Title)
	
	epubDocument.SetAuthor(article.Author)
	
	epubDocument.SetDescription(article.Blurb)
	
	epubDocument.AddSection(*article.Body, "Section 1", "", "")
	
	epubDocument.Write("test.epub")
}

type Article struct {
	Title           string
	Blurb           string
	Author          string
	LeadingImageUrl string
	Body			*string
}

func ConstructArticleFromURL(url string) (*Article, error) {
	document, err := goquery.NewDocument(url)

	articleRoot := document.Find("#postBody")

	title, err := ExtractDataFromeNode(articleRoot, ".post__title__title")
	if err != nil {
		return nil, err
	}

	blurb, err := ExtractDataFromeNode(articleRoot, ".post__title__excerpt")
	if err != nil {
		return nil, err
	}

	author, err := ExtractDataFromeNode(articleRoot, ".sidebar__author h3")
	if err != nil {
		return nil, err
	}
	
	imageNodes := articleRoot.Find("img").Nodes

	for _, node := range imageNodes {
		url, err = ExtractImageUrlFromImgNode(node)
		println(url, err)
	}
	
	body, err := GetArticleBody(articleRoot, ".post__content__section")
	if err != nil{
		return nil, err
	}
	
	article := Article {
		Title: title,
		Blurb: blurb,
		Author: author,
		Body: body}

	return &article, nil
}

func GetArticleBody(articleRoot *goquery.Selection, selector string) (*string, error) {
	
	bodyNodes := articleRoot.Find(selector).Nodes
	if len(bodyNodes) < 1{
		return nil, fmt.Errorf("No nodes found...")
	}
	bodyNode := bodyNodes[0]

	var b bytes.Buffer
	err := html.Render(&b, bodyNode)
	if (err != nil){
		return nil, err
	}

	body := b.String()
	
	return &body, nil
}

func ExtractDataFromeNode(s *goquery.Selection, selectors ...string) (string, error) {
	
	for _, selector := range selectors {
		s = s.Find(selector)
	}
	
	nodes := s.Nodes

	if len(nodes) < 1 {
		return "", fmt.Errorf("No nodes matched selector: " + strings.Join(selectors, " "))
	}

	return nodes[0].FirstChild.Data, nil
}

func ExtractImageUrlFromImgNode(node *html.Node) (string, error) {

	for _, attribute := range node.Attr {
		if attribute.Key == "src" {
			return attribute.Val, nil
		}
	}
	
	return "", fmt.Errorf("Could not find image URL")
}

