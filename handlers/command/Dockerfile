FROM golang:1.18-alpine

ADD ./CommandHandler /app
WORKDIR /app

RUN go mod download

RUN go build -o /handler

EXPOSE 8080

CMD [ "/handler" ]
