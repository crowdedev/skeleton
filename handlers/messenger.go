package handlers

import (
	"context"
	"fmt"
	"time"

	watermill "github.com/ThreeDotsLabs/watermill"
	amqp "github.com/ThreeDotsLabs/watermill-amqp/pkg/amqp"
	message "github.com/ThreeDotsLabs/watermill/message"
	configs "github.com/crowdeco/skeleton/configs"
)

type Messenger struct {
	publisher *amqp.Publisher
	consumer  *amqp.Subscriber
	Logger    *Logger
}

func NewMessenger(logger *Logger) *Messenger {
	address := fmt.Sprintf("amqp://%s:%s@%s:%d/", configs.Env.AmqpUser, configs.Env.AmqpPassword, configs.Env.AmqpHost, configs.Env.AmqpPort)
	config := amqp.NewDurableQueueConfig(address)

	publisher, err := amqp.NewPublisher(config, watermill.NewStdLogger(configs.Env.Debug, configs.Env.Debug))
	if err != nil {
		panic(err)
	}

	consumer, err := amqp.NewSubscriber(config, watermill.NewStdLogger(false, false))
	if err != nil {
		panic(err)
	}

	return &Messenger{
		Logger:    logger,
		publisher: publisher,
		consumer:  consumer,
	}
}

func (m *Messenger) Publish(queueName string, data []byte) error {
	for {
		msg := message.NewMessage(watermill.NewUUID(), data)
		err := m.publisher.Publish(queueName, msg)
		if err != nil {
			m.Logger.Error(err.Error())

			return err
		}

		time.Sleep(time.Second)
	}
}

func (m *Messenger) Consume(queueName string) (<-chan *message.Message, error) {
	messages, err := m.consumer.Subscribe(context.Background(), queueName)
	if err != nil {
		m.Logger.Error(err.Error())

		return nil, err
	}

	return messages, nil
}
