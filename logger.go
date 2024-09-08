package gomodbus

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

// Define color codes
const (
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorReset  = "\033[0m"
)

// CustomColorEncoder is a custom zapcore.Encoder that adds color to log levels
type CustomColorEncoder struct {
	zapcore.Encoder
}

func (c *CustomColorEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	buf, err := c.Encoder.EncodeEntry(entry, fields)
	if err != nil {
		return nil, err
	}

	// Add color based on log level
	var color string
	switch entry.Level {
	case zapcore.DebugLevel:
		color = colorBlue
	case zapcore.InfoLevel:
		color = colorGreen
	case zapcore.WarnLevel:
		color = colorYellow
	case zapcore.ErrorLevel, zapcore.DPanicLevel, zapcore.PanicLevel, zapcore.FatalLevel:
		color = colorRed
	default:
		color = colorReset
	}

	// Prepend color and append reset color code
	coloredBuf := buffer.NewPool().Get()
	coloredBuf.AppendString(color)
	coloredBuf.AppendString(buf.String())
	coloredBuf.AppendString(colorReset)

	return coloredBuf, nil
}

// InitializeLogger initializes a reusable zap logger with color support
func InitializeLogger() *zap.Logger {
	// Create a custom encoder config
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Create a custom color encoder
	customEncoder := &CustomColorEncoder{
		Encoder: zapcore.NewConsoleEncoder(encoderConfig),
	}

	// Create a core with the custom encoder
	core := zapcore.NewCore(customEncoder, zapcore.AddSync(os.Stdout), zap.DebugLevel)

	// Create and return a logger with the custom core
	return zap.New(core)
}

// Global logger instance
var Logger *zap.Logger

func init() {
	Logger = InitializeLogger()
}

// EnableDebug enables debug level logging for the global logger
func EnableDebug() {
	if Logger != nil {
		Logger = Logger.WithOptions(zap.IncreaseLevel(zap.DebugLevel))
	}
}