package screenshoter

import (
	"context"
	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"time"
)

// elementScreenshot takes a screenshot of a specific element.
func elementScreenshot(urlStr, sel string, delay time.Duration,res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlStr),
		chromedp.Sleep(delay * time.Second),
		chromedp.Screenshot(sel, res, chromedp.NodeVisible, chromedp.NodeReady),
	}
}

func MakeScreenshot(url, htmlElement, outPath string, delay time.Duration, level zerolog.Level) error{
	// create context with debug if needed
	ctx := context.Background()
	cancel := func(){}
	if level == zerolog.DebugLevel {
		ctx, cancel = chromedp.NewContext(
			context.Background(),
			chromedp.WithDebugf(log.Debug().Msgf),
		)
	} else {
		ctx, cancel = chromedp.NewContext(
			context.Background(),
		)
	}
	defer cancel()

	// capture screenshot of an element
	var buf []byte
	if err := chromedp.Run(ctx, elementScreenshot(url, htmlElement, delay, &buf)); err != nil {
		log.Error().Msgf("couldn't make a screenshot: %s", err)
		return err
	}
	if err := ioutil.WriteFile(outPath, buf, 0o644); err != nil {
		log.Error().Msgf("couldn't write screenshot file: %s", err)
		return err
	}

	log.Printf("wrote elementScreenshot.png and fullScreenshot.png")
	return nil
}