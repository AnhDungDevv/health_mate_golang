package interfaces

import mttq "github.com/eclipse/paho.mqtt.golang"

type MQTTClient interface {
	Connect() error
	Publish(topic string, payload interface{}) error
	Subcribe(topic string, callback mttq.MessageHandler) error
	Disconnect()
}
