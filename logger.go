package logger

import (
	"github.com/tkandal/ntnuzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

/*
 * Copyright (c) 2020 Norwegian University of Science and Technology
 */

// Logger logs JSON to stdout and stderr, and logs to a file if o filename is specified.
type Logger struct {
	logUTC      bool
	logFile     string
	maxSize     int
	maxBack     int
	maxDays     int
	development bool
	level       zapcore.Level
}

// LogOption is the type for options.
type LogOption func(*Logger) LogOption

// NewLogger allocates and returns a pointer to a logger-object, but the actual
// logger is not returned.  Use Sugar to return the actual logger.
func NewLogger(opts ...LogOption) *Logger {
	l := &Logger{
		logUTC:      true,
		maxSize:     128,
		maxBack:     3,
		maxDays:     28,
		development: false,
	}
	for _, opt := range opts {
		opt(l)
	}

	return l
}

// LogUTC set the logger to use UTC timestamps or local timestamp.
func LogUTC(b bool) LogOption {
	return func(l *Logger) LogOption {
		prev := l.logUTC
		l.logUTC = b
		return LogUTC(prev)
	}
}

// LogFile specify a file for logs in addition to stdout and stderr.
func LogFile(lf string) LogOption {
	return func(l *Logger) LogOption {
		prev := l.logFile
		l.logFile = lf
		return LogFile(prev)
	}
}

// MaxSize sets the maximum size for log-file in MB, before it is rotated.
func MaxSize(d int) LogOption {
	return func(l *Logger) LogOption {
		prev := l.maxSize
		l.maxSize = d
		return MaxSize(prev)
	}
}

// MaxBack sets the number of maximum number of generations for the log-file.
func MaxBack(d int) LogOption {
	return func(l *Logger) LogOption {
		prev := l.maxBack
		l.maxBack = d
		return MaxBack(prev)
	}
}

// MaxDays sets maximum age for the log-file before it is rotated.
func MaxDays(d int) LogOption {
	return func(l *Logger) LogOption {
		prev := l.maxDays
		l.maxDays = d
		return MaxDays(prev)
	}
}

// Development puts the logging in developer mode.
func Development(b bool) LogOption {
	return func(l *Logger) LogOption {
		prev := l.development
		l.development = b
		return Development(prev)
	}
}

func Level(level zapcore.Level) LogOption {
	return func(l *Logger) LogOption {
		prev := l.level
		l.level = level
		return Level(prev)
	}
}

// Sugar returns a sugared logger with log-lines formatted as JSON.
func (l *Logger) Sugar() (*zap.SugaredLogger, error) {
	if len(l.logFile) > 0 {
		logger, err := ntnuzap.NTNULumberjack(l.logFile, l.maxSize, l.maxBack, l.maxDays, l.logUTC)
		if err != nil {
			return nil, err
		}
		return logger.Sugar(), nil
	}
	files := []string{"stdout"}
	logger, err := ntnuzap.NTNUZap(files, l.level, l.development, l.logUTC)
	if err != nil {
		return nil, err
	}
	return logger.Sugar(), nil
}
