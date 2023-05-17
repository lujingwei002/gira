
all: cli service framework

framework:
	make -C framework/smallgame 

cli:
	go build -o gira-cli bin/cli/main.go
	go build -o protoc-gen-go-gclient bin/gen_gclient/*.go
	cp -rf protoc-gen-go-gclient ~/go/bin/protoc-gen-go-gclient

service:
	protoc --go_out=service/admin --go-grpc_out=service/admin --go-gclient_out=service/admin  service/admin/admin.proto
	protoc --go_out=service/peer --go-grpc_out=service/peer --go-gclient_out=service/peer service/peer/peer.proto
	protoc --go_out=service/hall --go-grpc_out=service/hall --go-gclient_out=service/hall service/hall/hall.proto


.PHONY: gen cli service framework
