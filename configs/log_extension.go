package configs

import (
	"github.com/sirupsen/logrus"
)

type LoggerExtension struct {
	Extensions []logrus.Hook
}

func (l *LoggerExtension) Register(extensions []logrus.Hook) {
	l.Extensions = extensions
}
