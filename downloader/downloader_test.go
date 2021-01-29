package downloader

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_matchLogPath(t *testing.T) {
	tests := []struct {
		name    string
		logLine string
		want    string
		wantOK  bool
	}{
		{"download1", "[download] downloads/'99 Percent' Miss This. What Is The Length.mp4 has already been downloaded and merged",
			"downloads/'99 Percent' Miss This. What Is The Length.mp4", true},
		{"merge1", "[ffmpeg] Merging formats into \"downloads/Dumbbell Romanian deadlift.mp4\"",
			"downloads/Dumbbell Romanian deadlift.mp4", true},
		{"download2", "[download] Destination: downloads/Dumbbell Romanian deadlift.f135.mp4",
			"downloads/Dumbbell Romanian deadlift.f135.mp4", true},
		{"download3", "[download] Destination: downloads/How to Protect Your Shopping Trolley From Improvised Explosives.webm",
			"downloads/How to Protect Your Shopping Trolley From Improvised Explosives.webm", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := matchLogPath(tt.logLine)
			if got != tt.want {
				t.Errorf("matchLogPath() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.wantOK {
				t.Errorf("matchLogPath() got1 = %v, want %v", got1, tt.wantOK)
			}
		})
	}
}

func TestNew(t *testing.T) {
	downloader := New("path", "")
	want := Downloader{downloadsDir: "path", format: stringOption{"--format", "bestvideo[ext=mp4][height<=?1080]+bestaudio[ext=m4a]/best"}}
	assert.Equal(t, downloader, want)

	downloader = New("/downloads", "bestvideo")
	want = Downloader{downloadsDir: "/downloads", format: stringOption{"--format", "bestvideo"}}
	assert.Equal(t, downloader, want)

}
