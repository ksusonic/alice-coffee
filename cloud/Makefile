PROJECT_NAME=alice-coffee

build:
	go build -o ${PROJECT_NAME} cmd/main.go

run: build
	./${PROJECT_NAME}

clean:
	go clean
	rm ${PROJECT_NAME}

deploy:
	go build -o /usr/local/bin/${PROJECT_NAME} cmd/main.go
	sudo systemctl restart ${PROJECT_NAME}