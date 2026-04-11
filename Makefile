.PHONY: test fmt lint integration build clean up down

test:
	go test ./...

fmt:
	gofmt -w .

lint:
	go vet ./...

integration:
	AWS_ENDPOINT_URL=http://localhost:4566 go test -tags integration -v -count=1 ./integration/

build:
	docker compose up --build build-lambdas

clean:
	rm -rf lambda-build/

up:
	docker compose up --build

down:
	docker compose down
