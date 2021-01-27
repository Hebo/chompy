package config

// Config contains Chompy's config. What were you expecting?
type Config struct {
	Port int `env:"PORT" envDefault:"8000"`
	// Directory for video downloads
	DownloadsDir    string `env:"DOWNLOADS_DIR" envDefault:"./downloads"`
	PlaylistSyncURL string `env:"PLAYLIST_SYNC"`
	Format          string `env:"FORMAT"`
}
