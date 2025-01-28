package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func InitLogger() {
	Log = logrus.New()

	// Set the output to stdout
	Log.Out = os.Stdout

	// Set the log level (info by default)
	Log.SetLevel(logrus.InfoLevel)

	// Use JSON formatter for structured logging
	Log.SetFormatter(&logrus.JSONFormatter{})
}
