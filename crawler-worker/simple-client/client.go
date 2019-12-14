package main

import (
	"flag"
	"log"
	"net"
)

func main() {

	address := flag.String("address", "localhost:8080", "server address")
	rootURL := flag.String("url to crawl", "http://www.nba.com", "url to crawl")

	//resp, err := http.Get("http://www.example.com")

	//root, err := html.Parse(resp.Body)
	//fmt.Println(root)
	//content, err := ioutil.ReadAll(resp.Body)

	//fmt.Println("hello")
	//fmt.Println(string(content))

	conn, err := net.Dial("tcp", *address)
	if err != nil {
		log.Print("cant connect")
		return
	}

	conn.Write([]byte(*rootURL))
	conn.Close()

	
}
