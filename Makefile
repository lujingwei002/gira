
gen:
	protoc --go_out=. --go-grpc_out=.  service/admin/admin.proto

cli:
	go build -o gira-cli bin/cli/main.go

.PHONY: gen cli
