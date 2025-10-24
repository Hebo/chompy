# build stage
FROM golang AS build

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build ./cmd/chompy

FROM debian:bookworm-slim AS install

RUN set -x \
      && apt-get update \
      && apt-get install -y --no-install-recommends \
            coreutils \
            tzdata \
            curl \
            ffmpeg \
            python3 \
            ca-certificates \
      # Install yt-dlp
      # https://github.com/yt-dlp/yt-dlp
      && curl -sSLo /usr/local/bin/yt-dlp https://github.com/yt-dlp/yt-dlp-nightly-builds/releases/latest/download/yt-dlp \
      && chmod a+rx /usr/local/bin/yt-dlp \
      # Clean-up
      && apt-get remove -y curl \
      && apt-get autoremove -y \
      && apt-get clean


WORKDIR /app

RUN /usr/local/bin/yt-dlp --version > /app/YTDLP_VERSION

COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=denoland/deno:bin /deno /usr/local/bin/deno
COPY --from=build /app/ /app/

VOLUME [ "/downloads" ]
ENV DOWNLOADS_DIR="/downloads"
ENV PORT=8000
EXPOSE 8000
ENTRYPOINT ["/app/entrypoint.sh"]
CMD ["/app/chompy"]
