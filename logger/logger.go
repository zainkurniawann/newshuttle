package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

var Log zerolog.Logger

func InitLogger() {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	Log = zerolog.New(output).With().Timestamp().Logger()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

// For logging error
func LogError(err error, message string, details map[string]interface{}) {
	Log.Error().
		Err(err).
		Str("message", message).
		Interface("details", details).
		Msg(message)
}

// For logging info
func LogInfo(message string, details map[string]interface{}) {
	Log.Info().
		Str("message", message).
		Interface("details", details).
		Msg(message)
}

// For logging warning
func LogWarn(message string, details map[string]interface{}) {
	Log.Warn().
		Str("message", message).
		Interface("details", details).
		Msg(message)
}

// For logging debug
func LogDebug(message string, details map[string]interface{}) {
	Log.Debug().
		Str("message", message).
		Interface("details", details).
		Msg(message)
}

// For logging fatal
func LogFatal(err error, message string, details map[string]interface{}) {
	Log.Fatal().
		Err(err).
		Str("message", message).
		Interface("details", details).
		Msg(message)
}