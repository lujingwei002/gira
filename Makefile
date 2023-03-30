
gen:
	protoc --go_out=. --go-grpc_out=.  service/admin/admin.proto
	make -C framework/smallgame gen


cli:
	go build -o gira-cli bin/cli/main.go

.PHONY: gen cli
