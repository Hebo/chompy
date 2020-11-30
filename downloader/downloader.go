package downloader

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"path"
	"regexp"

	"github.com/pkg/errors"
)

// Downloader handles downloading of video urls, wrapping youtube-dl
type Downloader struct {
	downloadsDir string
}

// New creates a new downloader
func New(path string) Downloader {
	dl := Downloader{downloadsDir: path}
	log.Println("creating downloader for path ", path)
	return dl
}

// Download fetches a single URL with youtube-dl, assuming youtube-dl is
// in $PATH
func (d Downloader) Download(url string) (string, error) {
	cmd := exec.Command("youtube-dl", url)
	opts := defaultOptions()
	opts = append(opts, stringOption{"--output", path.Join(d.downloadsDir, "%(title)s.%(ext)s")})

	cmd.Args = append(cmd.Args, opts.ToCmdArgs()...)

	log.Println("Running cmd", cmd.String())

	cmd.Stderr = os.Stderr
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		return "", errors.Wrap(err, "failed to create pipe")
	}

	scanner := bufio.NewScanner(cmdReader)
	pathChan := make(chan string)
	go capturingLogger(*scanner, pathChan)

	err = cmd.Start()
	if err != nil {
		return "", errors.Wrap(err, "error starting cmd")
	}

	var filePath string
	foundFilePath := false
	for p := range pathChan {
		foundFilePath = true
		filePath = p
	}

	err = cmd.Wait()
	if err != nil {
		return "", errors.Wrap(err, "error waiting on cmd")
	}

	if !foundFilePath {
		return "", errors.New("no output file path found")
	}

	return filePath, nil
}

var pathPatterns = []*regexp.Regexp{
	// [download] downloads/'99 Percent' Miss This. What Is The Length.mp4 has already been downloaded and merged
	regexp.MustCompile(`^\[download\][\s](.*?)[\s]has already.+$`),
	regexp.MustCompile(`^\[ffmpeg\] Merging formats into "(.*?)"$`),
}

// capturingLogger prints and scans for the output filepath. The most recent
// path found is assumed correct.
func capturingLogger(s bufio.Scanner, out chan<- string) {
	for s.Scan() {
		log.Println("youtube-dl ->", s.Text())
		if path, ok := matchLogPath(s.Text()); ok {
			out <- path
		}
	}

	close(out)
}

// matchLogPath looks for file paths in log lines by matching against regex
func matchLogPath(logLine string) (string, bool) {
	var path string
	found := false
	// log.Println("matching line", logLine)
	for _, r := range pathPatterns {
		if matches := r.FindStringSubmatch(logLine); matches != nil {
			log.Println("Matched path, ", matches[1])
			path = matches[1]
			found = true
		}
	}

	return path, found
}
