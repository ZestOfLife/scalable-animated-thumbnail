FROM golang:1.18

RUN apt-get update
RUN apt-get -y install imagemagick

ADD ./CompileWorker /app
WORKDIR /app

RUN go mod download

RUN go build -o /compile-worker

CMD [ "/compile-worker" ]
