package testsupport

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"testing"
	"time"
)

type MQTTClientOptions struct {
	Port int

	Username string
	Password string
}

type MQTTClient interface {
	Publish(topic string, payload string)
}

type pahoMqttClient struct {
	client mqtt.Client
	t      *testing.T
}

func (c *pahoMqttClient) Publish(topic string, payload string) {
	token := c.client.Publish(topic, 0, false, payload)

	if token.Wait() && token.Error() != nil {
		c.t.Fatalf("Failed to publish message: %s", token.Error())
	}
}

func CreateMQTTClient(t *testing.T, options MQTTClientOptions) MQTTClient {
	errC := make(chan error)

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("127.0.0.1:%d", options.Port))
	opts.SetClientID("acceptance-test")

	if options.Username != "" {
		opts.SetUsername(options.Username)
		opts.SetPassword(options.Password)
	}

	opts.SetCleanSession(true)
	opts.SetAutoReconnect(false) // The broker should not go offline during tests, if that does happen, it's indicative of a problem that needs solving.
	opts.SetConnectTimeout(time.Millisecond * 100)
	opts.SetOnConnectHandler(func(client mqtt.Client) {
		errC <- nil
	})

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		t.Fatalf("Failed to connect to MQTT broker: %s", token.Error())
	}

	err := <-errC

	if err != nil {
		t.Fatalf("Failed to connect to MQTT broker: %s", err)
	}

	return &pahoMqttClient{
		client: client,
	}
}
