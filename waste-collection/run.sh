FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod ./
COPY main.go ./
COPY run.sh ./

RUN go build -o waste-collection main.go

CMD ["./run.sh"]
