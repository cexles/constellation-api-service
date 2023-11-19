
migrate:
	migrate -path migration/ -database "postgresql://cexles:cexles@localhost:55002/cexles-test?sslmode=disable" -verbose up

migrate_down:
	migrate -path migration/ -database "postgresql://cexles:cexles@localhost:55002/cexles-test?sslmode=disable" -verbose down

migrate_fix:
	migrate -path migration/ -database "postgresql://cexles:cexles@localhost:55002/cexles-test?sslmode=disable" force VERSION