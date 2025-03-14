package work_redis

import "testing"

func Test_Redis(t *testing.T) {
	client, _ := getRedisClient()
	//if err != nil {
	//	return
	//}
	checkNetwork(client)
}
