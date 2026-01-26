package logs

import "log"

type LoggerBase struct {
	logger *log.Logger
	Level  int
}

func newLoggerBase(logger *log.Logger, level int) *LoggerBase {
	return &LoggerBase{logger: logger, Level: level}
}

func (l *LoggerBase) Println(v ...interface{}) {
	if shouldLog(l.Level) {
		l.logger.Println(v...)
	}
}

func (l *LoggerBase) Printf(format string, v ...interface{}) {
	if shouldLog(l.Level) {
		l.logger.Printf(format, v...)
	}
}

func (l *LoggerBase) Fatal(v ...interface{}) {
	if shouldLog(l.Level) {
		l.logger.Fatal(v...)
	}
}

func (l *LoggerBase) Fatalf(format string, v ...interface{}) {
	if shouldLog(l.Level) {
		l.logger.Fatalf(format, v...)
	}
}

func (l *LoggerBase) Panic(v ...interface{}) {
	if shouldLog(l.Level) {
		l.logger.Panic(v...)
	}
}

func (l *LoggerBase) Panicf(format string, v ...interface{}) {
	if shouldLog(l.Level) {
		l.logger.Panicf(format, v...)
	}
}

func (l *LoggerBase) Fatalln(v ...interface{}) {
	if shouldLog(l.Level) {
		l.logger.Fatalln(v...)
	}
}

func (l *LoggerBase) Panicln(v ...interface{}) {
	if shouldLog(l.Level) {
		l.logger.Panicln(v...)
	}
}
