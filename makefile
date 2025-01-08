setup:
	sudo apt install -y protobuf-compiler
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/air-verse/air@latest

generate:
	protoc\
		--go_out=gateway/grpc/gen\
		--go_opt=paths=source_relative\
		--go-grpc_out=gateway/grpc/gen\
		--go-grpc_opt=paths=source_relative\
		othello.proto

run:
	air
