package server

import (
	"log"

	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
)

func (s Server) taskPlaylistSync() {
	if s.playlistSyncURL == "" {
		return
	}

	log.Println("PlaylistSync task triggered")
	if err := s.downloader.DownloadPlaylist(s.playlistSyncURL); err != nil {
		log.Println("Error downloading playlist:", err)
	}
}

func (s Server) startWorkers() error {

	// Startup tasks
	if s.playlistSyncURL != "" {
		log.Printf("Tracking playlist: %s\n", s.playlistSyncURL)
		s.taskPlaylistSync()
	}

	// Scheduled tasks
	scheduler := cron.New(
		cron.WithChain(
			cron.SkipIfStillRunning(cron.DiscardLogger),
		))

	_, err := scheduler.AddFunc("@every 31m", s.taskPlaylistSync)
	if err != nil {
		return errors.Wrap(err, "failed to schedule task")
	}

	scheduler.Start()
	return nil
}
