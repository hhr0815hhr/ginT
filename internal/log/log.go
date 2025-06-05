package log

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"sync"
	"time"
)

var Logger *logrus.Logger

func init() {
	initSysLog()
}

// MyLog is a custom logging function that includes filename, level, message, and data.
func MyLog(filename string) *logrus.Logger {
	logger, _ := GetLoggerForFileWithRotation(filename)
	return logger
}

func initSysLog() {
	Logger = logrus.New()
	Logger.SetOutput(os.Stdout)
	Logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: time.DateTime,
	})
	Logger.SetLevel(logrus.DebugLevel)
}

var (
	loggers   = make(map[string]*logrus.Logger)
	loggersMu sync.Mutex
)

// GetLoggerForFileWithRotation returns a Logrus logger instance for the given filename with lumberjack for rotation.
// It creates the logger and configures lumberjack for the log file.
func GetLoggerForFileWithRotation(filename string) (*logrus.Logger, error) {
	loggersMu.Lock()
	defer loggersMu.Unlock()

	if logger, ok := loggers[filename]; ok {
		return logger, nil
	}

	newLogger := logrus.New()
	newLogger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.DateTime,
	}) // 使用 JSON 格式

	// 配置 lumberjack 进行日志轮转
	lumberjackLogger := &lumberjack.Logger{
		Filename:   "./logs/" + filename + ".log", // 日志文件路径
		MaxSize:    100,                           // 每个日志文件最大大小（MB）
		MaxBackups: 5,                             // 保留旧日志文件的最大个数
		MaxAge:     30,                            // 保留旧日志文件的最大天数
		Compress:   true,                          // 是否压缩旧日志文件
	}

	newLogger.SetOutput(lumberjackLogger)

	loggers[filename] = newLogger
	return newLogger, nil
}
