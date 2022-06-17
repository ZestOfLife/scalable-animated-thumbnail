FROM golang:1.18-alpine

ADD ./StatusProjector /app
WORKDIR /app

RUN go mod download

RUN go build -o /status

CMD [ "/status" ]
