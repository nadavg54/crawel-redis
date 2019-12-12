package main

import (
	"flag"
	_ "io"
	"io/ioutil"
	"log"
	"net"
	"net/http"

	"github.com/crawel-redis/config"
	"github.com/crawel-redis/redisclientwrapper"

	parser "github.com/crawel-redis/crawler-parser/urls-parser"
)

//const PORT_LISTEN_TO = 1234
//const VISITED_URLS_SET_NAME = "visted"

// func crawler(urlsStream chan string, parser Parser) {
// 	for url := range urlsStream {

// 	}
// }

func main() {

	//logfile, err := os.OpenFile("/tmp/redis-crawler-worker.log", os.O_RDWR|os.O_CREATE, 0666)
	//log.SetOutput(logfile)
	address := flag.String("address", "localhost", "redis address")
	port := flag.Int("port", 6379, "port number of redis")
	redisClient := redisclientwrapper.ClientWrapperFactory(*address, *port)

	urls := make(chan string)

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		// handle error
	}

	go miniWorker(urls, redisClient)
	go miniWorker(urls, redisClient)

	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
		}
		url, err := ioutil.ReadAll(conn)
		if err != nil {

		}
		urls <- string(url)

	}

}

func miniWorker(input <-chan string, redisClient *redisclientwrapper.ClientWrapper) {
	var myParser parser.URLParser

	for {
		url := <-input
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("couldn't fetch content of page " + url)
		}
		//htmlReader := io.Reader(resp.Body)
		urls, err := myParser.GetURLS(resp.Body)
		if err != nil {
			log.Printf("coudn't parse page " + url + " with error " + err.Error())
		}
		err = redisClient.AddToSet(config.URLSToVisitSetName, urls)
		if err != nil {
			log.Printf("coudn't add to set " + config.URLSToVisitSetName + " with error " + err.Error())
		}
	}

}
