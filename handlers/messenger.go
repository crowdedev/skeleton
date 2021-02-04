package handlers

import (
	"context"
	"time"

	watermill "github.com/ThreeDotsLabs/watermill"
	amqp "github.com/ThreeDotsLabs/watermill-amqp/pkg/amqp"
	message "github.com/ThreeDotsLabs/watermill/message"
)

type Messenger struct {
	Publisher *amqp.Publisher
	Consumer  *amqp.Subscriber
	Logger    *Logger
}

func (m *Messenger) Publish(queueName string, data []byte) error {
	for {
		msg := message.NewMessage(watermill.NewUUID(), data)
		err := m.Publisher.Publish(queueName, msg)
		if err != nil {
			m.Logger.Error(err.Error())

			return err
		}

		time.Sleep(time.Second)
	}
}

func (m *Messenger) Consume(queueName string) (<-chan *message.Message, error) {
	messages, err := m.Consumer.Subscribe(context.Background(), queueName)
	if err != nil {
		m.Logger.Error(err.Error())

		return nil, err
	}

	return messages, nil
}
