package main

import (
	_ "github.com/joho/godotenv/autoload"

	// standard
	"os"
	"strconv"

	// internal
	. "github.com/caian-org/lastfm-webp-widgets/internal/logger"

	"github.com/caian-org/lastfm-webp-widgets/internal/lastfm"
	"github.com/caian-org/lastfm-webp-widgets/internal/storage"
	"github.com/caian-org/lastfm-webp-widgets/internal/util"
)

func doRoutine() {
	Log.Info("program started")
	titleNeedsScroll := false

	framedir, err := os.MkdirTemp("/tmp", "frames-*")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(framedir)

	_, widgetPath := getWidgetLocation()
	Log.Info("got widget path", "path", widgetPath)

	Log.Info("initializing browser")
	initBrowser(widgetPath)

	Log.Info("initializing lastfm client")
	lastfm.Init()

	Log.Info("initializing storage")
	storage.Init()

	Log.Info("getting lastfm user info")
	lastFmUser := lastfm.GetUserInfo()

	Log.Info("getting lastfm user recent tracks")
	lastFmUserRecentTracks := lastfm.GetUserRecentTracks()
	if lastFmUserRecentTracks == nil {
		return
	}
	lastTrack := lastFmUserRecentTracks.Track[0]

	setUserURL(lastFmUser.Name)
	setUserStats(lastFmUser.PlayCount, util.UnixToHumanReadable(lastFmUser.Registered.Text))
	setTrackTitle(lastTrack.Name)
	setArtistName(lastTrack.Artist.Text)
	setAlbumCoverSource(lastTrack.Image[len(lastTrack.Image)-1].Text)

	if getTrackTitleSizeInPixels() > TRACK_TITLE_MAX_SIZE_PIXELS {
		Log.Info("size is too long, making it scrollable")
		setTrackTitleScrollable(lastTrack.Name)
		titleNeedsScroll = true
	}

	if lastTrack.Attr.NowPlaying == "true" {
		setUserListeningNow()
	} else {
		timeSince, err := strconv.ParseInt(lastTrack.Date.UTS, 10, 64)
		if err != nil {
			panic(err)
		}
		setUserListeningLastPlayed(util.UnixToRelativeHumanTime(timeSince))
	}

	waitPageToLoad()

	if titleNeedsScroll {
		animateWithScroll(framedir, getScrollTicksAmount())
	} else {
		animateSimple(framedir)
	}

	Log.Info("program finished")
}
