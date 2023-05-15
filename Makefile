
gen:
	protoc --go_out=. --go-grpc_out=.  service/admin/admin.proto
	make -C framework/smallgame gen

cli:
	go build -o gira-cli bin/cli/main.go
	go build -o protoc-gen-go-gira bin/gen_grpc/*.go
	cp -rf protoc-gen-go-gira ~/go/bin/protoc-gen-go-gira

.PHONY: gen cli
