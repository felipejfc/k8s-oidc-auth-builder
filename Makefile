build:
	@mkdir -p bin
	@go build -o bin/k8s-oidc-auth-builder

build-cross: build-cross-darwin build-cross-linux build-exec

build-exec:
	@chmod u+x bin/*

build-cross-darwin:
	@mkdir -p ./bin
	@echo "Building for darwin-i386..."
	@env GOOS=darwin GOARCH=386 go build -o ./bin/k8s-oidc-auth-builder-darwin-i386 ./main.go
	@echo "Building for darwin-x86_64..."
	@env GOOS=darwin GOARCH=amd64 go build -o ./bin/k8s-oidc-auth-builder-darwin-x86_64 ./main.go

build-cross-linux:
	@mkdir -p ./bin
	@echo "Building for linux-i386..."
	@env GOOS=linux GOARCH=386 go build -o ./bin/k8s-oidc-auth-builder-linux-i386 ./main.go
	@echo "Building for linux-x86_64..."
	@env GOOS=linux GOARCH=amd64 go build -o ./bin/k8s-oidc-auth-builder-linux-x86_64 ./main.go

image:
	docker build -t k8s-oidc-auth-builder .

run:
	@go run main.go start -d
