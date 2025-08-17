package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

func InnitLogger(level string, isProduction bool) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	parsedLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		log.Warn().Msgf("Invalid log level '%s', defaulting to 'info'", level)
		parsedLevel = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(parsedLevel)

	if !isProduction {
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.RFC3339,
		}).With().Caller().Logger()
	} else {
		log.Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	}

	log.Info().Msgf("Logger initialized with level: %s, production mode: %t", parsedLevel.String(), isProduction)
}
