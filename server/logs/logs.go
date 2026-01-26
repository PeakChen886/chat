// Package logs exposes info, warning and error loggers.
package logs

import (
	"io"
	"log"
	"strings"
)

var (
	// Info is a logger at the 'info' logging level.
	Info *LoggerBase
	// Warn is a logger at the 'warning' logging level.
	Warn *LoggerBase
	// Err is a logger at the 'error' logging level.
	Err *LoggerBase
)

func parseFlags(logFlags string) int {
	flags := 0
	for _, v := range strings.Split(logFlags, ",") {
		switch {
		case v == "date":
			flags |= log.Ldate
		case v == "time":
			flags |= log.Ltime
		case v == "microseconds":
			flags |= log.Lmicroseconds
		case v == "longfile":
			flags |= log.Llongfile
		case v == "shortfile":
			flags |= log.Lshortfile
		case v == "UTC":
			flags |= log.LUTC
		case v == "msgprefix":
			flags |= log.Lmsgprefix
		case v == "stdFlags":
			flags |= log.LstdFlags
		default:
			log.Fatalln("Invalid log flags string: ", logFlags)
		}
	}
	if flags == 0 {
		flags = log.LstdFlags
	}
	return flags
}

// Init initializes info, warning and error loggers given the flags and the output.
func Init(output io.Writer, logFlags string, level int) {
	flags := parseFlags(logFlags)
	Info = newLoggerBase(log.New(output, "I", flags), LevelInfo)
	Warn = newLoggerBase(log.New(output, "W", flags), LevelWarn)
	Err = newLoggerBase(log.New(output, "E", flags), LevelError)
	currentLevel = level
}

// 日志级别定义
const (
	LevelDebug = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

var currentLevel = LevelInfo // 默认日志级别

// SetLevel 设置日志级别
func SetLevel(level int) {
	currentLevel = level
}

// 检查是否应该输出指定级别的日志
func shouldLog(level int) bool {
	return level >= currentLevel
}
