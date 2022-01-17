database:
	docker-compose up db-init

backend:
	air -c scripts/.air.toml

frontend:
	cd web && ember serve
