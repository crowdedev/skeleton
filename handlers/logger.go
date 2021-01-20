package handlers

import (
	"fmt"

	configs "github.com/crowdeco/skeleton/configs"
	logrus "github.com/sirupsen/logrus"
	mongodb "github.com/weekface/mgorus"
)

type Logger struct {
	logger *logrus.Logger
}

func NewLogger() *Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	mongodb, err := mongodb.NewHooker(fmt.Sprintf("%s:%d", configs.Env.MongoDbHost, configs.Env.MongoDbPort), configs.Env.MongoDbName, "logs")
	if err == nil {
		logger.AddHook(mongodb)
	} else {
		fmt.Print(err)
	}

	return &Logger{
		logger: logger,
	}
}

func (l *Logger) Trace(message string) {
	if configs.Env.Debug {
		fields := logrus.Fields{
			"ServiceName": configs.Env.ServiceName,
			"Debug":       true,
		}

		l.logger.WithFields(fields).Trace(message)
	}
}

func (l *Logger) Debug(message string) {
	if configs.Env.Debug {
		fields := logrus.Fields{
			"ServiceName": configs.Env.ServiceName,
			"Debug":       true,
		}

		l.logger.WithFields(fields).Debug(message)
	}
}

func (l *Logger) Info(message string) {
	if configs.Env.Debug {
		fields := logrus.Fields{
			"ServiceName": configs.Env.ServiceName,
			"Debug":       true,
		}

		l.logger.WithFields(fields).Info(message)
	}
}

func (l *Logger) Warning(message string) {
	if configs.Env.Debug {
		fields := logrus.Fields{
			"ServiceName": configs.Env.ServiceName,
			"Debug":       true,
		}

		l.logger.WithFields(fields).Warning(message)
	}
}

func (l *Logger) Error(message string) {
	fields := logrus.Fields{
		"ServiceName": configs.Env.ServiceName,
		"Debug":       configs.Env.Debug,
	}

	l.logger.WithFields(fields).Error(message)
}

func (l *Logger) Fatal(message string) {
	fields := logrus.Fields{
		"ServiceName": configs.Env.ServiceName,
		"Debug":       configs.Env.Debug,
	}

	l.logger.WithFields(fields).Fatal(message)
}

func (l *Logger) Panic(message string) {
	fields := logrus.Fields{
		"ServiceName": configs.Env.ServiceName,
		"Debug":       configs.Env.Debug,
	}

	l.logger.WithFields(fields).Panic(message)
}
