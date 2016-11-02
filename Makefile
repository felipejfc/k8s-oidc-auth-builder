build:
	@mkdir -p bin
	@go build -o bin/k8s-oidc-auth-builder

run:
	@go run main.go start
