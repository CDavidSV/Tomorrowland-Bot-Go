FROM golang:1.22 as builder

# Install necessary dependencies
RUN apt-get update && apt-get -y upgrade
RUN apt-get -y install curl

# Install ffmpeg
RUN apt -y install ffmpeg

# Install yt-dlp
RUN curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /usr/local/bin/yt-dlp
RUN chmod a+rx /usr/local/bin/yt-dlp

# Install nodejs
RUN curl -fsSL https://deb.nodesource.com/setup_22.x -o nodesource_setup.sh
RUN bash nodesource_setup.sh
RUN apt-get install -y nodejs

WORKDIR /app
COPY . .

WORKDIR /app/tmrl-web

RUN npm install && npx playwright install && npx playwright install-deps
RUN chmod +x index.js
RUN npm link

WORKDIR /app

# Download dependencies and build project
RUN go mod download
RUN go build ./cmd/bot

# Run the compiled binary
CMD ["./bot"]