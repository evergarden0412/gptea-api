build-GPTeaFunction:
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o main cmd/api/main.go
	cp ./main $(ARTIFACTS_DIR)/bootstrap