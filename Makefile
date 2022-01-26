database:
	docker-compose up db-init broker

backend:
	air -c scripts/.air.toml

frontend:
	cd web && ember serve
