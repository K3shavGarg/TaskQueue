package pkg

import (
	"Job-Queue/internal/config"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Log *logrus.Logger

func init() {
	Log = logrus.New()

	// JSON output like Winston
	Log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05", // Go's date layout
		PrettyPrint:     false,
	})

	// Write to rotating file
	logFile := &lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    20, // MB
		MaxBackups: 0,  // Keep all backups
		MaxAge:     14, // days
		Compress:   true,
	}

	// Output to file
	Log.SetLevel(logrus.InfoLevel)

	Log.SetOutput(logFile)
	// Or combine with console:
	// Log.SetOutput(io.MultiWriter(os.Stdout, logFile))

	Log.WithFields(logrus.Fields{
		"service":     "queue-service",
		"environment": config.Env.Environment,
	}).Info("Logger initialized")
}
