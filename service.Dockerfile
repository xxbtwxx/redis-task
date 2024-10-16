FROM golang:1.23.2-alpine3.20

COPY . .

RUN go build

CMD ["./redis-task"]