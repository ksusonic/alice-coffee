PROJECT_NAME=alice-coffee

build:
	protoc --go_out=$GOPATH/pkg/mod/ ../protos/vending/*.proto -I=../
	go build -o ${PROJECT_NAME} cmd/main.go

run: build
	./${PROJECT_NAME}

clean:
	go clean
	rm ${PROJECT_NAME}

deploy:
	go build -o /usr/local/bin/${PROJECT_NAME} cmd/main.go
	sudo systemctl restart ${PROJECT_NAME}