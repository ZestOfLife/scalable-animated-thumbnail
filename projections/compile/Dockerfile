FROM golang:1.18-alpine

ADD ./CompileProjector /app
WORKDIR /app

RUN go mod download

RUN go build -o /compiler

CMD [ "/compiler" ]
