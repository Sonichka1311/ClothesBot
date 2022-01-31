build:
	GOOS=linux GOARCH=amd64 go build -o bin/bot/bot cmd/bot/main.go

run: build
	docker-compose up --build --force-recreate

buildmacos:
	go build -o bin/bot/bot cmd/bot/main.go

localrun: buildmacos
	./bin/bot/bot

cloudrun: build
	./bin/bot/bot