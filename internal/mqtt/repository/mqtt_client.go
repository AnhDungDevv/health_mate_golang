package repositoryMQTT

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTClientImpl struct {
	client mqtt.Client
}

func NewMQTTClient(brokerURL, clientID string) *MQTTClientImpl {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(brokerURL)
	opts.SetClientID(clientID)
	opts.SetKeepAlive(60 * time.Second)
	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	})
	client := mqtt.NewClient(opts)
	return &MQTTClientImpl{client: client}
}

func (m *MQTTClientImpl) Connect() error {
	if token := m.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	fmt.Println("Connected to MQTT broker")
	return nil
}

func (m *MQTTClientImpl) Publish(topic string, payload interface{}) error {
	token := m.client.Publish(topic, 0, false, payload)
	token.Wait()
	return token.Error()
}

func (m *MQTTClientImpl) Subcribe(topic string, callback mqtt.MessageHandler) error {
	token := m.client.Subscribe(topic, 0, callback)
	token.Wait()
	return token.Error()
}

func (m *MQTTClientImpl) Disconnect() {
	m.client.Disconnect(250)
}
