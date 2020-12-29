package server

import (
	"log"

	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
)

func (s Server) startWorkers() error {

	// general design thoughts
	//
	//  - explore how to handle config
	// - do we need to create config folder? contains:
	//        1. archive.txt
	//          2. config options
	//       needs research on best way to manage config and be easy to run
	//  - set archive.txt to config folder?
	//  - spin off workers for cleanup, regular playlist downloads
	//   - for GC, perhaps a return a channel / fn to call from api fns after each download
	//
	// TODOs
	// - set playlist url via config
	// - handle private playlists
	// - remove capturinglogger and replace w/ normal logger
	//
	var playlistTask func() = func() {}
	if true { // TODO
		playlistTask = func() {
			log.Println("Playlist task triggered")
			s.downloader.DownloadPlaylist("https://www.youtube.com/playlist?list=PLMM9FcCPG72z8fGbr-R4mLXebKcV45tkR")
		}
	}

	// Startup tasks
	playlistTask()

	scheduler := cron.New(
		cron.WithChain(
			cron.SkipIfStillRunning(cron.DiscardLogger),
		))

	_, err := scheduler.AddFunc("@every 31m", playlistTask)
	if err != nil {
		return errors.Wrap(err, "failed to schedule task")
	}

	scheduler.Start()
	return nil
}