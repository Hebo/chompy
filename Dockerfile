# build stage
FROM golang AS build

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build ./cmd/chompy

FROM alpine:3.21.3 AS install

RUN set -x \
   && apk add --no-cache \
         # base64
         coreutils \
         tzdata \
         curl \
         ffmpeg \
         python3 \
   # Install yt-dlp
   # https://github.com/yt-dlp/yt-dlp
   && curl -sSLo /usr/local/bin/yt-dlp https://github.com/yt-dlp/yt-dlp-nightly-builds/releases/latest/download/yt-dlp \
   && chmod a+rx /usr/local/bin/yt-dlp \
   # Clean-up
   && apk del curl

RUN addgroup -g 1000 app \
    && adduser -u 1000 -G app -D app
USER app
WORKDIR /app

RUN /usr/local/bin/yt-dlp --version > /app/YTDLP_VERSION

COPY --chown=app:app --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --chown=app:app --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --chown=app:app --from=build /app/ /app/

VOLUME [ "/downloads" ]
ENV DOWNLOADS_DIR="/downloads"
ENV PORT=8000
EXPOSE 8000
ENTRYPOINT ["/app/entrypoint.sh"]
CMD ["/app/chompy"]
