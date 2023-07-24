FROM golang:1.21rc3-bullseye

WORKDIR /app

# Add docker-compose-wait tool -------------------
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.7.2/wait ./wait
RUN chmod +x ./wait

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o bin/app cmd/web/*.go

EXPOSE 8080

CMD ["./bin/app"]
