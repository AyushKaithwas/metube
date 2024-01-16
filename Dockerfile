FROM golang:1.21

# Install FFmpeg
RUN apt-get update && apt-get install -y ffmpeg

WORKDIR /app

COPY . .

RUN go mod download

# RUN CGO_ENABLED=0 go build -o server .
RUN go build -o main cmd/api/main.go

EXPOSE 8080

CMD ["./main"]