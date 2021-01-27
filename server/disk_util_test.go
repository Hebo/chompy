package server

import (
	"path"
	"testing"
	"time"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func Test_isOverSizeLimit(t *testing.T) {
	testVideos := []videoFile{
		{"2.mp4", time.Time{}, 5},
		{"3.mp4", time.Time{}, 10},
		{"4.mp4", time.Time{}, 20},
		{"5.mp4", time.Time{}, 40},
		{"6.mp4", time.Time{}, 1100},
	}

	type args struct {
		videos    []videoFile
		sizeLimit int64
	}
	tests := []struct {
		name  string
		args  args
		want  bool
		want1 int64
	}{
		{"no limit", args{testVideos, 0}, false, 0},
		{"limited1", args{testVideos, 10}, true, 1165},
		{"limited2", args{testVideos, 1000}, true, 175},
		{"under limit", args{testVideos, 100000}, false, -98825},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := needsDeletion(tt.args.videos, tt.args.sizeLimit)
			if got != tt.want {
				t.Errorf("isOverSizeLimit() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("isOverSizeLimit() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestDeleteVideoFiles(t *testing.T) {
	fs := afero.NewMemMapFs()
	toDeleteVideos := []videoFile{
		{"2.mp4", time.Time{}, 5},
		{"3.mp4", time.Time{}, 10},
		{"4.mp4", time.Time{}, 20},
		{"5.mp4", time.Time{}, 40},
		{"6.mp4", time.Time{}, 1100},
	}

	notToDelete := []string{
		"hello.txt",
		"important_files.rar",
		"somethingelse.mp4",
		".ytdl-config.txt",
	}

	dir := "/downloads"
	for _, v := range toDeleteVideos {
		err := afero.WriteFile(fs, path.Join(dir, v.Filename), []byte("file "+v.Filename), 0644)
		assert.NoError(t, err)
	}

	for _, v := range notToDelete {
		err := afero.WriteFile(fs, path.Join(dir, v), []byte("file "+v), 0644)
		assert.NoError(t, err)
	}

	err := deleteVideoFiles(fs, toDeleteVideos, dir)
	assert.NoError(t, err)

	files, err := afero.ReadDir(fs, dir)
	assert.NoError(t, err)

	var fnames []string
	for _, f := range files {
		fnames = append(fnames, f.Name())
	}

	for _, f := range notToDelete {
		assert.Contains(t, fnames, f, "should retain files not in deletion list")
	}
	for _, vid := range toDeleteVideos {
		assert.NotContains(t, fnames, vid.Filename, "should not contain deleted file")
	}
}
