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

// taskCleanup limits the disk size of downloaded videos by deleting
// videos until the downloads directory size is within the specified size limit.
// Videos are deleted in strict reverse chronological order (i.e. taskCleanup
// will not attempt to delete larger videos before smaller ones).
func (s Server) taskCleanup() {
	if s.maxSize != 0 {
		log.Printf("Max downloads dir size set to %d MiB\n", s.maxSize)
	}

	videos, err := getVideoFiles(s.downloadsDir, createdAsc)
	if err != nil {
		log.Fatalln("Failed to get video files:", err)
	}

	if over, diff := needsDeletion(videos, int64(s.maxSize)); over {
		log.Printf("/!\\ MAX_SIZE is lower than current directory size (by %d MiB). "+
			"Please remove extra files manually", diff)

		s.maxSize = 0
	}

	for {
		if _, ok := <-s.cleanup; ok {
			log.Println("taskMaxSize task triggered")

			if s.maxSize == 0 {
				continue
			}

			vids, err := getVideoFiles(s.downloadsDir, createdAsc)
			if err != nil {
				log.Println("Failed to get video files:", err)
				continue
			}

			over, diff := needsDeletion(vids, int64(s.maxSize))
			if !over {
				continue
			}

			var reclaimed int64
			vidsToDelete := make([]videoFile, 0, 2)
			for _, v := range vids {
				if reclaimed >= diff {
					break
				}

				reclaimed += v.Size
				vidsToDelete = append(vidsToDelete, v)
			}

			log.Printf("Over max size by %d MiB. "+
				"Deleting %d videos to free %d MiB", diff, len(vidsToDelete), reclaimed)

			if err = deleteVideoFiles(s.fs, vidsToDelete, s.downloadsDir); err != nil {
				log.Println("Failed to delete videos:", err)
			}
		}
	}
}

func (s Server) triggerCleanup() {
	s.cleanup <- struct{}{}
}

func (s Server) startWorkers() error {
	// Scheduled tasks
	// ---------------
	scheduler := cron.New()

	jobPlaylistSync := cron.NewChain(
		cron.SkipIfStillRunning(cron.DiscardLogger),
	).Then(cron.FuncJob(s.taskPlaylistSync))

	if _, err := scheduler.AddJob("@every 31m", jobPlaylistSync); err != nil {
		return errors.Wrap(err, "failed to schedule task")
	}

	// Startup tasks
	// --------------
	if s.playlistSyncURL != "" {
		log.Printf("Tracking playlist: %s\n", s.playlistSyncURL)
	}
	go jobPlaylistSync.Run()
	go s.taskCleanup()

	scheduler.Start()
	return nil
}
