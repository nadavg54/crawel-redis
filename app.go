package main

import (
	"flag"
	"github.com/crawel-redis/redisClientWrapper"
	"fmt"
)

func main(){
	address := flag.String("address","localhost","redis address")
	port := flag.Int("port",6379,"port number of redis")
	client:= redisClientWrapper.ClientWrapperFactory(*address,*port)
	if (client.IsAlive()){
		fmt.Println("redis is working")
	} else {
		fmt.Println("redis is down")
	}
}