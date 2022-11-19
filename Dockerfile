FROM golang:alpine

COPY . /github.com/Alishreder/binanceSignalBot/
WORKDIR /github.com/Alishreder/binanceSignalBot/

RUN go mod download
RUN go build -o ./bin/bot cmd/bot/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=0 /github.com/Alishreder/binanceSignalBot/bin/bot .

EXPOSE 80

CMD ["./bot"]