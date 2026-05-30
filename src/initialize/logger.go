package initialize

import (
	"fmt"
	"gin-admin/config"
	"gin-admin/global"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func initLogger(cfg *config.Config) {
	var err error
	global.Log, err = NewLogger(cfg.Log)
	if err != nil {
		panic(fmt.Sprintf("日志初始化失败: %v", err))
	}
}

// 定义日志分级配置：文件名 + 对应最低输出级别
type levelLogger struct {
	level zapcore.Level
	name  string
}

func NewLogger(logConfig config.LogConfig) (*zap.Logger, error) {
	// 1. 解析日志级别
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(logConfig.Level)); err != nil {
		level = zapcore.InfoLevel
	}

	// 2. 编码器配置（保留原格式：json/console）
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
	var encoder zapcore.Encoder
	if logConfig.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// 3. 生成【按天】的日志目录：logs/2026-05-30
	today := time.Now().Format("2006-01-02")
	logDir := fmt.Sprintf("logs/%s", today)
	// 自动创建目录（不存在则创建）
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("创建日志目录失败: %w", err)
	}

	// 4. 定义分级日志规则（行业标准：高等级日志只输出对应及以上级别）
	levelLoggers := []levelLogger{
		{zapcore.DebugLevel, "debug"}, // 所有日志
		{zapcore.InfoLevel, "info"},   // info及以上
		{zapcore.WarnLevel, "warn"},   // warn及以上
		{zapcore.ErrorLevel, "error"}, // 仅error
	}

	var cores []zapcore.Core
	// 5. 为每个级别创建独立的日志写入器（带切割）
	for _, ll := range levelLoggers {
		// 日志文件完整路径
		logFilePath := fmt.Sprintf("%s/%s.log", logDir, ll.name)

		// 复用lumberjack实现日志切割/备份/压缩
		writer := &lumberjack.Logger{
			Filename:   logFilePath,
			MaxSize:    logConfig.MaxSize,    // 单个文件最大大小(MB)
			MaxBackups: logConfig.MaxBackups, // 最大备份数
			MaxAge:     logConfig.MaxAge,     // 保留天数
			Compress:   logConfig.Compress,   // 是否压缩
		}

		// 创建带【级别过滤】的core
		core := zapcore.NewCore(
			encoder,
			zapcore.AddSync(writer),
			ll.level, // 核心：按级别过滤日志
		)
		cores = append(cores, core)
	}

	// 6. 合并所有core（多输出核心）
	core := zapcore.NewTee(cores...)

	// 7. 修复调用栈错误！！！原skip=2导致显示runtime/asm，改为1即可显示业务代码行号
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	return logger, nil
}
