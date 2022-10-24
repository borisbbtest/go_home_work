.PHONY: all go

all: go

go: go/ServiceShortURL.pb.go

clean:
	rm -f shortserver/*.pb.go
	rmdir shortserver* 2>/dev/null || true


go/ServiceShortURL.pb.go: ServiceShortURL.proto
	mkdir -p ./shortserver # make directory for go package
	protoc $$PROTO_PATH --go_opt=paths=source_relative --go_out=./shortserver ServiceShortURL.proto

