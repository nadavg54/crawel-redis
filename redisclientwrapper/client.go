package redisclientwrapper

import (
	"errors"
	"log"
	"strconv"

	"github.com/go-redis/redis"
)

//ClientWrapper basic redis client wrapper
type ClientWrapper struct {
	internalClient *redis.Client
}

//IsAlive checking the connection with redis
func (c *ClientWrapper) IsAlive() bool {
	_, err := c.internalClient.Ping().Result()
	if err != nil {
		log.Println("can reach redis with error: " + err.Error())
		return false
	}
	return true
}

//AddToSet adding array line to set
func (c *ClientWrapper) AddToSet(set string, lines []string) error {
	_, err := c.internalClient.SAdd(set, lines).Result()
	if err != nil {
		return errors.New("cant add to set " + set + " with error " + err.Error())
		//add log here
	}
	return nil
}

//RemoveInterSectionAndRetrieve remove all the elements in set1 that intersect set2
func (c *ClientWrapper) RemoveInterSectionAndRetrieve(set1 string, set2 string, numberOfEleToRet int64) ([]string, error) {
	//processed pending
	pipeliner := c.internalClient.TxPipeline()
	pipeliner.SDiffStore(set1, set1, set2)
	exeResult := pipeliner.SPopN(set1, numberOfEleToRet)
	_, err := pipeliner.Exec()

	if err != nil {
		return nil, err
	}
	elements, err := exeResult.Result()
	if err != nil {
		return nil, err
	}
	return elements, nil
}

//AddToMultipleSets transaction adding values to two sets
func (c *ClientWrapper) AddToMultipleSets(set1 string, val1 []string, set2 string, val2 []string) error {
	pipeliner := c.internalClient.TxPipeline()
	pipeliner.SAdd(set1, val1)
	pipeliner.SAdd(set2, val2)
	_, err := pipeliner.Exec()
	
	if err != nil {
		return  err
	}
	return nil
}

//ClientWrapperFactory builds new redis wrapper instance
func ClientWrapperFactory(add string, port int) *ClientWrapper {
	var opt redis.Options
	opt.Addr = add + ":" + strconv.Itoa(port)
	newInternalClient := redis.NewClient(&opt)
	instance := ClientWrapper{internalClient: newInternalClient}
	return &instance
}

// func main() {
// 	client := ClientWrapperFactory("localhost", 6379)
// 	res, _ := client.RemoveInterSectionAndRetrieve("set1", "set2", 300)
// 	fmt.Println(res)
// }
