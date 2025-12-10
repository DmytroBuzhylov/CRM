package main

import (
	"Test/config"
	app "Test/internal/app"
	"github.com/rs/zerolog/log"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		return
	}

	a, err := app.NewApp(&cfg)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	a.Run()

}
