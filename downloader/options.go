package downloader

import "fmt"

// Goals: <=1080p, reasonable size, iOS compatible, avoid merging if possible
// Below works decently well, but merges a lot?
// https://www.reddit.com/r/youtubedl/comments/fe08jx/can_youtubedl_download_only_mp4_files_at_1080_or/
var defaultFormat = stringOption{"--format", "bestvideo[height<=?1080]+bestaudio/best"}
var defaultFormatSort = stringOption{"--format-sort", "vcodec:h264,res,acodec:m4a"}

type ytdlopts []option

type option interface {
	toArg() string
}

type stringOption struct {
	Option string
	Value  string
}

func (o stringOption) toArg() string {
	return fmt.Sprintf("%s=%s", o.Option, o.Value)
}

type boolOption struct {
	Option string
}

func (o boolOption) toArg() string {
	return o.Option
}

func (o ytdlopts) ToCmdArgs() []string {
	var args []string
	for _, opt := range o {
		args = append(args, opt.toArg())
	}
	return args
}

func defaultOptions() ytdlopts {
	return ytdlopts{
		stringOption{"--retries", "3"},
		boolOption{"--no-progress"},
		boolOption{"--no-mtime"},
		boolOption{"--no-playlist"},
		boolOption{"--embed-subs"},
		stringOption{"--sub-langs", "all,-live_chat"},
		boolOption{"--netrc"},
		stringOption{"--match-filter", "!is_live & !was_live"},
		stringOption{"--sponsorblock-remove", "sponsor"},
		// Compatibility
		stringOption{"--compat-options", "multistreams"},
		stringOption{"--compat-options", "abort-on-error"},
	}
}
