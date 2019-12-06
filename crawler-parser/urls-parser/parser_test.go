package parser

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"
)

func Test_parse(t *testing.T) {

	var myParser URLParser

	file, _ := ioutil.ReadFile("htmlToParse.html")

	result, _ := myParser.GetURLS(bytes.NewReader(file))
	fmt.Println("result is " + result[1])
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println("result size \n %n", len(result))
	// }

}
