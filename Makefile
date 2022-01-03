build:
	go build -ldflags "-w -s" -trimpath -o bin/

all:
	go get -d github.com/mitchellh/gox
	go build -mod=readonly -o ./bin/ github.com/mitchellh/gox
	go mod tidy
	./bin/gox -mod="readonly" -output="bin/{{.Dir}}_{{.OS}}_{{.Arch}}" -osarch="windows/amd64 linux/amd64 linux/arm darwin/amd64 darwin/arm64"
	rm ./bin/gox*

test:
	go test ./...

test-v:
	go test -v ./...

clean:
	rm -rf ./bin/*
