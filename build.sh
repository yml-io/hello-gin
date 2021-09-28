CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -gcflags="-N -l" src/main.go
docker-compose up -d --build