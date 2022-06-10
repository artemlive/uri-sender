package main

import (
	"context"
	"github.com/artemlive/uri-sender/config"
	"github.com/artemlive/uri-sender/notifier"
	"github.com/artemlive/uri-sender/screenshoter"
	"github.com/go-co-op/gocron"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/alecthomas/kingpin.v2"
	"strings"
	"time"
)

var (
	logLevel   = kingpin.Flag("logLevel", "Log level").Default("info").Short('l').String()
	// I created another one flag to turn on the chromedp debug separately, cause it generates insane amount of messages
	debugChromeDP   = kingpin.Flag("dpDebug", "Turn on ChromeDP debug").Default("false").Short('d').Bool()
	configPath = kingpin.Flag("config", "Path to a config file").Default("config.json").Short('c').Envar("CONFIG_PATH").String()
)

func main() {
	kingpin.Parse()
	setLogLevel()
	creator := notifier.NewCreator()
	ctx := context.Background()
	cfg, err := config.NewConfig(*configPath)
	if err != nil {
		log.Fatal().Msgf("Couldn't read config %s, %s", *configPath, err)
	}
	log.Debug().Msgf("Successfully read the config file from %s", *configPath)
	// Read all notifiers and run the via scheduler
	// TODO: create a manager for this code
	for _, notify := range cfg.Notifiers {

		filePath := ""
		if len(notify.ScreenShot.URL) > 0 {
			filePath, err = screenshoter.MakeScreenshot(ctx, notify.ScreenShot.URL, notify.ScreenShot.HTMLElement, notify.ScreenShot.OutPath, notify.ScreenShot.Wait, *debugChromeDP)
			if err != nil {
				log.Error().Msgf("Couldn't make a screenshot '%s': %s, check the htmlElement property", notify.Type, err)
			}
		}
		s := gocron.NewScheduler(time.UTC)

		notificator, err := creator.CreateNotifier(notifier.Action(strings.ToUpper(notify.Type)), notifier.Message{Message: notify.Message, File: filePath}, notify.Recipients, *cfg)
		if err != nil {
			log.Error().Msgf("Couldn't initialize notificator '%s': %s", notify.Type, err)
			return
		}
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

func test() {
	log.Info().Msg("TEST")
}

func setLogLevel() {
	level := getLogLevel()
	zerolog.SetGlobalLevel(level)
}

func getLogLevel() zerolog.Level {
	level := zerolog.InfoLevel
	switch *logLevel {
	case "debug":
		level = zerolog.DebugLevel
	case "error":
		level = zerolog.ErrorLevel
	}
	return level
}
