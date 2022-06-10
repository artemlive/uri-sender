package screenshoter

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"os"
	"time"
)

const (
	contextWaitCoefficient = 3
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

func MakeScreenshot(ctxParent context.Context, url, htmlElement, outPath string, delay time.Duration, debugChromeDP bool) (string, error) {
	timeout := delay*contextWaitCoefficient*time.Second
	// create context with debug if needed
	ctxTimeout, cancel := context.WithTimeout(ctxParent, timeout)
	log.Debug().Msgf("Created context with the timeout: %s", timeout.String())
	ctx, cancel := chromedp.NewContext(
		ctxTimeout,
	)
	if debugChromeDP {
		ctx, cancel = chromedp.NewContext(
			ctxTimeout,
			// the debug is insane
			chromedp.WithDebugf(log.Debug().Msgf),
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
	log.Debug().Msgf("screenshot file has successfully been written: %s", outFile)

	return outFile, nil
}
