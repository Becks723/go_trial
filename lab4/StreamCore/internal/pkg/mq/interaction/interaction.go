package interaction

import (
	"StreamCore/internal/pkg/mq/model"
	"StreamCore/pkg/mq"
	"context"
	"errors"

	"github.com/bytedance/sonic"
	amqp "github.com/rabbitmq/amqp091-go"
)

type InteractionMQ interface {
	PublishLikeEvent(ctx context.Context, event *model.LikeEvent) error
	Consumer() (mq.Consumer, error)
}

func NewInteractionMQ(conn *amqp.Connection) InteractionMQ {
	queueName := "interaction_queue"
	sender, err := mq.NewRabbitSender(conn, queueName)
	if err != nil {
		sender = nil
		// TODO: log mq unavailable
	}
	return &iamq{
		sender:    sender,
		conn:      conn,
		queueName: queueName,
	}
}

var senderNotInitializedErr = errors.New("sender not initialized")

func (m *iamq) PublishLikeEvent(ctx context.Context, event *model.LikeEvent) error {
	if m.sender == nil {
		return senderNotInitializedErr
	}
	buffer, err := sonic.Marshal(event)
	if err != nil {
		return err
	}
	return m.sender.Send(ctx, buffer)
}

func (m *iamq) Consumer() (mq.Consumer, error) {
	return mq.NewRabbitConsumer(m.conn, m.queueName)
}

func (m *iamq) destroy() {
	if m.sender != nil {
		m.sender.Close()
	}
}

type iamq struct {
	sender    mq.Sender
	conn      *amqp.Connection
	queueName string
	exchange  string
}
