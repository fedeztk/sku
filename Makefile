all:
	go generate ./...
	go build -o sku cmd/sku/main.go

run:
	go run cmd/sku/main.go

clean:
	@if [ -f sku ] && [ -x sku ]; then \
		rm sku; \
	fi

docker-build:
	docker build -t sku .

docker-run:
	docker run -it -e "TERM=xterm-256color" sku

install:
	go generate ./...
	go build -o sku cmd/sku/main.go
	mv -f sku `go env GOPATH`/bin/

uninstall:
	rm `go env GOPATH`/bin/sku
