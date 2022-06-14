FROM golang:1.18-alpine

ADD ./ExtractProjector /app
WORKDIR /app

RUN go mod download

RUN go build -o /extractor

CMD [ "/extractor" ]
