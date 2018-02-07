package logger

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	logrustash "github.com/bshuster-repo/logrus-logstash-hook"
	"github.com/ld100/goblet/pkg/util/config"
	"github.com/sirupsen/logrus"
)

// Logging wrapper on top of sirupsen/logrus
// Wraps https://godoc.org/github.com/sirupsen/logrus
type Logger struct {
	*logrus.Logger
}

func (lg *Logger) SetLogLevel(level logrus.Level) {
	lg.Logger.Level = level
}

func (lg *Logger) SetLogFormatter(formatter logrus.Formatter) {
	lg.Logger.Formatter = formatter
}

func (lg *Logger) SetOutput(out io.Writer) {
	logrus.SetOutput(out)
}

// Debug logs a message at level Debug on the standard logger.
func (lg *Logger) Debug(args ...interface{}) {
	if lg.Logger.Level >= logrus.DebugLevel {
		entry := lg.Logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(2)
		entry.Debug(args)
	}
}

// Debug logs a message with fields at level Debug on the standard logger.
func (lg *Logger) DebugWithFields(l interface{}, f Fields) {
	if lg.Logger.Level >= logrus.DebugLevel {
		entry := lg.Logger.WithFields(logrus.Fields(f))
		entry.Data["file"] = fileInfo(2)
		entry.Debug(l)
	}
}

// Info logs a message at level Info on the standard logger.
func (lg *Logger) Info(args ...interface{}) {
	if lg.Logger.Level >= logrus.InfoLevel {
		entry := lg.Logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(2)
		entry.Info(args...)
	}
}

// Debug logs a message with fields at level Debug on the standard logger.
func (lg *Logger) InfoWithFields(l interface{}, f Fields) {
	if lg.Logger.Level >= logrus.InfoLevel {
		entry := lg.Logger.WithFields(logrus.Fields(f))
		entry.Data["file"] = fileInfo(2)
		entry.Info(l)
	}
}

// Warn logs a message at level Warn on the standard logger.
func (lg *Logger) Warn(args ...interface{}) {
	if lg.Logger.Level >= logrus.WarnLevel {
		entry := lg.Logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(2)
		entry.Warn(args...)
	}
}

// Debug logs a message with fields at level Debug on the standard logger.
func (lg *Logger) WarnWithFields(l interface{}, f Fields) {
	if lg.Logger.Level >= logrus.WarnLevel {
		entry := lg.Logger.WithFields(logrus.Fields(f))
		entry.Data["file"] = fileInfo(2)
		entry.Warn(l)
	}
}

// Error logs a message at level Error on the standard logger.
func (lg *Logger) Error(args ...interface{}) {
	if lg.Logger.Level >= logrus.ErrorLevel {
		entry := lg.Logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(2)
		entry.Error(args...)
	}
}

// Debug logs a message with fields at level Debug on the standard logger.
func (lg *Logger) ErrorWithFields(l interface{}, f Fields) {
	if lg.Logger.Level >= logrus.ErrorLevel {
		entry := lg.Logger.WithFields(logrus.Fields(f))
		entry.Data["file"] = fileInfo(2)
		entry.Error(l)
	}
}

// Fatal logs a message at level Fatal on the standard logger.
func (lg *Logger) Fatal(args ...interface{}) {
	if lg.Logger.Level >= logrus.FatalLevel {
		entry := lg.Logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(2)
		entry.Fatal(args...)
	}
}

// Debug logs a message with fields at level Debug on the standard logger.
func (lg *Logger) FatalWithFields(l interface{}, f Fields) {
	if lg.Logger.Level >= logrus.FatalLevel {
		entry := lg.Logger.WithFields(logrus.Fields(f))
		entry.Data["file"] = fileInfo(2)
		entry.Fatal(l)
	}
}

// Panic logs a message at level Panic on the standard logger.
func (lg *Logger) Panic(args ...interface{}) {
	if lg.Logger.Level >= logrus.PanicLevel {
		entry := lg.Logger.WithFields(logrus.Fields{})
		entry.Data["file"] = fileInfo(2)
		entry.Panic(args...)
	}
}

// Debug logs a message with fields at level Debug on the standard logger.
func (lg *Logger) PanicWithFields(l interface{}, f Fields) {
	if lg.Logger.Level >= logrus.PanicLevel {
		entry := lg.Logger.WithFields(logrus.Fields(f))
		entry.Data["file"] = fileInfo(2)
		entry.Panic(l)
	}
}

// Exit runs all the Logrus atexit handlers and then terminates the program using os.Exit(code)
func (lg *Logger) Exit(code int) {
	logrus.Exit(code)
}

func New(cfg *config.Config) (*Logger) {
	ls := logrus.New()
	logger := &Logger{ls}

	// Log as JSON instead of the default ASCII formatter.
	logger.SetLogFormatter(&logrus.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logger.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	// TODO: Take debug level from cfg
	logger.SetLogLevel(logrus.DebugLevel)

	//Enabling logstash hook
	logstashEnabled := cfg.GetBool("LOGSTASH_ENABLED")
	if logstashEnabled {
		logstashPort := cfg.GetInt("LOGSTASH_PORT")
		logstashHost := cfg.GetString("LOGSTASH_HOST")
		logstashUrl := fmt.Sprintf("%v:%v", logstashHost, logstashPort)
		appName := cfg.GetString("APP_NAME")

		hook, err := logrustash.NewHook("tcp", logstashUrl, appName)

		if err != nil {
			logger.Fatal(err)
		}
		logger.Logger.Hooks.Add(hook)
	}

	return logger
}

// Fields wraps logrus.Fields, which is a map[string]interface{}
type Fields logrus.Fields

func fileInfo(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, "/")
		if slash >= 0 {
			file = file[slash+1:]
		}
	}
	return fmt.Sprintf("%s:%d", file, line)
}
