package work_rabbitMQ

import "testing"

func Test_GetRabbitMQClient(t *testing.T) {
	err, _ := getRabbitMQClient()
	if err != nil {
		return
	}
}
