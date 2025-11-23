.PHONY: test docker-run clean

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
