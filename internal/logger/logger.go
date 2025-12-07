package logger

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type PrefixHook struct {
	Prefix string
}

func (h *PrefixHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *PrefixHook) Fire(e *logrus.Entry) error {
	e.Message = h.Prefix + e.Message
	return nil
}

func SetupLogger(prefix string) *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	logger.SetLevel(logrus.InfoLevel)

	logger.AddHook(&PrefixHook{Prefix: fmt.Sprintf("[%s]", prefix)})

	logger.Info("Logger has been set up")
	return logger
}
