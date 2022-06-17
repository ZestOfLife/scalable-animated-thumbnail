FROM golang:1.18-alpine

ADD ./ResizeProjector /app
WORKDIR /app

RUN go mod download

RUN go build -o /resizor

CMD [ "/resizor" ]
