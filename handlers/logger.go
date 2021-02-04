package handlers

import (
	"runtime"

	configs "github.com/crowdeco/skeleton/configs"
	logrus "github.com/sirupsen/logrus"
)

type Logger struct {
	Env    *configs.Env
	Logger *logrus.Logger
}

func (l *Logger) Trace(message string) {
	if l.Env.Debug {
		var file string
		var line int
		var caller string

		pc, file, line, ok := runtime.Caller(1)
		detail := runtime.FuncForPC(pc)
		if ok || detail != nil {
			caller = detail.Name()
		}

		fields := logrus.Fields{
			"ServiceName": l.Env.ServiceName,
			"Debug":       true,
			"Caller":      caller,
			"File":        file,
			"Line":        line,
		}

		l.Logger.WithFields(fields).Trace(message)
	}
}

func (l *Logger) Debug(message string) {
	if l.Env.Debug {
		var file string
		var line int
		var caller string

		pc, file, line, ok := runtime.Caller(1)
		detail := runtime.FuncForPC(pc)
		if ok || detail != nil {
			caller = detail.Name()
		}

		fields := logrus.Fields{
			"ServiceName": l.Env.ServiceName,
			"Debug":       true,
			"Caller":      caller,
			"File":        file,
			"Line":        line,
		}

		l.Logger.WithFields(fields).Debug(message)
	}
}

func (l *Logger) Info(message string) {
	if l.Env.Debug {
		var file string
		var line int
		var caller string

		pc, file, line, ok := runtime.Caller(1)
		detail := runtime.FuncForPC(pc)
		if ok || detail != nil {
			caller = detail.Name()
		}

		fields := logrus.Fields{
			"ServiceName": l.Env.ServiceName,
			"Debug":       true,
			"Caller":      caller,
			"File":        file,
			"Line":        line,
		}

		l.Logger.WithFields(fields).Info(message)
	}
}

func (l *Logger) Warning(message string) {
	if l.Env.Debug {
		var file string
		var line int
		var caller string

		pc, file, line, ok := runtime.Caller(1)
		detail := runtime.FuncForPC(pc)
		if ok || detail != nil {
			caller = detail.Name()
		}

		fields := logrus.Fields{
			"ServiceName": l.Env.ServiceName,
			"Debug":       true,
			"Caller":      caller,
			"File":        file,
			"Line":        line,
		}

		l.Logger.WithFields(fields).Warning(message)
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
		"ServiceName": l.Env.ServiceName,
		"Debug":       l.Env.Debug,
		"Caller":      caller,
		"File":        file,
		"Line":        line,
	}

	l.Logger.WithFields(fields).Error(message)
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
		"ServiceName": l.Env.ServiceName,
		"Debug":       l.Env.Debug,
		"Caller":      caller,
		"File":        file,
		"Line":        line,
	}

	l.Logger.WithFields(fields).Fatal(message)
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
		"ServiceName": l.Env.ServiceName,
		"Debug":       l.Env.Debug,
		"Caller":      caller,
		"File":        file,
		"Line":        line,
	}

	l.Logger.WithFields(fields).Panic(message)
}
