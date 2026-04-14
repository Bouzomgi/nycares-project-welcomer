.PHONY: test fmt lint e2e build clean up down

test:
	go test ./...

fmt:
	gofmt -w .

lint:
	go vet ./...

e2e:
	AWS_ENDPOINT_URL=http://localhost:4566 go test -tags e2e -v -count=1 ./e2e/

build:
	docker compose up --build build-lambdas

clean:
	rm -rf lambda-build/

up:
	docker compose up --build

down:
	docker compose down
