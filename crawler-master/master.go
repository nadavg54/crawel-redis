package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/crawel-redis/config"
	"github.com/crawel-redis/redisclientwrapper"
)

const ()

func main() {

	//logfile, err := os.OpenFile("/tmp/redis-crawler-master.log", os.O_RDWR|os.O_CREATE, 0666)
	//log.SetOutput(logfile)
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
		fmt.Fprintf(os.Stderr, "port isn't a number")
		os.Exit(1)
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
			os.Exit(1)
		}
		for _, currURL := range urlsToVisist {
			urls <- currURL
		}

	}

	//resp, err := http.Get("http://www.example.com")

	//root, err := html.Parse(resp.Body)
	//fmt.Println(root)
	//content, err := ioutil.ReadAll(resp.Body)

	//fmt.Println("hello")
	//fmt.Println(string(content))

	// conn, err := net.Dial("tcp", *address)
	// if err != nil {
	// 	log.Print("cant connect")
	// 	return
	// }

	// conn.Write([]byte(*rootURL))
	// conn.Close()

}

func miniMaster(urlsChan <-chan string, workerAdd string) {

	conn, err := net.Dial("tcp", workerAdd)
	//connWriter := bufio.NewWriter(conn)
	if err != nil {
		//log.Fatal("mimi master cant reach worker "+workerAdd, err)
		return
		//os.Exit(1)
	}

	for {

		url := <-urlsChan
		log.Println("sending url " + url + " from minimaster")
		conn.Write([]byte(url))
		conn.Close()
	}
}
