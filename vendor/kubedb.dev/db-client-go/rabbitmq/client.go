/*
Copyright AppsCode Inc. and Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package rabbitmq

import (
	rmqhttp "github.com/michaelklishin/rabbit-hole/v3"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Client struct {
	AMQPClient
	HTTPClient
	Channel
}

type AMQPClient struct {
	*amqp.Connection
}

type HTTPClient struct {
	*rmqhttp.Client
}

type Channel struct {
	*amqp.Channel
}

type ConnectionQueue struct {
	conn map[string]*Client
}

func NewConnectionQueue() *ConnectionQueue {
	return &ConnectionQueue{
		conn: make(map[string]*Client),
	}
}

func (c *ConnectionQueue) GetAMQPConnection(key string) *AMQPClient {
	return &c.conn[key].AMQPClient
}

func (c *ConnectionQueue) GetHTTPConnection(key string) *HTTPClient {
	return &c.conn[key].HTTPClient
}

func (c *ConnectionQueue) GetAMQPChannel(key string) *Channel {
	return &c.conn[key].Channel
}

func (c *ConnectionQueue) GetClientWithKey(key string) *Client {
	return c.conn[key]
}

func (c *ConnectionQueue) SetClientWithKey(key string, client *Client) {
	c.conn[key] = client
}
