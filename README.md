# Chompy

![Build](https://github.com/hebo/chompy/workflows/gobuild/badge.svg)

**Download and watch videos easily on iOS**

Chompy wraps [youtube-dl](https://youtube-dl.org/) in an API, allowing ad-free downloading and streaming on devices that can't run youtube-dl directly, such as iOS.

## Usage

Better docs to come soon(TM)

Deploy me via Docker, and call `/download` and `/videos`

### Video Formats

The default format for downloaded videos is mp4, at resolutions up to 1080p. You can see the exact format string in [downloader/options.go](downloader/options.go).

Set the format for a download with the `format` request parameter. For instance:

```
http -v post localhost:8000/download url="https://www.youtube.com/watch?v=L5emxkKNf9Y" format='worstvideo'
```

## Development

### Run locally

**Dependencies:** [ffmpeg](https://ffmpeg.org/) and [youtube-dl](https://youtube-dl.org/). API examples use [HTTPie](https://httpie.io/)

Run the app
```
go run ./cmd/chompy
```

Download something exciting

```
http -v post localhost:8000/download url="https://www.youtube.com/watch?v=L5emxkKNf9Y"
```

```
HTTP/1.1 200 OK
Content-Length: 83
Content-Type: application/json; charset=UTF-8
Date: Mon, 30 Nov 2020 19:27:53 GMT

{
    "filename": "How to Protect Your Shopping Trolley From Improvised Explosives.mp4"
}
```

Then play it
```
http -v localhost:8000/videos filename=='How to Protect Your Shopping Trolley From Improvised Explosives.mp4'
```




### Docker

```
docker build -t chompy .

docker run -p 8000:8000 chompy
```
