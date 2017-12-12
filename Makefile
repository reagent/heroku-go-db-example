run: .env bin/heroku-go-db-example
	@PATH="$(PWD)/bin:$(PATH)" heroku local

.env:
	cp .env.dev .env

bin/heroku-go-db-example: main.go
	go build -o bin/heroku-go-db-example main.go

clean:
	rm -rf bin

