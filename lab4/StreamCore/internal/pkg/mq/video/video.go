package video

import (
	"StreamCore/internal/pkg/mq/model"
	"StreamCore/pkg/mq"
	"context"
	"errors"

	"github.com/bytedance/sonic"
	amqp "github.com/rabbitmq/amqp091-go"
)

type VideoMQ interface {
	PublishVisitEvent(ctx context.Context, event *model.VisitEvent) error
	Consumer() (mq.Consumer, error)
}

func NewVideoMQ(conn *amqp.Connection) VideoMQ {
	queueName := "video_queue"
	sender, err := mq.NewRabbitSender(conn, queueName)
	if err != nil {
		sender = nil
		// TODO: log mq unavailable
	}
	return &videomq{
		sender:    sender,
		conn:      conn,
		queueName: queueName,
	}
}

var senderNotInitializedErr = errors.New("sender not initialized")

func (m *videomq) PublishVisitEvent(ctx context.Context, event *model.VisitEvent) error {
	if m.sender == nil {
		return senderNotInitializedErr
	}
	buffer, err := sonic.Marshal(event)
	if err != nil {
		return err
	}
	return m.sender.Send(ctx, buffer)
}

func (m *videomq) Consumer() (mq.Consumer, error) {
	return mq.NewRabbitConsumer(m.conn, m.queueName)
}

func (m *videomq) destroy() {
	if m.sender != nil {
		m.sender.Close()
	}
}

type videomq struct {
	sender    mq.Sender
	conn      *amqp.Connection
	queueName string
	exchange  string
}
