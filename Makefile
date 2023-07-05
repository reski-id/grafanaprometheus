swag-init:
	swag init --parseInternal --exclude build,deploy,docs,scripts,vendor -g main.go

swag-fmt:
	swag fmt --exclude build,developments,docs,scripts -g main.go 

up:
	docker-compose up

down:
	docker-compose down

run:
	go run .\my-app\main.go
	