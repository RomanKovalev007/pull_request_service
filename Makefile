.PHONY: test docker-run clean full_clean

test:
	docker-compose -f docker-compose.test.yml up --build

docker-run:
	docker-compose up --build

clean:
	docker-compose down
	docker-compose -f docker-compose.test.yml down -v 

full_clean:
	docker-compose down -v
	docker-compose -f docker-compose.test.yml down -v 

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: lint-fix
lint-fix:
	golangci-lint run --fix ./...