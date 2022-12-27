.PHONY: all go

all: go

go: go/ServiceShortURL.pb.go

clean:
	rm -f shortrpc/*.pb.go
	rmdir shortrpc* 2>/dev/null || true


go/ServiceShortURL.pb.go: ServiceShortURL.proto
	mkdir -p ./shortrpc # make directory for go package
	protoc $$PROTO_PATH --go_opt=paths=source_relative   --go-grpc_out=. --go_out=./shortrpc ServiceShortURL.proto

