# ==============================================================================
# Docker
dev:
	echo "Starting local environment"
	docker-compose -f docker-compose-dev.yml up --build -d
down:
	docker-compose -f docker-compose-dev.yml down
# ==============================================================================
# Migrations
migrate:
	migrate -path migration/ -database "postgresql://cexles_user:cexles_password@localhost:55005/cexles_db?sslmode=disable" -verbose up

migrate_down:
	migrate -path migration/ -database "postgresql://cexles_user:cexles_password@localhost:55005/cexles_db?sslmode=disable" -verbose down

migrate_fix:
	migrate -path migration/ -database "postgresql://cexles_user:cexles_password@localhost:55005/cexles_db?sslmode=disable" force VERSION