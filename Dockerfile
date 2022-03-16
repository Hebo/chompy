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
         # base64
         coreutils \
         tzdata \
         curl \
         ffmpeg \
         python3 \
   # Install yt-dlp
   # https://github.com/yt-dlp/yt-dlp
   && curl -sSLo /usr/local/bin/yt-dlp https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp \
   && chmod a+rx /usr/local/bin/yt-dlp \
   # Requires python -> python3.
   && ln -s /usr/bin/python3 /usr/bin/python \
   # Clean-up
   && apk del curl

WORKDIR /app

COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app/ /app/

VOLUME [ "/downloads" ]
ENV DOWNLOADS_DIR="/downloads"
ENV PORT=8000
EXPOSE 8000
ENTRYPOINT ["/app/entrypoint.sh"]
CMD ["/app/chompy"]
