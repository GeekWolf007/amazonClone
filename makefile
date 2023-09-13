build:
	@go build -o bin/golang_BANK

run: build
	@./bin/golang_BANK
