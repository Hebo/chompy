# build stage
FROM golang as build

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build ./cmd/chompy

FROM alpine:3.12.1

RUN set -x \
   && apk add --no-cache \
         tzdata \
         curl \
         ffmpeg \
         gnupg \
         python3 \
   # Install youtube-dl
   # https://github.com/rg3/youtube-dl
   && curl -sSLo /usr/local/bin/youtube-dl https://yt-dl.org/downloads/latest/youtube-dl \
   && curl -sSLo youtube-dl.sig https://yt-dl.org/downloads/latest/youtube-dl.sig \
   && gpg --keyserver keyserver.ubuntu.com --recv-keys '7D33D762FD6C35130481347FDB4B54CBA4826A18' \
   && gpg --keyserver keyserver.ubuntu.com --recv-keys 'ED7F5BF46B3BBED81C87368E2C393E0F18A9236D' \
   && gpg --verify youtube-dl.sig /usr/local/bin/youtube-dl \
   && chmod a+rx /usr/local/bin/youtube-dl \
   # Install yt-dlp
   && curl -sSLo /usr/local/bin/yt-dlp https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp \
   && chmod a+rx /usr/local/bin/yt-dlp \
   # Requires python -> python3.
   && ln -s /usr/bin/python3 /usr/bin/python \
   # Clean-up
   && rm youtube-dl.sig \
   && apk del curl gnupg

WORKDIR /app

COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app/ /app/

VOLUME [ "/downloads" ]
ENV DOWNLOADS_DIR="/downloads"
ENV PORT=8000
EXPOSE 8000
ENTRYPOINT ["/app/chompy"]
