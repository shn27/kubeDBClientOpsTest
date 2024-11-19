package work

import (
	"fmt"
	"time"
)

func Ans() {
	fmt.Println("=================")
	mongodbClient, err := GetMongoDBClientUsingAppRef()
	if err != nil {
		fmt.Printf("get db client error: %v\n", err)
		return
	}
	for {
		shard(mongodbClient)
		time.Sleep(10 * time.Second)
	}
}
