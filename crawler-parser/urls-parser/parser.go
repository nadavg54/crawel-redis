package parser

import (
	"errors"
	"io"

	"golang.org/x/net/html"
)

//URLParser returns all the urls inside a html page
type URLParser struct {
}

//GetURLS returns all
func (URLParser) GetURLS(htmlInput io.Reader) ([]string, error) {

	result := []string{}
	htmlRoot, err := html.Parse(htmlInput)
	//result := make([]string, 1000)

	if err != nil {
		parserError := errors.New("html parser failed with error" + err.Error())
		return nil, parserError
	}

	finalResult := traverseHTMLTree(htmlRoot, result)
	return finalResult, nil
}

func traverseHTMLTree(node *html.Node, result []string) []string {

	if node == nil {
		return nil
	}
	if isLinkElement(node) {
		currlink, shouldUse := getURLAttrb(node.Attr)
		if shouldUse {

			result = append(result, currlink)
		}

	}

	for currNode := node.FirstChild; currNode != nil; currNode = currNode.NextSibling {
		result = traverseHTMLTree(currNode, result)
	}
	return result
}

func getURLAttrb(attr []html.Attribute) (string, bool) {
	for i := 0; i < len(attr); i++ {
		if attr[i].Key == "href" {
			return attr[i].Val, true
		}
	}
	return "", false
}

func isLinkElement(node *html.Node) bool {
	if node.Type == html.ElementNode {
		if node.Data == "a" {
			return true
		}
	}
	return false
}
