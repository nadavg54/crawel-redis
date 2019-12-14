package main

import (
	"flag"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/crawel-redis/config"
	"github.com/crawel-redis/redisclientwrapper"
	log "github.com/sirupsen/logrus"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logfile, err := os.OpenFile("/tmp/redis-crawler-master.log", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		os.Exit(1)
	}
	log.SetOutput(logfile)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)

	log.SetOutput(logfile)

	log.SetReportCaller(true)
}

func main() {

	address := flag.String("address", "localhost:6379", "server address")
	workersAdresses := flag.String("workers", "localhost:12345", "comma seperated list of workers")
	rootURL := flag.String("url to crawl", "http://www.nba.com", "url to crawl")

	flag.Parse()

	workersAdd := strings.Split(*workersAdresses, ",")

	urls := make(chan string, 500)

	for _, add := range workersAdd {
		go miniMaster(urls, add)
	}

	addressArr := strings.Split(*address, ":")
	port, err := strconv.Atoi(addressArr[1])
	if err != nil {
		log.Fatal("port isn't a number")
	}

	visitRoot := true

	redisClient := redisclientwrapper.ClientWrapperFactory(addressArr[0], port)
	for {
		if visitRoot {
			urls <- *rootURL
			visitRoot = false
			continue
		}
		urlsToVisist, err := redisClient.RemoveInterSectionAndRetrieve(config.URLSToVisitSetName, config.VistedURLSSetName, 300)
		for len(urlsToVisist) == 0 && err == nil {
			time.Sleep(1 * time.Second)
			urlsToVisist, err = redisClient.RemoveInterSectionAndRetrieve(config.URLSToVisitSetName, config.VistedURLSSetName, 300)
		}
		if err != nil {
			log.Fatal("master got error when trying to reach redis", err)
		}
		for _, currURL := range urlsToVisist {
			urls <- currURL
		}
	}
}

func miniMaster(urlsChan <-chan string, workerAdd string) {

	//connWriter := bufio.NewWriter(conn)

	for {
		conn, err := net.Dial("tcp", workerAdd)
		if err != nil {
			log.Warn("got error connection to worker " + err.Error())
			time.Sleep(time.Second)
			continue
		}
		url := <-urlsChan
		log.Info("sending url " + url + " from minimaster")
		conn.Write([]byte(url))
		conn.Close()
	}
}
