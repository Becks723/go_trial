package mq

import (
	"StreamCore/internal/pkg/mq/interaction"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MQSet struct {
	Interaction interaction.InteractionMQ
}

func NewMQSet(conn *amqp.Connection) *MQSet {
	return &MQSet{
		Interaction: interaction.NewInteractionMQ(conn),
	}
}
