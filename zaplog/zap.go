package zaplog

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/pkg/errors"

	"go.uber.org/zap/zapcore"

	"github.com/go-kratos/kratos/v2/log"
	"go.uber.org/zap"
)

type LogWrapper func(level log.Level, keyvals ...interface{}) error

func (f LogWrapper) Log(level log.Level, keyvals ...interface{}) error {
	return f(level, keyvals...)
}

func NewZapLogger(opts ...Option) (log.Logger, error) {
	options := defaultOptions()
	for _, opt := range opts {
		opt(options)
	}

	cfg := zap.NewProductionConfig()

	if options.filePath != "" {
		path := options.filePath
		if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
			return nil, errors.WithStack(err)
		}

		if options.rotate != nil {
			// 创建编码器配置
			encoderConfig := zap.NewProductionEncoderConfig()
			encoderConfig.TimeKey = ""
			encoderConfig.CallerKey = ""
			encoderConfig.MessageKey = ""
			encoderConfig.StacktraceKey = ""

			zapLevel := getZapLevel(options.level)

			// 创建文件输出
			fileWriter := zapcore.AddSync(&lumberjack.Logger{
				Filename:   path,
				MaxSize:    options.rotate.maxSize,
				MaxBackups: options.rotate.maxBackups,
				MaxAge:     options.rotate.maxAge,
				Compress:   options.rotate.compress,
			})

			// 创建标准输出
			consoleWriter := zapcore.AddSync(os.Stdout)

			// 创建多个输出核心
			core := zapcore.NewTee(
				zapcore.NewCore(
					zapcore.NewJSONEncoder(encoderConfig),
					fileWriter,
					zapLevel,
				),
				zapcore.NewCore(
					zapcore.NewJSONEncoder(encoderConfig),
					consoleWriter,
					zapLevel,
				),
			)

			logger := zap.New(core)
			return wrapLogger(logger), nil
		} else {
			// 不使用日志轮转时直接配置多个输出路径
			cfg.OutputPaths = []string{"stdout", path}
		}
	} else {
		// 仅输出到标准输出
		cfg.OutputPaths = []string{"stdout"}
	}

	cfg.Level = zap.NewAtomicLevelAt(getZapLevel(options.level))
	cfg.EncoderConfig.TimeKey = ""
	cfg.EncoderConfig.CallerKey = ""
	cfg.EncoderConfig.MessageKey = ""
	cfg.EncoderConfig.StacktraceKey = ""

	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return wrapLogger(logger), nil
}

// 其他辅助函数保持不变
func getZapLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	case "fatal":
		return zap.FatalLevel
	default:
		return zap.InfoLevel
	}
}

func wrapLogger(logger *zap.Logger) LogWrapper {
	return func(level log.Level, keyvals ...interface{}) error {
		zapLevelMap := map[log.Level]zapcore.Level{
			log.LevelDebug: zap.DebugLevel,
			log.LevelInfo:  zap.InfoLevel,
			log.LevelWarn:  zap.WarnLevel,
			log.LevelError: zap.ErrorLevel,
			log.LevelFatal: zap.FatalLevel,
		}[level]

		fields := make([]zap.Field, len(keyvals)/2)
		for i := 0; i < len(keyvals); i += 2 {
			fields[i/2] = zap.String(fmt.Sprintf("%v", keyvals[i]), fmt.Sprintf("%v", keyvals[i+1]))
		}
		logger.Log(zapLevelMap, "", fields...)
		return nil
	}
}

func MaskPassword(input string) string {
	re := regexp.MustCompile(`(password:")([^"]*)(")`)
	result := re.ReplaceAllString(input, `${1}******${3}`)
	return result
}
