package work_redis

import "testing"

func Test_Redis(t *testing.T) {
	_, err := getRedisClient()
	if err != nil {
		return
	}
}
