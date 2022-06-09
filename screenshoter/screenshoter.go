package screenshoter

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"os"
	"time"
)

// elementScreenshot takes a screenshot of a specific element.
func elementScreenshot(urlStr, sel string, delay time.Duration, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlStr),
		chromedp.Sleep(delay * time.Second),
		chromedp.Screenshot(sel, res, chromedp.NodeVisible, chromedp.NodeReady),
	}
}

func fullScreenshot(urlStr string, quality int, delay time.Duration, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlStr),
		chromedp.FullScreenshot(res, quality),
	}
}

func MakeScreenshot(ctxParent context.Context, url, htmlElement, outPath string, delay time.Duration, level zerolog.Level) (string, error) {
	// create context with debug if needed
	ctx := context.Background()
	cancel := func() {}
	ctxTimeout, cancel := context.WithTimeout(ctxParent, delay*2*time.Minute)
	if level == zerolog.DebugLevel {
		ctx, cancel = chromedp.NewContext(
			ctxTimeout,
			chromedp.WithDebugf(log.Debug().Msgf),
		)
	} else {
		ctx, cancel = chromedp.NewContext(
			ctxTimeout,
		)
	}
	defer cancel()

	// capture screenshot of an element
	var buf []byte
	if len(htmlElement) > 0 {
		log.Debug().Msgf("taking an element screenshot")
		if err := chromedp.Run(ctx, elementScreenshot(url, htmlElement, delay, &buf)); err != nil {
			return "", err
		}
	} else {
		// capture entire browser viewport, returning png with quality=90
		log.Debug().Msgf("taking a full screenshot")
		if err := chromedp.Run(ctx, fullScreenshot(url, 90, delay, &buf)); err != nil {
			return "", err
		}
	}
	if err := os.MkdirAll(outPath, os.ModePerm); err != nil {
		return "", err
	}
	unixTs := time.Now().Unix()
	outFile := fmt.Sprintf("%s/%d.png", outPath, unixTs)
	if err := ioutil.WriteFile(outFile, buf, 0o644); err != nil {
		return "", err
	}

	return outFile, nil
}
