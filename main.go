package main

import (
	"github.com/artemlive/uri-sender/config"
	"github.com/artemlive/uri-sender/notifier"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/alecthomas/kingpin.v2"
	"strings"
	"time"
)

var (
	logLevel   = kingpin.Flag("logLevel", "Log level").Default("info").String()
	configPath = kingpin.Flag("config", "Path to a config file").Default("config.json").Short('c').Envar("CONFIG_PATH").String()
)

func main() {
	kingpin.Parse()
	setLogLevel()
	creator := notifier.NewCreator()

	cfg, err := config.NewConfig(*configPath)
	if err != nil {
		log.Fatal().Msgf("Couldn't read config %s, %s", *configPath, err)
	}

	// Read all notifiers and run the via scheduler
	for _, notify := range cfg.Notifiers {
		notificator, err := creator.CreateNotifier(notifier.Action(strings.ToUpper(notify.Type)), notifier.Message{Message: notify.Message}, notify.Recipients, *cfg)
		if err != nil {
			log.Error().Msgf("Couldn't initialize notificator '%s': %s", notify.Type, err)
			return
		}
		s := gocron.NewScheduler(time.UTC)
		_, err = s.Cron(notify.Cron).Do(notificator.Send)
		if err != nil {
			log.Fatal().Msgf("Couldn't send message via '%s' %s", notify.Type, err)
			return
		}
		s.StartAsync()
	}
	// infinite loop for scheduler
	for {
		time.Sleep(1 * time.Second)
	}

}

func test(){
	log.Info().Msg("TEST")
}
func setLogLevel() {
	level := zerolog.InfoLevel
	switch *logLevel {
	case "debug":
		level = zerolog.DebugLevel
	case "error":
		level = zerolog.ErrorLevel
	}
	zerolog.SetGlobalLevel(level)
}
