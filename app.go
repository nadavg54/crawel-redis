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
	client.AddToSet("pending",[]string{"a","b"})
	client.AddToSet("processed",[]string{"a"})
	result := client.RemoveInterSectionAndRetrieve("pending","processed",1)
	fmt.Println(result)
}