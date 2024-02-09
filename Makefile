compile:
	PATH=$$PWD/local/bin:$$PATH \
		protoc --go_out=. --go_opt=paths=source_relative \
		internal/test/proto2_test.proto \
		internal/test/proto3_test.proto

prepare:
	mkdir -p local/bin
	GOBIN=$$PWD/local/bin go install google.golang.org/protobuf/cmd/protoc-gen-go
