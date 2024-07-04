build:
	docker build . -f Dockerfile -t antohachaban/news-alligator-web

run:
	docker run -d -v news-aggregator-backups:/root/backups -p 8080:8080 antohachaban/news-alligator-web

stop:
	docker stop $(shell docker ps -q --filter ancestor=antohachaban/news-alligator-web)

push:
	docker push antohachaban/news-alligator-web

test:
	go test -v ./...

pull:
	docker pull antohachaban/news-alligator-web

dev-up:
	docker-compose -f devbox/docker-compose.yml up -d

dev-down:
	docker-compose -f devbox/docker-compose.yml down

dev-build:
	docker-compose -f devbox/docker-compose.yml build

swag:
	swag init -d .,../../../pkg/handler/web,../../../pkg/model

run-postgres:
	docker run --name=news-alligator-db -e POSTGRES_PASSWORD='qwerty' -p 5436:5432 -d --rm postgres

migrate:
	migrate -path ./schema -database 'postgres://postgres:qwerty@0.0.0.0:5436/postgres?sslmode=disable' up