package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// InitLogger initializes the structured logger
func InitLogger() {
	println("📝 Logger başlatılıyor...")

	// Set log level from environment variable
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	println("📊 Log seviyesi:", logLevel)

	// Parse log level
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		println("⚠️ Log seviyesi parse edilemedi, varsayılan 'info' kullanılıyor")
		level = zerolog.InfoLevel
	}

	// Configure logger
	zerolog.SetGlobalLevel(level)
	zerolog.TimeFieldFormat = time.RFC3339

	// Pretty console output for development
	if os.Getenv("ENV") == "development" {
		println("🖥️ Development modu - console output")
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		})
	} else {
		println("🏭 Production modu - JSON output")
		// JSON output for production
		log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	}

	// Set global logger
	zerolog.DefaultContextLogger = &log.Logger
	println("✅ Logger başlatıldı")
}

// GetLogger returns the configured logger
func GetLogger() zerolog.Logger {
	return log.Logger
}

// Log levels
func Info(msg string, fields map[string]interface{}) {
	log.Info().Fields(fields).Msg(msg)
}

func Error(msg string, err error, fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	if err != nil {
		fields["error"] = err.Error()
	} else {
		fields["error"] = "unknown error"
	}
	log.Error().Fields(fields).Msg(msg)
}

func Warn(msg string, fields map[string]interface{}) {
	log.Warn().Fields(fields).Msg(msg)
}

func Debug(msg string, fields map[string]interface{}) {
	log.Debug().Fields(fields).Msg(msg)
}

func Fatal(msg string, err error, fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	fields["error"] = err.Error()
	log.Fatal().Fields(fields).Msg(msg)
}
