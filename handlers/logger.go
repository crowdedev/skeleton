package handlers

import (
	"fmt"
	"runtime"

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
		var file string
		var line int
		var caller string

		pc, file, line, ok := runtime.Caller(1)
		detail := runtime.FuncForPC(pc)
		if ok || detail != nil {
			caller = detail.Name()
		}

		fields := logrus.Fields{
			"ServiceName": configs.Env.ServiceName,
			"Debug":       true,
			"Caller":      caller,
			"File":        file,
			"Line":        line,
		}

		l.logger.WithFields(fields).Trace(message)
	}
}

func (l *Logger) Debug(message string) {
	if configs.Env.Debug {
		var file string
		var line int
		var caller string

		pc, file, line, ok := runtime.Caller(1)
		detail := runtime.FuncForPC(pc)
		if ok || detail != nil {
			caller = detail.Name()
		}

		fields := logrus.Fields{
			"ServiceName": configs.Env.ServiceName,
			"Debug":       true,
			"Caller":      caller,
			"File":        file,
			"Line":        line,
		}

		l.logger.WithFields(fields).Debug(message)
	}
}

func (l *Logger) Info(message string) {
	if configs.Env.Debug {
		var file string
		var line int
		var caller string

		pc, file, line, ok := runtime.Caller(1)
		detail := runtime.FuncForPC(pc)
		if ok || detail != nil {
			caller = detail.Name()
		}

		fields := logrus.Fields{
			"ServiceName": configs.Env.ServiceName,
			"Debug":       true,
			"Caller":      caller,
			"File":        file,
			"Line":        line,
		}

		l.logger.WithFields(fields).Info(message)
	}
}

func (l *Logger) Warning(message string) {
	if configs.Env.Debug {
		var file string
		var line int
		var caller string

		pc, file, line, ok := runtime.Caller(1)
		detail := runtime.FuncForPC(pc)
		if ok || detail != nil {
			caller = detail.Name()
		}

		fields := logrus.Fields{
			"ServiceName": configs.Env.ServiceName,
			"Debug":       true,
			"Caller":      caller,
			"File":        file,
			"Line":        line,
		}

		l.logger.WithFields(fields).Warning(message)
	}
}

func (l *Logger) Error(message string) {
	var file string
	var line int
	var caller string

	pc, file, line, ok := runtime.Caller(1)
	detail := runtime.FuncForPC(pc)
	if ok || detail != nil {
		caller = detail.Name()
	}

	fields := logrus.Fields{
		"ServiceName": configs.Env.ServiceName,
		"Debug":       configs.Env.Debug,
		"Caller":      caller,
		"File":        file,
		"Line":        line,
	}

	l.logger.WithFields(fields).Error(message)
}

func (l *Logger) Fatal(message string) {
	var file string
	var line int
	var caller string

	pc, file, line, ok := runtime.Caller(1)
	detail := runtime.FuncForPC(pc)
	if ok || detail != nil {
		caller = detail.Name()
	}

	fields := logrus.Fields{
		"ServiceName": configs.Env.ServiceName,
		"Debug":       configs.Env.Debug,
		"Caller":      caller,
		"File":        file,
		"Line":        line,
	}

	l.logger.WithFields(fields).Fatal(message)
}

func (l *Logger) Panic(message string) {
	var file string
	var line int
	var caller string

	pc, file, line, ok := runtime.Caller(1)
	detail := runtime.FuncForPC(pc)
	if ok || detail != nil {
		caller = detail.Name()
	}

	fields := logrus.Fields{
		"ServiceName": configs.Env.ServiceName,
		"Debug":       configs.Env.Debug,
		"Caller":      caller,
		"File":        file,
		"Line":        line,
	}

	l.logger.WithFields(fields).Panic(message)
}
