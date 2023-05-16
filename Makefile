
gen:
	protoc --go_out=service/admin --go-grpc_out=service/admin --go-gclient_out=service/admin  service/admin/admin.proto
	make -C framework/smallgame gen

cli:
	go build -o gira-cli bin/cli/main.go
	go build -o protoc-gen-go-gclient bin/gen_gclient/*.go
	cp -rf protoc-gen-go-gclient ~/go/bin/protoc-gen-go-gclient

.PHONY: gen cli
