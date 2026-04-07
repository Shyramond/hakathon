# Все тесты
go test ./internal/... -v -count=1

# Только сервисы
go test ./internal/service/... -v

# Только middleware
go test ./internal/middleware/... -v

# С покрытием
go test ./internal/... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html

# Конкретный тест
go test ./internal/service/ -run TestPurchase -v