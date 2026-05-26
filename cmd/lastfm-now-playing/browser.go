package main

import (
	// standard
	"fmt"
	"strings"
	"time"

	// internal
	. "github.com/caian-org/lastfm-webp-widgets/internal/logger"
	B "github.com/caian-org/lastfm-webp-widgets/internal/widget_browser"
)

var browser *B.Browser

func initBrowser(fp string) {
	browser = B.New(WIDGET_WIDTH, WIDGET_HEIGHT, WIDGET_PIXEL_RATIO)
	Log.Info("browser initialized")

	widgetURI := fmt.Sprintf("file://%s", fp)
	browser.GoTo(widgetURI)
	Log.Info("went to widget location", "uri", widgetURI)
}

/* ---------------------------------------------------------------------------------------------- */

func getTrackTitleSizeInPixels() int {
	size := browser.
		GetCurrentPage().
		MustEval("() => music.getTrackTitleSizeInPixelsRounded()").
		Int()

	Log.Info("got track title size in pixels", "size", size)
	return size
}

func getScrollTicksAmount() int {
	ticks := browser.
		GetCurrentPage().
		MustEval("() => music.getScrollTicksAmount()").
		Int()

	Log.Info("got scroll ticks amount", "ticks", ticks)
	return ticks
}

func setTrackTitle(title string) {
	tt := browser.
		GetCurrentPage().
		MustEval("(t) => music.setTrackTitle(t)", title).
		Str()

	Log.Info("track title setted", "title", tt)
}

func setTrackTitleScrollable(title string) {
	whitespace := strings.Repeat("&nbsp;", TRACK_TITLE_SCROLLABLE_WHITESPACES)
	titleScrollable := fmt.Sprintf("%s%s%s", title, whitespace, title)

	setTrackTitle(titleScrollable)
}

func setArtistName(name string) {
	artist := browser.
		GetCurrentPage().
		MustEval("(n) => music.setArtistName(n)", name).
		Str()

	Log.Info("artist name setted", "name", artist)
}

func setAlbumCoverSource(src string) {
	cover := browser.
		GetCurrentPage().
		MustEval("(s) => music.setAlbumCoverSource(s)", src).
		Str()

	Log.Info("album cover source setted", "src", cover)
}

func tickTitleTrackScroll() {
	browser.
		GetCurrentPage().
		MustEval("() => music.tickScroll()")
}

/* ---------------------------------------------------------------------------------------------- */

func setUserURL(username string) {
	url := browser.
		GetCurrentPage().
		MustEval("(u) => user.setUrl(u)", username).
		Str()

	Log.Info("user URL setted", "url", url)
}

func setUserStats(scrobbles string, creation string) {
	stats := browser.
		GetCurrentPage().
		MustEval("(s, c) => user.setStats(s, c)", scrobbles, creation).
		Str()

	Log.Info("user stats setted", "stats", stats)
}

func setUserListeningNow() {
	status := browser.
		GetCurrentPage().
		MustEval("() => user.setListeningStatusNowPlaying()").
		Str()

	Log.Info("user listening status setted", "status", status)
}

func setUserListeningLastPlayed(relativeTime string) {
	status := browser.
		GetCurrentPage().
		MustEval("(t) => user.setListeningStatusLastPlayed(t)", relativeTime).
		Str()

	Log.Info("user listening status setted", "status", status)
}

/* ---------------------------------------------------------------------------------------------- */

func randomizeSoundWave() {
	browser.
		GetCurrentPage().
		MustEval("() => waves.randomize()")
}

/* ---------------------------------------------------------------------------------------------- */

func setThemeMode(mode string) {
	browser.
		GetCurrentPage().
		MustEval("(m) => theme.set(m)", mode)
}

/* ---------------------------------------------------------------------------------------------- */

func waitPageToLoad() {
	Log.Info("waiting for page to load")

	start := time.Now()
	browser.GetCurrentPage().MustWaitStable()
	diffMS := time.Since(start).Milliseconds()

	Log.Info("page loaded", "time", diffMS)
}
