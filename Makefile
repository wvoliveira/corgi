database:
	docker-compose up db-init broker

backend:
	# air -c scripts/.air.toml
	go run .\cmd\corgi\main.go

frontend:
	cd web && ember serve
