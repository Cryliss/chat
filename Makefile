build:
	go mod download
	go build --race -o chat bin/chat/chat.go

run:
	./chat -port 8888
