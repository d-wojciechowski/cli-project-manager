run:
	go run cmd/main.go
fmt:
	go mod tidy
	go fmt ./...
bench:
	go test -v -bench . -run=^# ./...
cover:
	go test -v -cover ./...
test:
	go test -v ./...
build:
	go build  -o dist/app.exe ./main.go
build-small:
	go build  -ldflags '-w -s' -a -installsuffix cgo -o dist/app.exe ./cmd/main.go
