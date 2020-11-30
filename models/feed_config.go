package models

import "fmt"

type FeedConfigMap map[string]FeedConfig

func (c FeedConfig) Validate() error {
	if c.NumItems == 0 {
		return fmt.Errorf("feed Config %q: NumItems is not set", c.Title)
	}

	if c.Schedule == "" {
		return fmt.Errorf("feed Config %q: Schedule is not set", c.Title)
	}
	return nil
}

// FeedConfig is the configuration for a single Feed
type FeedConfig struct {
	Title    string
	Reddits  []string
	NumItems int `toml:"num_items"`
	// Schedule is in crontab syntax
	Schedule string `toml:"schedule"`
}
