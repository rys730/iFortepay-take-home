migrate-up:
# 	goose -dir db/migrations postgres "postgres://user:pass@localhost:5432/app?sslmode=disable" up
	goose -dir db/migrations postgres "postgres://postgres:postgres@localhost:3033/iforte?sslmode=disable" up