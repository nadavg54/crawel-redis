package redisClientWrapper

import (
	"strconv"
	"github.com/go-redis/redis"
	"log"
)

type ClientWrapper struct {
	
	internalClient *redis.Client
}

func (c *ClientWrapper) IsAlive() bool{
	_,err := c.internalClient.Ping().Result()
	if (err != nil){
		log.Println("can reach redis with error: " + err.Error())
		return false
	}
	return true
}

func ClientWrapperFactory(add string, port int) *ClientWrapper {
	var opt redis.Options
	opt.Addr = add + ":" + strconv.Itoa(port)
	newInternalClient := redis.NewClient(&opt)
	instance := ClientWrapper{internalClient:newInternalClient}
	return &instance
}