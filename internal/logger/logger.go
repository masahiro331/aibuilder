package logger

import (
    "github.com/sirupsen/logrus"
)

type Logger struct {
    *logrus.Logger
}

func NewLogger(debugMode bool) *Logger {
    logger := logrus.New()

    logger.SetFormatter(&logrus.TextFormatter{
        FullTimestamp: true,
    })

    if debugMode {
        logger.SetLevel(logrus.DebugLevel)
    } else {
        logger.SetLevel(logrus.InfoLevel)
    }

    return &Logger{logger}
}
