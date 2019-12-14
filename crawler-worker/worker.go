package main

import (
	"flag"
	_ "io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/crawel-redis/config"
	"github.com/crawel-redis/redisclientwrapper"
	log "github.com/sirupsen/logrus"

	parser "github.com/crawel-redis/crawler-parser/urls-parser"
)

//const PORT_LISTEN_TO = 1234
//const VISITED_URLS_SET_NAME = "visted"

// func crawler(urlsStream chan string, parser Parser) {
// 	for url := range urlsStream {

// 	}
// }

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logfile, err := os.OpenFile("/tmp/redis-crawler-worker.log", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		os.Exit(1)
	}
	log.SetOutput(logfile)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)

	log.SetReportCaller(true)
}

func main() {

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

	i := 0
	for {
		conn, err := ln.Accept()
		log.Info("worker iteration: " + strconv.Itoa(i))
		if err != nil {
			// handle error
			log.Fatal("failed to get url: " + err.Error())
		}
		url, err := ioutil.ReadAll(conn)
		conn.Close()
		if err != nil {
			log.Warn("failed to read url content: " + err.Error())
			continue
		}
		urls <- string(url)
	}

}

func miniWorker(input <-chan string, redisClient *redisclientwrapper.ClientWrapper) {
	var myParser parser.URLParser

	for {
		url := <-input
		log.Info("got url" + url)
		resp, err := http.Get(url)
		if err != nil {
			log.Warn("couldn't fetch content of page " + url)
			continue
		}

		urls, err := myParser.GetURLS(resp.Body)
		if err != nil {
			log.Warn("coudn't parse page " + url + " with error " + err.Error())
			continue
		}
		log.Printf("parse results: " + strconv.Itoa(len(urls)))
		err = redisClient.AddToMultipleSets(config.URLSToVisitSetName, urls, config.VistedURLSSetName, []string{url})
		if err != nil {
			log.Warn("coudn't add to set " + config.URLSToVisitSetName + " with error " + err.Error())
		}
	}
}
