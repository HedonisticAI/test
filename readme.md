# Test
Test for effective mobile
## Requered
1) Goose for migrations
2) Docker compose
3) Golang (my ver 1.24.2)
4) Insomnia (optional)
## How to start
1) Run DB 
```bash
docker compose up
```
2) Run migrations 
```bash
goose postgres "postgresql://goose:password@127.0.0.1:5432/test?sslmode=disable" -dir db/migrations -table public.goose_migrations up
```
```bash
goose -dir db/migrations postgres  "postgresql://goose:password@127.0.0.1:5432/test?sslmode=disable" up
```
If you changed your db - change
>"postgresql://goose:password@127.0.0.1:5432/test?sslmode=disable"
as well
3) Install dependencies
```bash
go mod tidy
```
4) Run App
```bash 
go run cmd/main.go
```
5) Export request collection to insomnia (optional)