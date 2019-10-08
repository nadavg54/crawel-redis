package parser

import (
	"errors"
	"io"

	"golang.org/x/net/html"
)

// type LineFilter struct{

// }

// func (*LineFilter) isValid(line string) bool{

// }

//URLParser returns all the urls inside a html page
type URLParser struct {
}

//GetURLS returns all
func (URLParser) GetURLS(htmlInput io.Reader) ([]string, error) {

	htmlRoot, err := html.Parse(htmlInput)
	var result []string

	if err != nil {
		parserError := errors.New("html parser failed with error" + err.Error())
		return nil, parserError
	}

	for currNode := htmlRoot.FirstChild; currNode != nil; {
		if isLinkElement(currNode) {
			result = append(result, getURLAttrb(currNode.Attr))
		}
	}
	return result, nil
}

func getURLAttrb(attr []html.Attribute) string {
	for i := 0; i < len(attr); i++ {
		if attr[i].Key == "href" {
			return attr[i].Val
		}
	}
	return ""
}

func isLinkElement(node *html.Node) bool {
	if node.Type == html.ElementNode {
		if node.Data == "a" {
			return true
		}
	}
	return false
}
