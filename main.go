package main

import (
	"github.com/artemlive/uri-sender/config"
	"github.com/artemlive/uri-sender/notifier"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/alecthomas/kingpin.v2"
	"strings"
)


var (
	logLevel   = kingpin.Flag("logLevel", "Log level").Default("info").String()
	configPath   = kingpin.Flag("config", "Path to a config file").Default("config.json").Short('c').Envar("CONFIG_PATH").String()
)

func main(){
	kingpin.Parse()
	setLogLevel()
	creator := notifier.NewCreator()

	cfg, err := config.NewConfig(*configPath)
	if err != nil {
		log.Fatal().Msgf("Couldn't read config %s, %s", *configPath, err)
	}
	for _, notify := range cfg.Notifiers {
		notificator, err := creator.CreateNotifier(notifier.Action(strings.ToUpper(notify.Type)), notifier.Message{Message: notify.Message}, notify.Recipients, *cfg)
		if err != nil {
			log.Error().Msgf("Couldn't initialize notificator '%s': %s", notify.Type, err)
			return
		}
		err = notificator.Send()
		if err != nil {
			log.Fatal().Msgf("Couldn't send message via '%s' %s", notify.Type, err)
			return
		}
	}
}

func setLogLevel(){
	level := zerolog.InfoLevel
	switch *logLevel {
	case "debug":
		level = zerolog.DebugLevel
	case "error":
		level = zerolog.ErrorLevel
	}
	zerolog.SetGlobalLevel(level)
}