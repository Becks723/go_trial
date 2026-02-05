package infra

import (
	"StreamCore/config"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func InitRabbitMQ() (*amqp.Connection, error) {
	conf := config.Instance().RabbitMQ
	url := fmt.Sprintf("amqp://%s:%s@localhost:5672/", conf.Username, conf.Password)
	return amqp.Dial(url)
}
