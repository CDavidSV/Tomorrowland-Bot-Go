FROM golang:1.22 as builder

WORKDIR /app
COPY . .

RUN go mod download

RUN go build ./cmd/bot

RUN apt-get update
RUN apt-get -y install ffmpeg
RUN apt-get -y install yt-dlp

CMD ["./bot"]