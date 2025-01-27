# Проект "Go Metrics Alerting"

### Инструкция по развертыванию

1. клонировать репозиторий
2. запустить сервер ```go run cmd/server/main.go -a http://localhost:8080```
3. запустить агент ```go run cmd/agent/main.go -a http://localhost:8080 -r 5 -p 3```