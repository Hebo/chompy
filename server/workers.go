package server

import (
	"log"

	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
)

func (s Server) startWorkers() error {
	var playlistTask func() = func() {}
	if s.playlistSyncURL != "" {
		log.Printf("Tracking playlist: %s\n", s.playlistSyncURL)
		playlistTask = func() {
			log.Println("Playlist task triggered")
			s.downloader.DownloadPlaylist(s.playlistSyncURL)
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
