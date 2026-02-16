.PHONY: test fmt lint integration build clean up down

test:
	go test ./...

fmt:
	gofmt -w .

lint:
	go vet ./...

integration:
	go test -tags integration -v -count=1 ./integration/

build:
	docker compose up build-lambdas

clean:
	rm -rf lambda-build/

up:
	docker compose up --build

down:
	docker compose down
