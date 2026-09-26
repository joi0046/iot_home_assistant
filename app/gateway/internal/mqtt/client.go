package mqtt

import (
	"fmt"

	paho "github.com/eclipse/paho.mqtt.golang"
)

type Client struct {
	client paho.Client
}

func New(brokerURL string) (*Client, error) {
	opts := paho.NewClientOptions()
	opts.AddBroker(brokerURL)

	opts.AutoReconnect = true
	opts.ConnectRetry = true

	client := paho.NewClient(opts)

	token := client.Connect()
	token.Wait()

	if err := token.Error(); err != nil {
		return nil, fmt.Errorf("MQTT connect: %w", err)
	}

	return &Client{
		client: client,
	}, nil
}

func (c *Client) Publish(topic string, payload string) error {
	token := c.client.Publish(topic, 0, false, payload)
	token.Wait()

	if err := token.Error(); err != nil {
		return fmt.Errorf("MQTT publish: %w", err)
	}

	return nil
}

func (c *Client) IsConnected() bool {
	return c.client.IsConnected()
}

func (c *Client) Close() {
	c.client.Disconnect(250)
}
