package downloader

import "testing"

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
			"", false},
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
