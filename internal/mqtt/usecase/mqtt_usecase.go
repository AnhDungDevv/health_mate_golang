package usecase

import (
	"encoding/json"
	"fmt"
	"health_backend/internal/mqtt/interfaces"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTUsecase struct {
	mqttClient interfaces.MQTTClient
}

func NewMQTTUsecase(mqttClient interfaces.MQTTClient) *MQTTUsecase {
	return &MQTTUsecase{mqttClient: mqttClient}
}

func (m *MQTTUsecase) SendMessage(topic string, message interface{}) error {
	msg, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return m.mqttClient.Publish(topic, msg)
}

func (m *MQTTUsecase) ReceiveMessage(topic string) error {
	return m.mqttClient.Subcribe(topic, func(c mqtt.Client, msg mqtt.Message) {
		fmt.Printf("Received message on %s: %s\n", msg.Topic(), msg.Payload())

	})
}
