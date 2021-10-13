package xlog

type Option func(logger *Logger)

func WithTrace(traceName string) Option {
	return func(l *Logger) {
		l.traceName = traceName
	}
}

func WithFeishu(feishuUrl string) Option {
	return func(l *Logger) {
		l.feishuUrl = feishuUrl
	}
}

func WithInfo(serviceName, env string) Option {
	return func(l *Logger) {
		l.serviceName = serviceName
		l.env = env
	}
}

func WithDebug() Option {
	return func(l *Logger) {
		l.debug = true
	}
}
