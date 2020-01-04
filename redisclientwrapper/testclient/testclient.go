package main

import (
	"github.com/nadavg54/crawel-redis/redisclientwrapper"
)

func main() {
	client := redisclientwrapper.ClientWrapperFactory("localhost", 6379)
	_ = client.AddToMultipleSets("set11", []string{"val1"}, "set22", []string{"val2"})
}
