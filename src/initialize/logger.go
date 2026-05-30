package initialize

import (
	"fmt"
	"gin-admin/config"
	"gin-admin/global"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"strings"
)

func initLogger(cfg *config.Config) {
	var err error
	global.Log, err = NewLogger(cfg.Log)
	if err != nil {
		panic(fmt.Sprintf("日志初始化失败: %v", err))
	}
}

// levelLogger 用于按级别拆分日志文件
type levelLogger struct {
	level zapcore.Level
	name  string
}

func NewLogger(logConfig config.LogConfig) (*zap.Logger, error) {
	globalLevel := parseLogLevel(logConfig.Level)
	logDir := logConfig.FilePath

	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("创建日志目录失败: %w", err)
	}

	encoder := newEncoder(logConfig.Format)

	var cores []zapcore.Core

	switch strings.ToLower(logConfig.Output) {
	case "stdout":
		cores = append(cores, newStdoutCore(encoder, globalLevel))
	case "file":
		fileCores := newFileCores(logConfig, encoder, globalLevel, logDir)
		cores = append(cores, fileCores...)
	case "both":
		cores = append(cores, newStdoutCore(encoder, globalLevel))
		fileCores := newFileCores(logConfig, encoder, globalLevel, logDir)
		cores = append(cores, fileCores...)
	default:
		return nil, fmt.Errorf("不支持的输出类型: %s (可选: stdout/file/both)", logConfig.Output)
	}

	if len(cores) == 0 {
		return nil, fmt.Errorf("未能创建任何日志核心(Core)")
	}

	core := zapcore.NewTee(cores...)
	logger := zap.New(core, zap.AddCaller())

	return logger, nil
}

// parseLogLevel 将字符串日志级别转换为 zapcore.Level
func parseLogLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	case "panic":
		return zapcore.PanicLevel
	default:
		return zapcore.InfoLevel
	}
}

// newEncoder 根据格式创建编码器
func newEncoder(format string) zapcore.Encoder {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	if strings.ToLower(format) == "json" {
		return zapcore.NewJSONEncoder(encoderConfig)
	}
	return zapcore.NewConsoleEncoder(encoderConfig)
}

// newStdoutCore 创建标准输出 Core
func newStdoutCore(encoder zapcore.Encoder, level zapcore.Level) zapcore.Core {
	return zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout),
		level,
	)
}

// newFileCores 创建文件输出 Cores（按级别分文件）
func newFileCores(logConfig config.LogConfig, encoder zapcore.Encoder, globalLevel zapcore.Level, logDir string) []zapcore.Core {
	levelLoggers := []levelLogger{
		{zapcore.DebugLevel, "debug"},
		{zapcore.InfoLevel, "info"},
		{zapcore.WarnLevel, "warn"},
		{zapcore.ErrorLevel, "error"},
	}

	var cores []zapcore.Core
	for _, ll := range levelLoggers {
		if ll.level < globalLevel {
			continue
		}

		writer := &lumberjack.Logger{
			Filename:   fmt.Sprintf("%s/%s.log", logDir, ll.name),
			MaxSize:    logConfig.MaxSize,
			MaxBackups: logConfig.MaxBackups,
			MaxAge:     logConfig.MaxAge,
			Compress:   logConfig.Compress,
		}

		core := zapcore.NewCore(
			encoder,
			zapcore.AddSync(writer),
			ll.level,
		)
		cores = append(cores, core)
	}

	return cores
}
