FROM golang:1.21rc3-bullseye

WORKDIR /app

# Add docker-compose-wait tool -------------------
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.7.2/wait ./wait
RUN chmod +x ./wait

COPY go.mod go.sum ./

RUN go mod download

COPY .env ./

COPY . .

RUN go build -tags=jsoniter -o bin/app cmd/web/*.go 

CMD ["./bin/app"]

