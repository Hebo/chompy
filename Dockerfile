# build stage
FROM golang as build

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build ./cmd/chompy

FROM alpine:3.12.1

RUN apk add --no-cache ffmpeg youtube-dl

WORKDIR /app

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app/ /app/

EXPOSE 8000
ENTRYPOINT ["/app/chompy", "-port", "8000"]
