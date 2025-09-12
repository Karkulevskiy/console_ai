.PHONY: all server client run

all: server client

server:
	go run ./

client:
	go run ./client/

run:
	go run ./
	go run ./client/
