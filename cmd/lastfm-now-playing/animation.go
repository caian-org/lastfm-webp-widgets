package main

import (
	// standard
	"bytes"
	"fmt"
	"path/filepath"

	// internal
	. "github.com/caian-org/lastfm-webp-widgets/internal/logger"

	"github.com/caian-org/lastfm-webp-widgets/internal/storage"
	"github.com/caian-org/lastfm-webp-widgets/internal/util"
	"github.com/caian-org/lastfm-webp-widgets/pkg/webpanimation"
)

const (
	THEME_MODE_LIGHT = "light"
	THEME_MODE_DARK  = "dark"
)

var (
	THEME_MODES = []string{THEME_MODE_LIGHT, THEME_MODE_DARK}
)

func takeScreenshot(framedir string, counter int, tm string, frames *map[string][]string) {
	ffp := filepath.Join(framedir, fmt.Sprintf("%d-%s.png", counter, tm))
	browser.TakeScreenshot(ffp)

	(*frames)[tm] = append((*frames)[tm], ffp)
	Log.Info("took screenshot", "type", "scroll", "mode", tm, "path", ffp)
}

func animateWithScroll(framedir string, ticks int) {
	frames := map[string][]string{}
	soundWaveRefreshesBetweenTicks := 3

	for _, tm := range THEME_MODES {
		frames[tm] = []string{}
	}

	counter := 0

	for i := 0; i <= ticks; i++ {
		for j := 1; j <= soundWaveRefreshesBetweenTicks; j++ {
			randomizeSoundWave()

			if j == soundWaveRefreshesBetweenTicks {
				tickTitleTrackScroll()
			}

			counter += 1

			for _, tm := range THEME_MODES {
				setThemeMode(tm)
				takeScreenshot(framedir, counter, tm, &frames)
			}
		}
	}

	for k := 0; k < soundWaveRefreshesBetweenTicks*8; k++ {
		randomizeSoundWave()
		counter += 1

		for _, tm := range THEME_MODES {
			setThemeMode(tm)
			takeScreenshot(framedir, counter, tm, &frames)
		}
	}

	for _, tm := range THEME_MODES {
		animate(frames[tm], fmt.Sprintf("lastfm-now-playing-%s.webp", tm))
	}
}

func animateSimple(framedir string) {
	frames := map[string][]string{}

	for _, tm := range THEME_MODES {
		frames[tm] = []string{}
	}

	for i := 0; i < 10; i++ {
		randomizeSoundWave()

		for _, tm := range THEME_MODES {
			setThemeMode(tm)
			takeScreenshot(framedir, i, tm, &frames)
		}
	}

	for _, tm := range THEME_MODES {
		animate(frames[tm], fmt.Sprintf("lastfm-now-playing-%s.webp", tm))
	}
}

func animate(frames []string, output string) {
	wpa := webpanimation.NewWebpAnimation(
		WIDGET_WIDTH*WIDGET_PIXEL_RATIO,
		WIDGET_HEIGHT*WIDGET_PIXEL_RATIO,
		0,
	)

	// see: <https://developers.google.com/speed/webp/docs/gif2webp#:~:text=%2Dkmin%20int%2C%20%2Dkmax%20int>
	wpa.WebPAnimEncoderOptions.SetKmin(9)
	wpa.WebPAnimEncoderOptions.SetKmax(17)

	defer wpa.ReleaseMemory()

	webpConfig := webpanimation.NewWebpConfig()
	webpConfig.SetLossless(1)

	timeline := 0
	for _, f := range frames {
		frame, err := util.LoadImageFromFS(f)
		if err != nil {
			panic(err)
		}

		err = wpa.AddFrame(frame, timeline, webpConfig)
		if err != nil {
			panic(err)
		}

		timeline += 200
	}

	err := wpa.AddFrame(nil, timeline, webpConfig)
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	err = wpa.Encode(&buf)
	if err != nil {
		panic(err)
	}

	storage.Save(output, buf)
	Log.Info("webp animation saved")
}
