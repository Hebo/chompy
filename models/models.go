package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Story is a single piece of content (post, link, etc.)
type Story struct {
	Title        string
	Link         string
	Hostname     string
	CommentsLink string
	NumComments  int
	Subreddit    string
	Text         string
}

// Block is a collection of stories. In the future, a digest may have multiple blocks.
type Block struct {
	Title   string
	Stories []Story
}

// ContentBlocks is a slice of Blocks
type ContentBlocks []Block

func (c ContentBlocks) Value() (driver.Value, error) {
	matched := true
	if !matched {
		return driver.Value(""), fmt.Errorf("number '%s' not a valid PhoneNumber format", c)
	}

	marshalled, err := json.Marshal(c)
	if err != nil {
		return driver.Value(""), fmt.Errorf("failed to marshal json: %e", err)
	}

	return driver.Value(marshalled), nil
}

func (c *ContentBlocks) Scan(src interface{}) error {
	var source []byte
	// let's support string and []byte
	switch src := src.(type) {
	case string:
		source = []byte(src)
	case []byte:
		source = src
	default:
		return errors.New("incompatible type for ContentBlocks")
	}
	err := json.Unmarshal(source, &c)
	if err != nil {
		return err
	}

	return nil
}

// Digest is a single item in the RSS feed, an individual newsletter
type Digest struct {
	ID        int           `db:"id"`
	FeedName  string        `db:"feed_name"`
	Title     string        `db:"title"`
	Content   ContentBlocks `db:"content"`
	CreatedAt time.Time     `db:"created_at"`
}

func (c ContentBlocks) String() string {
	var storyCount int
	var previewStory string
	for _, block := range c {
		if previewStory == "" && len(block.Stories) > 0 {
			previewStory = block.Stories[0].Title
		}
		storyCount += len(block.Stories)

	}

	return fmt.Sprintf("<%d Stories - %q>", storyCount, Truncate(previewStory, 20))
}

// Truncate truncates a string
func Truncate(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length] + "..."
}

var _ fmt.Stringer = ContentBlocks{}
