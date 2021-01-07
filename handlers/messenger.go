package handlers

import (
	"errors"
	"fmt"
	"time"

	configs "github.com/crowdeco/todo-service/configs"
	amqp "github.com/streadway/amqp"
)

type Messenger struct {
	queueName       string
	logger          *Logger
	connection      *amqp.Connection
	channel         *amqp.Channel
	done            chan bool
	notifyConnClose chan *amqp.Error
	notifyChanClose chan *amqp.Error
	notifyConfirm   chan amqp.Confirmation
	isReady         bool
}

func NewMessenger(queueName string) *Messenger {
	m := Messenger{
		logger:    NewLogger(),
		queueName: queueName,
		done:      make(chan bool),
	}

	address := fmt.Sprintf("amqp://%s:%s@%s:%d/", configs.Env.AmqpUser, configs.Env.AmqpPassword, configs.Env.AmqpHost, configs.Env.AmqpPort)

	go m.handleReconnect(address)

	return &m
}

func (m *Messenger) handleReconnect(address string) {
	for {
		fmt.Println("Attempting to connect")

		m.isReady = false
		connection, err := m.connect(address)

		if err != nil {
			fmt.Println("Failed to connect. Retrying...")

			select {
			case <-m.done:
				return
			case <-time.After(5 * time.Second):
			}

			continue
		}

		if done := m.handleReInit(connection); done {
			break
		}
	}
}

func (m *Messenger) connect(address string) (*amqp.Connection, error) {
	connection, err := amqp.Dial(address)
	if err != nil {
		return nil, err
	}

	m.changeConnection(connection)

	fmt.Println("Connected!")

	return connection, nil
}

func (m *Messenger) handleReInit(connection *amqp.Connection) bool {
	for {
		m.isReady = false

		err := m.init(connection)
		if err != nil {
			fmt.Println("Failed to initialize channel. Retrying...")

			select {
			case <-m.done:
				return true
			case <-time.After(2 * time.Second):
			}
			continue
		}

		select {
		case <-m.done:
			return true
		case <-m.notifyConnClose:
			fmt.Println("Connection closed. Reconnecting...")
			return false
		case <-m.notifyChanClose:
			fmt.Println("Channel closed. Re-running init...")
		}
	}
}

func (m *Messenger) init(connection *amqp.Connection) error {
	channel, err := connection.Channel()
	if err != nil {
		return err
	}

	err = channel.Confirm(false)
	if err != nil {
		return err
	}

	_, err = channel.QueueDeclare(
		m.queueName,
		false, // Durable
		false, // Delete when unused
		false, // Exclusive
		false, // No-wait
		nil,   // Arguments
	)

	if err != nil {
		return err
	}

	m.changeChannel(channel)
	m.isReady = true

	fmt.Println("Setup!")

	return nil
}

func (m *Messenger) changeConnection(connection *amqp.Connection) {
	m.connection = connection
	m.notifyConnClose = make(chan *amqp.Error)
	m.connection.NotifyClose(m.notifyConnClose)
}

func (m *Messenger) changeChannel(channel *amqp.Channel) {
	m.channel = channel
	m.notifyChanClose = make(chan *amqp.Error)
	m.notifyConfirm = make(chan amqp.Confirmation, 1)
	m.channel.NotifyClose(m.notifyChanClose)
	m.channel.NotifyPublish(m.notifyConfirm)
}

func (m *Messenger) Push(data []byte) error {
	if !m.isReady {
		return errors.New("failed to push: not connected")
	}

	for {
		err := m.UnsafePush(data)
		if err != nil {
			m.logger.Error("Push failed. Retrying...")
			select {
			case <-m.done:
				return errors.New("Session is shutting down")
			case <-time.After(5 * time.Second):
			}
			continue
		}

		select {
		case confirm := <-m.notifyConfirm:
			if confirm.Ack {
				m.logger.Info("Push confirmed!")
				return nil
			}
		case <-time.After(5 * time.Second):
		}

		m.logger.Error("Push didn't confirm. Retrying...")
	}
}

func (m *Messenger) UnsafePush(data []byte) error {
	if !m.isReady {
		return errors.New("not connected to a server")
	}

	m.logger.Info(fmt.Sprintf("Send to queue: %+v", json))

	return m.channel.Publish(
		"",          // Exchange
		m.queueName, // Routing key
		false,       // Mandatory
		false,       // Immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        data,
		},
	)
}

func (m *Messenger) Consume() (<-chan amqp.Delivery, error) {
	if !m.isReady {
		return nil, errors.New("not connected to a server")
	}

	return m.channel.Consume(
		m.queueName,
		"",    // Consumer
		false, // Auto-Ack
		false, // Exclusive
		false, // No-local
		false, // No-Wait
		nil,   // Args
	)
}

func (m *Messenger) Close() error {
	if !m.isReady {
		return errors.New("already closed: not connected to the server")
	}

	err := m.channel.Close()
	if err != nil {
		return err
	}

	err = m.connection.Close()
	if err != nil {
		return err
	}

	close(m.done)

	m.isReady = false

	return nil
}
