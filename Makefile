include configs/dc.env

.PHONY: init dc-local-up dc-local-down local-run dev-run dev-stop run stop

dc-local-up:
	docker-compose -f ./deployments/docker-compose.local.yaml --env-file configs/dc.env up -d

dc-local-down:
	docker-compose -f ./deployments/docker-compose.local.yaml --env-file configs/dc.env down

local-run: dc-local-up
	go run ./cmd/archive_bot/main.go -config ./configs/local.yaml

dev-run:
	docker-compose -f ./deployments/docker-compose.dev.yaml --env-file configs/dc.env up

dev-stop:
	docker-compose -f ./deployments/docker-compose.dev.yaml --env-file configs/dc.env down --rmi local

run:
	docker-compose -f ./deployments/docker-compose.yaml --env-file configs/dc.env up

stop:
	docker-compose -f ./deployments/docker-compose.yaml --env-file configs/dc.env down --rmi local
