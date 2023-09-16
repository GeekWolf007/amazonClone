build:
	@go build -o bin/golang_amazon

run: build
	@./bin/golang_amazon
