.PHONY:
.SILENT:

build:
	go build -o ./.bin/bot cmd/bot/main.go

run: build
	./.bin/bot

build-image:
	sudo docker build -t binance_signal_bot:v0.1 .

start-container:
	sudo docker run --name binance_signal_bot -p 80:80 --env-file .env binance_signal_bot:v0.1
