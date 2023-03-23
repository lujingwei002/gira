gen:
	protoc --go_out=. --go-grpc_out=.  service/admin/admin.proto

.PHONY: gen 
