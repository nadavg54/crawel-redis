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

func (c *ClientWrapper) AddToSet(set string,lines []string){
	_,err := c.internalClient.SAdd(set,lines).Result()
	if (err != nil){
		//add log here
		return
	}
}

//remove all the elements in set1 that intersect set2
func (c *ClientWrapper) RemoveInterSectionAndRetrieve(set1 string, set2 string, numberOfEleToRet int64) []string{
	//processed pending 
	pipeliner := c.internalClient.TxPipeline()
	pipeliner.SDiffStore(set1,set1,set2)
	exeResult := pipeliner.SRandMemberN(set1,numberOfEleToRet)
	_,err := pipeliner.Exec()
	
	if(err != nil){
		
	}
	elements,err := exeResult.Result()
	if(err != nil){

	}
	return elements
}

func ClientWrapperFactory(add string, port int) *ClientWrapper {
	var opt redis.Options
	opt.Addr = add + ":" + strconv.Itoa(port)
	newInternalClient := redis.NewClient(&opt)
	instance := ClientWrapper{internalClient:newInternalClient}
	return &instance
}