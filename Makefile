build:
	go build -o ./bin/aura cmd/aura/main.go

docker:
	go build -o /bin/aura cmd/aura/main.go

all: build
