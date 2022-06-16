FROM golang:1.18

RUN apt-get update
RUN apt-get install ffmpeg

ADD ./ExtractWorker /app
WORKDIR /app

RUN go mod download

RUN go build -o /extract-worker

CMD [ "/extract-worker" ]
