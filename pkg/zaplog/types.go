package zaplog

type Options struct {
	filePath string
	level    string
	rotate   *RotateOptions
}

// RotateOptions 日志轮转的配置选项
type RotateOptions struct {
	maxSize    int  // 每个日志文件的最大大小(MB)
	maxBackups int  // 保留的旧日志文件最大数量
	maxAge     int  // 保留的旧日志文件最大天数
	compress   bool // 是否压缩旧日志文件
}

type Option func(*Options)

// defaultOptions 默认配置
func defaultOptions() *Options {
	return &Options{
		filePath: "",
		level:    "info",
		rotate:   nil, // 默认不启用日志轮转
	}
}

// WithFilePath 设置日志文件路径
func WithFilePath(path string) Option {
	return func(o *Options) {
		o.filePath = path
	}
}

// WithLevel 设置日志级别
func WithLevel(level string) Option {
	return func(o *Options) {
		o.level = level
	}
}

// WithRotate 启用日志轮转并设置相关配置
func WithRotate(maxSize, maxBackups, maxAge int, compress bool) Option {
	return func(o *Options) {
		o.rotate = &RotateOptions{
			maxSize:    maxSize,
			maxBackups: maxBackups,
			maxAge:     maxAge,
			compress:   compress,
		}
	}
}
