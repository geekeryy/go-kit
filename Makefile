
g:
	stringer -type Code  -linecomment ./pkg/xerror/code.go
proto:
	protoc -I=. --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. ./grpc/*/*.proto

t:export APP_ENV=dev
t:
	go test -coverprofile=c.out -coverpkg=./... ./...  && go tool cover -html=c.out