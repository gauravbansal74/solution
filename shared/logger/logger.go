package logger

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gauravbansal74/solution/config"
	"os"
)

var (
	AppName = ""
)

type Fields map[string]interface{}

// Logger Config using Env Vars.
func LoggerConfig(conf config.Config) {
	// Set Log Formatter to JSON formatter
	log.SetFormatter(&log.JSONFormatter{})

	//Get and Configure the name on Application Level
	AppName = conf.SystemLog

	// Set Logger Output
	log.SetOutput(os.Stdout)

	// Set Log Level -  Debugging is enable/disable
	if conf.Debug {
		log.SetLevel(log.DebugLevel) // Debug Level +
	} else {
		log.SetLevel(log.InfoLevel) // Warn Level +
	}
}

// ("payments", "Msg", {"key":value})
func Debug(message string, other ...map[string]interface{}) {
	fields := map[string]interface{}{}
	if other != nil {
		fields = other[0]
	}
	log.WithFields(log.Fields{"name": AppName}).WithFields(fields).Debug(message)
}

// ("payments", "Msg", {"key":value})
func Info(who string, message string, other ...map[string]interface{}) {
	fields := map[string]interface{}{}
	if other != nil {
		fields = other[0]
	}
	log.WithFields(log.Fields{"user": who, "name": AppName}).WithFields(fields).Info(message)
}

// ("payments", error, "err", {"key":value})
func Fatal(who string, err error, message string, other ...map[string]interface{}) {
	fields := map[string]interface{}{}
	if other != nil {
		fields = other[0]
	}
	log.WithFields(log.Fields{"user": who, "name": AppName}).WithFields(fields).WithError(err).Fatal(message)
}

// ("payments", error, "err", {"key":value})
func Error(who string, err error, message string, other ...map[string]interface{}) {
	fields := map[string]interface{}{}
	if other != nil {
		fields = other[0]
	}
	log.WithFields(log.Fields{"user": who, "name": AppName}).WithFields(fields).WithError(err).Error(message)
}
