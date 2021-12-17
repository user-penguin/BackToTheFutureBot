FROM golang:1.17-alpine
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY command ./command
COPY config ./config
COPY message ./message
COPY reader ./reader
COPY state ./state
COPY currency ./currency
COPY main.go ./

RUN go build -o /docker-gs-ping

CMD [ "/docker-gs-ping" ]
