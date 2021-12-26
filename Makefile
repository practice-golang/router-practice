build:
	go build -o bin/

all:
	go get -d github.com/mitchellh/gox
	go build -mod=readonly -o ./bin/ github.com/mitchellh/gox
	go mod tidy
	./bin/gox -mod="readonly" -output="bin/{{.Dir}}_{{.OS}}_{{.Arch}}" -osarch="windows/amd64 linux/amd64 linux/arm darwin/amd64 darwin/arm64"
	rm ./bin/gox*

vendor:
	go build -o ./ -mod vendor

modvendor:
	go mod tidy
	go mod vendor

clean:
	rm -rf ./bin/*
