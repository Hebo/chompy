package server

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/spf13/afero"
)

type ordering int

const (
	createdAsc ordering = iota + 1
	createdDesc
)

const toMiB = 1024 * 1024

func getVideoFiles(path string, order ordering) ([]videoFile, error) {
	var vids []videoFile

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return vids, err
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), ".") || file.IsDir() {
			continue
		}
		vids = append(vids, videoFile{
			Filename: file.Name(),
			Created:  file.ModTime(),
			Size:     file.Size() / toMiB,
		})
	}

	switch order {
	case createdAsc:
		sort.Slice(vids, func(i, j int) bool { return vids[i].Created.Before(vids[j].Created) })
	case createdDesc:
		sort.Slice(vids, func(i, j int) bool { return vids[i].Created.After(vids[j].Created) })
	}

	return vids, nil
}

// needsDeletion checks if the total size of videos is over the specified limit
func needsDeletion(videos []videoFile, max int64) (bool, int64) {
	var size, diff int64

	if max == 0 {
		return false, diff
	}

	for _, v := range videos {
		size += v.Size
	}

	diff = size - max
	if diff <= 0 {
		return false, diff
	}

	return true, diff
}

// deleteVideoFiles performs deletion on the specified videos
func deleteVideoFiles(fs afero.Fs, videos []videoFile, dir string) error {
	for _, v := range videos {
		log.Println("Removing", v.Filename)
		if err := fs.Remove(path.Join(dir, v.Filename)); err != nil {
			log.Printf("Error deleting file %q: %s\n", v.Filename, err)
		}
	}
	return nil
}

func touch(dir, filename string) error {
	videoPath := filepath.Join(dir, filepath.Clean(filename))

	_, err := os.Stat(videoPath)
	if err != nil {
		return err
	}
	currentTime := time.Now().Local()
	err = os.Chtimes(videoPath, currentTime, currentTime)
	if err != nil {
		return err
	}

	return nil
}
