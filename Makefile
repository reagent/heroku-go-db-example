DATABASE_NAME = heroku-go-db-example

run: .env create bin/heroku-go-db-example
	@PATH="$(PWD)/bin:$(PATH)" heroku local

.env:
	cp .env.dev .env

bin/heroku-go-db-example: main.go
	go build -o bin/heroku-go-db-example main.go

create:
	@psql -lqt | cut -d \| -f 1 | grep -qw $(DATABASE_NAME) || createdb $(DATABASE_NAME)

drop:
	dropdb $(DATABASE_NAME)

clean:
	rm -rf bin

