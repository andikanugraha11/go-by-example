protoc user\v1\user.proto --go_out=paths=source_relative:./gen 
protoc user\v1\*.proto --go_out=paths=source_relative:./gen --go-grpc_out=paths=source_relative:./gen