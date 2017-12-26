package gocd

import (
	log "github.com/Sirupsen/logrus"
	"os"
)

// Set logging level and type constants
const (
	LogLevelEnvVarName = "GOCD_LOG_LEVEL"
	LogLevelDefault    = "WARNING"
	LogTypeEnvVarName  = "GOCD_LOG_TYPE"
	LogTypeDefault     = "TEXT"
)

var logLevels = map[string]log.Level{
	"PANIC":   log.PanicLevel,
	"FATAL":   log.FatalLevel,
	"ERROR":   log.ErrorLevel,
	"WARNING": log.WarnLevel,
	"INFO":    log.InfoLevel,
	"DEBUG":   log.DebugLevel,
}

var logFormat = map[string]log.Formatter{
	"JSON": &log.JSONFormatter{},
	"TEXT": &log.TextFormatter{},
}

// SetupLogging based on Environment Variables
//
//  Set Logging level with $GOCD_LOG_LEVEL
//  Allowed Values:
//    - DEBUG
//    - INFO
//    - WARNING
//    - ERROR
//    - FATAL
//    - PANIC
//
//  Set Logging type  with $GOCD_LOG_TYPE
//  Allowed Values:
//    - JSON
//    - TEXT
func SetupLogging() {
	log.SetLevel(logLevels[getLogLevel()])

	log.SetFormatter(logFormat[getLogType()])
}

// Get the log type from env variables
func getLogType() string {
	logType := os.Getenv(LogTypeEnvVarName)
	if len(logType) == 0 {
		// If no env is set, return the default
		return LogTypeDefault
	}
	return logType
}

// Get the log level from env variables
func getLogLevel() string {
	loglevel := os.Getenv(LogLevelEnvVarName)
	if len(loglevel) == 0 {
		// If no env is set, return the default
		return LogLevelDefault
	}
	return loglevel
}
